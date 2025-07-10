package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	PROJECT_SCAN_DIR_NAME = "project-to-scan"
	REPORTS_DIR_NAME      = "reports"
	REPORTS_FILE_NAME     = "semgrep-report.json"

	APP_DATA_DIR_NAME = "GoCodeCheckDesktop"
)

func CheckIfFolderOrFileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {

		if os.IsNotExist(err) {
			return false, fmt.Errorf("the path '%s' does not exist", path)
		}

		return false, fmt.Errorf("error checking path '%s': %w", path, err)
	}

	return true, nil
}

// cleanDir removes all files and directories in the given directory without removing the directory itself
func CleanDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		if err := os.RemoveAll(filepath.Join(dir, name)); err != nil {
			return err
		}
	}

	return nil
}

// copyDir recursively copies a directory tree
func CopyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Skip .git directory to avoid unnecessary files
			if entry.Name() == ".git" {
				continue
			}
			if err = CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy regular files
			if err = CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file from src to dst
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Don't try to copy very large files
	if srcInfo.Size() > 100*1024*1024 { // Skip files larger than 100MB
		return nil
	}

	dstFile, err := os.CreateTemp(filepath.Dir(dst), "tmp_")
	if err != nil {
		return err
	}

	_, err = io.Copy(dstFile, srcFile)
	dstFile.Close()
	if err != nil {
		os.Remove(dstFile.Name())
		return err
	}

	if err = os.Chmod(dstFile.Name(), srcInfo.Mode()); err != nil {
		os.Remove(dstFile.Name())
		return err
	}

	if err = os.Rename(dstFile.Name(), dst); err != nil {
		os.Remove(dstFile.Name())
		return err
	}

	return nil
}

// getAbsolutePath converts a relative path to an absolute path
func GetAbsolutePath(path string) (string, error) {
	// If path is already absolute, return it
	if filepath.IsAbs(path) {
		return path, nil
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Join with the current directory to make it absolute
	absPath := filepath.Join(cwd, path)
	return absPath, nil
}

// getBaseDir returns the appropriate base directory for storing application data
// In development, it uses a local directory; in production, it uses a user-specific directory
func GetBaseDir() string {
	// Try to use a directory in user's home for production
	userConfigDir, err := os.UserConfigDir()

	if err == nil {
		// Use a directory in the user's config folder (e.g., AppData on Windows)
		appDataDir := filepath.Join(userConfigDir, APP_DATA_DIR_NAME)

		// Check/create the directory
		if err := os.MkdirAll(appDataDir, 0755); err == nil {
			return appDataDir
		}
	}

	// Fallback: try to use system temp directory
	tempBaseDir := filepath.Join(os.TempDir(), APP_DATA_DIR_NAME)
	if err := os.MkdirAll(tempBaseDir, 0755); err == nil {
		return tempBaseDir
	}

	// Last resort: use a local directory (for development)
	return filepath.Join(".", "backend", "core")
}

// ensureDirectory makes sure a directory exists, creating it if necessary
func EnsureDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, 0755)
	}
}
