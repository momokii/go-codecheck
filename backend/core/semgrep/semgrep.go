package semgrep

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/momokii/go-codecheck/backend/pkg/utils"
)

var (
	PROJECT_SCAN_TEMP_FOLDER_PATH string
	REPORTS_FOLDER_PATH           string

	// scan management
	currentScanCtx    context.Context
	currentScanCancel context.CancelFunc
	scanMutex         sync.Mutex
)

func init() {
	// Initialize paths based on environment (development vs production)
	baseDir := utils.GetBaseDir()

	PROJECT_SCAN_TEMP_FOLDER_PATH = filepath.Join(baseDir, utils.PROJECT_SCAN_DIR_NAME)

	REPORTS_FOLDER_PATH = filepath.Join(baseDir, utils.REPORTS_DIR_NAME)

	utils.EnsureDirectory(PROJECT_SCAN_TEMP_FOLDER_PATH)
	utils.EnsureDirectory(REPORTS_FOLDER_PATH)
}

func InitializeAndPrepareFolderScanSemgrep(projectDir string) error {
	// check directory input user is available
	if _, err := utils.CheckIfFolderOrFileExists(projectDir); err != nil {
		return fmt.Errorf("project path error: %w", err)
	}

	// make sure this temp folder exist
	if _, err := os.Stat(PROJECT_SCAN_TEMP_FOLDER_PATH); os.IsNotExist(err) {
		if err := os.MkdirAll(PROJECT_SCAN_TEMP_FOLDER_PATH, 0755); err != nil {
			return fmt.Errorf("failed to create scan folder: %w", err)
		}
	}

	// clean temp folder
	if err := utils.CleanDir(PROJECT_SCAN_TEMP_FOLDER_PATH); err != nil {
		return fmt.Errorf("failed to clean scan folder: %w", err)
	}

	// copy all files from user input dir path to project temp dir
	if err := utils.CopyDir(projectDir, PROJECT_SCAN_TEMP_FOLDER_PATH); err != nil {

		// if error happen here, make sure to cleanup the temp dir
		utils.CleanDir(PROJECT_SCAN_TEMP_FOLDER_PATH)

		return fmt.Errorf("failed to copy project files to scan folder: %w", err)
	}

	return nil
}

func RunSemgrepScan(usingAllRules bool) (*ScanResult, error) {
	scanMutex.Lock()

	// Create cancellable context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)

	// Store current scan context for cancellation
	currentScanCtx = ctx
	currentScanCancel = cancel

	scanMutex.Unlock()

	defer func() {
		scanMutex.Lock()
		currentScanCtx = nil
		currentScanCancel = nil
		scanMutex.Unlock()

		// Cleanup temp directory
		utils.CleanDir(PROJECT_SCAN_TEMP_FOLDER_PATH)
	}()

	// Make sure directories exist
	utils.EnsureDirectory(PROJECT_SCAN_TEMP_FOLDER_PATH)
	utils.EnsureDirectory(REPORTS_FOLDER_PATH)

	// Get absolute paths for mounting volumes - Docker requires absolute paths for volume mounts
	absProjectPath, err := filepath.Abs(PROJECT_SCAN_TEMP_FOLDER_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute project path: %w", err)
	}

	absReportsPath, err := filepath.Abs(REPORTS_FOLDER_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute reports path: %w", err)
	}

	// Build direct docker run command - equivalent to the compose file
	config_rules := "r/all"
	if !usingAllRules {
		config_rules = "auto"
	}
	cmd := exec.CommandContext(ctx,
		"docker", "run", "--rm",
		"--volume", fmt.Sprintf("%s:/src:ro", absProjectPath),
		"--volume", fmt.Sprintf("%s:/reports:rw", absReportsPath),
		"--workdir", "/src",
		"semgrep/semgrep:latest",
		"semgrep", "scan",
		"--config", config_rules,
		"--json", "-o", "/reports/semgrep-report.json", ".",
	)

	// get output
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// Run the command in a goroutine to handle cancellation properly
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	// Wait for either completion or cancellation
	select {
	case err := <-done:
		// Command completed normally
		exitCode := 0
		if err != nil {
			// Check if it was cancelled first
			if ctx.Err() == context.Canceled {
				return nil, fmt.Errorf("Scan Was Cancelled by User")
			}

			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				return nil, fmt.Errorf("Failed to run docker command: %w", err)
			}
		}

		return &ScanResult{
			ExitCode: exitCode,
			Stdout:   outBuf.String(),
			Stderr:   errBuf.String(),
		}, nil

	case <-ctx.Done():
		// Context was cancelled - kill the process immediately
		if cmd.Process != nil {
			cmd.Process.Kill()
		}

		// Also try to kill any running semgrep containers as backup
		go func() {
			exec.Command("docker", "stop", "$(docker ps -q --filter ancestor=semgrep/semgrep:latest)").Run()
		}()

		return nil, fmt.Errorf("Scan Was Cancelled by User")
	}
}

func CancelCurrentScan() error {
	scanMutex.Lock()
	defer scanMutex.Unlock()

	if currentScanCancel != nil {
		// Cancel the context immediately
		currentScanCancel()

		// Force kill any Docker containers that might be running
		go func() {
			// Kill any running semgrep containers with more aggressive approach
			exec.Command("docker", "kill", "$(docker ps -q --filter ancestor=semgrep/semgrep:latest)").Run()
			exec.Command("docker", "stop", "$(docker ps -q --filter ancestor=semgrep/semgrep:latest)").Run()
		}()

		return nil
	}

	return fmt.Errorf("No Active Scan to Cancel")
}

func IsScanRunning() bool {
	scanMutex.Lock()
	defer scanMutex.Unlock()

	// Check if context exists and is not cancelled/timed out
	if currentScanCtx == nil {
		return false
	}

	// If context is cancelled or timed out, consider scan as not running
	select {
	case <-currentScanCtx.Done():
		return false
	default:
		return true
	}
}
