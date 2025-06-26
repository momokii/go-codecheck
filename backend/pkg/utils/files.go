package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
