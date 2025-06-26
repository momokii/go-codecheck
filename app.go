package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/momokii/go-codecheck/backend/core/docker"
	"github.com/momokii/go-codecheck/backend/core/semgrep"
	"github.com/momokii/go-codecheck/backend/pkg/parser"
	"github.com/momokii/go-codecheck/backend/pkg/utils"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) CheckIfFolderOrFIleExists(path string) (bool, error) {
	exists, err := utils.CheckIfFolderOrFileExists(path)
	if err != nil {
		return false, fmt.Errorf("error checking path '%s': %w", path, err)
	}
	return exists, nil
}

func (a *App) CheckDockerIsAvailable() (string, error) {
	dockerVersion, err := docker.CheckDockerIsAvailable()
	if err != nil {
		return "", fmt.Errorf("failed to check Docker version: %w", err)
	}
	return dockerVersion, nil
}

func (a *App) CheckDockerImagesIsAvailable(imageName, imageTag string) (docker.DockerImage, error) {
	image, err := docker.CheckDockerImageIsAvailable(imageName, imageTag)
	if err != nil {
		return docker.DockerImage{}, fmt.Errorf("failed to check Docker image '%s:%s': %w", imageName, imageTag, err)
	}

	return docker.DockerImage{
		Repository: image.Repository,
		Tag:        image.Tag,
		ImageID:    image.ImageID,
	}, nil
}

func (a *App) InitAndPrepareFolderScanSemgrep(targetPath string) error {
	return semgrep.InitializeAndPrepareFolderScanSemgrep(targetPath)
}

func (a *App) RunSemgrepScan() (*semgrep.ScanResult, error) {
	return semgrep.RunSemgrepScan()
}

func (a *App) GetSemgrepReportData() (*parser.SemgrepReport, error) {

	path := filepath.Join(semgrep.REPORTS_FOLDER_PATH, semgrep.REPORTS_FILE_NAME)

	return parser.ParseSemgrepReport(path)
}
