package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/momokii/go-codecheck/backend/core/databases"
	"github.com/momokii/go-codecheck/backend/core/docker"
	"github.com/momokii/go-codecheck/backend/core/handlers"
	"github.com/momokii/go-codecheck/backend/core/models"
	repository "github.com/momokii/go-codecheck/backend/core/repository/repo"
	"github.com/momokii/go-codecheck/backend/core/repository/scan"
	"github.com/momokii/go-codecheck/backend/core/repository/user"
	"github.com/momokii/go-codecheck/backend/core/semgrep"
	"github.com/momokii/go-codecheck/backend/pkg/parser"
	"github.com/momokii/go-codecheck/backend/pkg/utils"
)

// App struct
type App struct {
	ctx         context.Context
	scanHandler handlers.ScanHandler
	repoHandler handlers.RepoHandler
	userHandler handlers.UserHandler
}

// NewApp creates a new App application struct
func NewApp() *App {

	// init db
	if err := databases.InitDatabaseSQLite(); err != nil {
		panic(fmt.Sprintf("Database initialization failed: %v", err))
	}

	db, err := databases.NewSQLiteDatabases(databases.DATABASE_SQLITE_PATH)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// init repo
	repoRepo := repository.NewRepository()
	scanRepo := scan.NewScanRepository()
	userRepo := user.NewUserRepository()

	//  init app with handler
	return &App{
		scanHandler: *handlers.NewScanHandler(db.GetDB(), scanRepo, repoRepo),
		repoHandler: *handlers.NewRepoHandler(db.GetDB(), repoRepo),
		userHandler: *handlers.NewUserHandler(db.GetDB(), userRepo),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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

	path := filepath.Join(semgrep.REPORTS_FOLDER_PATH, utils.REPORTS_FILE_NAME)

	return parser.ParseSemgrepReport(path)
}

// ! ============================  FUNCTION FROM HANDLERS

// ? =========== REPO
type RepoDataPaginationFE struct {
	Data  []models.Repository `json:"data"`
	Total int                 `json:"total"`
}

func (a *App) GetRepoDatas(userId, page, perPage int, search string, desc_sort bool) (RepoDataPaginationFE, error) {

	datas, total, err := a.repoHandler.GetReposData(userId, page, perPage, search, desc_sort)

	data := RepoDataPaginationFE{
		Data:  datas,
		Total: total,
	}

	return data, err
}

func (a *App) GetRepoById(userId, repoId int) (models.Repository, error) {
	return a.repoHandler.GetOneRepoById(userId, repoId)
}

func (a *App) CreateNewRepo(newRepoData models.RepositoryCreate) error {
	return a.repoHandler.CreateNewRepo(&newRepoData)
}

func (a *App) UpdateRepo(repoId int, updateData models.RepositoryUpdate) error {
	return a.repoHandler.UpdateRepo(&updateData)
}

func (a *App) DeleteRepo(repoId, userId int) error {
	return a.repoHandler.DeleteRepo(repoId, userId)
}

// ? =========== SCAN
type ScanDataPaginationFE struct {
	Data  []models.ScanFull `json:"data"`
	Total int               `json:"total"`
}

func (a *App) GetScanDatas(repoId, page, perPage int, search string, desc_sort bool) (ScanDataPaginationFE, error) {
	datas, total, err := a.scanHandler.GetScansData(repoId, page, perPage, search, desc_sort)

	data := ScanDataPaginationFE{
		Data:  datas,
		Total: total,
	}

	return data, err
}

func (a *App) GetScanById(repoId, scanId int) (models.ScanFull, error) {
	return a.scanHandler.GetScanById(repoId, scanId)
}

func (a *App) CreateNewScan(repoId int, newScanData models.ScanCreate) error {
	return a.scanHandler.CreateNewScan(&newScanData)
}

func (a *App) DeleteScan(repoId, scanId int) error {
	return a.scanHandler.DeleteScan(repoId, scanId)
}

func (a *App) GetScanDetail(jsonString string) (*parser.SemgrepReport, error) {
	return parser.ParseResultSemgrepFromDatabase(jsonString)
}

// ? =========== USER
func (a *App) GetAndValidateUserByToken(token string) (models.User, error) {
	return a.userHandler.GetAndValidateUserDataByToken(token)
}

func (a *App) Login(userLogin models.UserLogin) (token string, err error) {
	return a.userHandler.UserLogin(userLogin)
}

func (a *App) UpdateUserPassword(userUpdate models.UserUpdate, isCompletedSetup bool) error {
	return a.userHandler.UserUpdatePassword(userUpdate, isCompletedSetup)
}

func (a *App) UpdateUserUsername(userUpdate models.UserUpdate) error {
	return a.userHandler.UserUpdateUsername(userUpdate)
}

func (a *App) CompleteFirstSetup(userId int) error {
	return a.userHandler.UserCompleteFirstSetup(userId)
}

func (a *App) Logout(userId int) error {
	return a.userHandler.UserDeleteSession(userId)
}
