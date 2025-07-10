package semgrep

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/momokii/go-codecheck/backend/pkg/utils"
)

var (
	PROJECT_SCAN_TEMP_FOLDER_PATH string
	REPORTS_FOLDER_PATH           string
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

func RunSemgrepScan() (*ScanResult, error) {
	// make sure to cleanup the temp dir after this function
	defer utils.CleanDir(PROJECT_SCAN_TEMP_FOLDER_PATH)

	// context for timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

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
	cmd := exec.CommandContext(ctx,
		"docker", "run", "--rm",
		"--volume", fmt.Sprintf("%s:/src:ro", absProjectPath),
		"--volume", fmt.Sprintf("%s:/reports:rw", absReportsPath),
		"--workdir", "/src",
		"semgrep/semgrep:latest",
		"semgrep", "scan", "--config", "auto", "--json", "-o", "/reports/semgrep-report.json", ".",
	)

	// get output
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("failed to run docker command: %w", err)
		}
	}

	return &ScanResult{
		ExitCode: exitCode,
		Stdout:   outBuf.String(),
		Stderr:   errBuf.String(),
	}, nil
}
