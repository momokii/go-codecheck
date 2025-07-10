package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/momokii/go-codecheck/backend/core/databases"
	"github.com/momokii/go-codecheck/backend/core/models"
	repository "github.com/momokii/go-codecheck/backend/core/repository/repo"
	"github.com/momokii/go-codecheck/backend/pkg/utils"
)

type RepoHandler struct {
	sqliteDB       *databases.SQLiteDB
	RepositoryRepo *repository.Repository
}

func NewRepoHandler(sqliteDB *databases.SQLiteDB, repoRepo *repository.Repository) *RepoHandler {
	return &RepoHandler{
		sqliteDB:       sqliteDB,
		RepositoryRepo: repoRepo,
	}
}

func (h *RepoHandler) GetReposData(userId, page, perPage int, search string, desc_sort bool) ([]models.Repository, int, error) {

	var repos []models.Repository
	var total int

	if page <= 0 {
		page = 1
	}

	if perPage <= 0 {
		perPage = 10
	}

	sortBy := models.SortByAsc
	if desc_sort {
		sortBy = models.SortByDesc
	}

	paginationFilter := models.PaginationFiltering{
		Page:    page,
		PerPage: perPage,
		Search:  search,
		SortBy:  sortBy,
	}

	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		reposData, total_data, err := h.RepositoryRepo.Find(tx, userId, paginationFilter)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		total = total_data
		repos = *reposData

		return nil, http.StatusOK

	}); err != nil {
		return repos, total, err
	}

	return repos, total, nil
}

func (h *RepoHandler) GetOneRepoById(userId, repoId int) (models.Repository, error) {

	var repo models.Repository

	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		repoData, err := h.RepositoryRepo.FindByParam(
			tx,
			userId,
			repoId,
			"",
			"",
		)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		if repoData.Id == 0 {
			return fmt.Errorf("repository with id '%d' not found", repoId), http.StatusNotFound
		}

		// pass the data
		repo = *repoData

		return nil, http.StatusOK

	}); err != nil {
		return repo, err
	}

	return repo, nil
}

func (h *RepoHandler) CreateNewRepo(newRepoData *models.RepositoryCreate) error {

	if err := utils.ValidateStruct(newRepoData); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Name":
				return fmt.Errorf("validation error: %s is required", err.Field())
			case "Description":
				return fmt.Errorf("validation error: %s is required", err.Field())
			case "UserId":
				return fmt.Errorf("validation error: %s is required", err.Field())
			case "Path":
				return fmt.Errorf("validation error: %s is required", err.Field())
			}
		}
	}

	// check if path is valid
	if _, err := utils.CheckIfFolderOrFileExists(newRepoData.Path); err != nil {
		return fmt.Errorf("validation error: %s is not a valid path", newRepoData.Path)
	}

	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		// check if exist on same name
		existingRepo, err := h.RepositoryRepo.FindByParam(tx, newRepoData.UserId, 0, newRepoData.Name, newRepoData.Path)
		if err != nil {
			return err, http.StatusBadRequest
		}

		if existingRepo.Id != 0 {
			return fmt.Errorf("repository with name '%s' | '%s' already exists, make sure to have unique value for name and path", newRepoData.Name, newRepoData.Path), http.StatusConflict
		}

		// if not, create new one
		if _, err := h.RepositoryRepo.Create(tx, newRepoData); err != nil {
			return fmt.Errorf("failed to create repository: %w", err), http.StatusInternalServerError
		}

		return nil, http.StatusCreated

	}); err != nil {
		return err
	}

	return nil
}

func (h *RepoHandler) UpdateRepo(dataEditRepo *models.RepositoryUpdate) error {

	if err := utils.ValidateStruct(dataEditRepo); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Id":
				return fmt.Errorf("validation error: %s is required", err.Field())
			case "UserId":
				return fmt.Errorf("validation error: %s is required", err.Field())
			case "Name":
				return fmt.Errorf("validation error: %s is required", err.Field())
			case "Description":
				return fmt.Errorf("validation error: %s is required", err.Field())
			case "Path":
				return fmt.Errorf("validation error: %s is required", err.Field())
			}
		}
	}

	// check path validation
	if _, err := utils.CheckIfFolderOrFileExists(dataEditRepo.Path); err != nil {
		return fmt.Errorf("validation error: %s is not a valid path", dataEditRepo.Path)
	}

	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		// check if data exist
		repoData, err := h.RepositoryRepo.FindByParam(
			tx,
			dataEditRepo.UserId,
			dataEditRepo.Id,
			"",
			"",
		)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		if repoData.Id == 0 {
			return fmt.Errorf("repository with id '%d' not found", dataEditRepo.Id), http.StatusNotFound
		}

		// check if name or path already exist
		dataCheck, err := h.RepositoryRepo.FindByParam(
			tx,
			dataEditRepo.UserId,
			dataEditRepo.Id,
			dataEditRepo.Name,
			dataEditRepo.Path,
		)
		if err != nil {
			return fmt.Errorf("failed to check repository: %w", err), http.StatusInternalServerError
		}

		if dataCheck.Id != 0 && dataCheck.Id != repoData.Id {
			return fmt.Errorf("repository with name '%s' | '%s' already exists in another data, make sure to have unique value for name and path", dataEditRepo.Name, dataEditRepo.Path), http.StatusConflict
		}

		// update data
		if err := h.RepositoryRepo.Update(tx, dataEditRepo); err != nil {
			return fmt.Errorf("failed to update repository: %w", err), http.StatusInternalServerError
		}

		return nil, http.StatusOK

	}); err != nil {
		return err
	}

	return nil
}

func (h *RepoHandler) DeleteRepo(repoId, userId int) error {

	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		// check if repo exist
		repoData, err := h.RepositoryRepo.FindByParam(
			tx,
			userId,
			repoId,
			"",
			"",
		)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		if repoData.Id == 0 {
			return fmt.Errorf("repository with id '%d' not found in your data", repoId), http.StatusNotFound
		}

		// delete repo
		if err := h.RepositoryRepo.Delete(tx, repoId, userId); err != nil {
			return fmt.Errorf("failed to delete repository: %w", err), http.StatusInternalServerError
		}

		return nil, http.StatusOK

	}); err != nil {
		return err
	}

	return nil
}
