package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/momokii/go-codecheck/backend/core/databases"
	"github.com/momokii/go-codecheck/backend/core/models"
	repository "github.com/momokii/go-codecheck/backend/core/repository/repo"
	"github.com/momokii/go-codecheck/backend/core/repository/scan"
	"github.com/momokii/go-codecheck/backend/pkg/utils"
)

type ScanHandler struct {
	sqliteDB *databases.SQLiteDB
	ScanRepo *scan.ScanRepository
	RepoRepo *repository.Repository
}

func NewScanHandler(sqliteDB *databases.SQLiteDB, scanRepo *scan.ScanRepository, repoRepo *repository.Repository) *ScanHandler {
	return &ScanHandler{
		sqliteDB: sqliteDB,
		ScanRepo: scanRepo,
		RepoRepo: repoRepo,
	}
}

func (h *ScanHandler) GetScansData(repoId, page, perPage int, search string, desc_sort bool) ([]models.ScanFull, int, error) {

	var scans []models.ScanFull
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

		scanDatas, total_data, err := h.ScanRepo.Find(tx, repoId, paginationFilter)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		total = total_data
		scans = *scanDatas

		return nil, http.StatusOK

	}); err != nil {
		return scans, total, err
	}

	return scans, total, nil
}

func (h *ScanHandler) GetScanById(repoId, scanId int) (models.ScanFull, error) {

	var scan models.ScanFull

	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		scanData, err := h.ScanRepo.FindById(tx, repoId, scanId)
		if err != nil {
			return fmt.Errorf("failed to find scan: %w", err), http.StatusInternalServerError
		}

		if scanData.Scan.Id == 0 {
			return fmt.Errorf("scan not found"), http.StatusNotFound
		}

		scan = *scanData

		return nil, http.StatusOK

	}); err != nil {
		return scan, err
	}

	return scan, nil
}

func (h *ScanHandler) CreateNewScan(newScanData *models.ScanCreate) error {

	if err := utils.ValidateStruct(newScanData); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "RepositoryId":
				return fmt.Errorf("repository ID is required")
			case "UserId":
				return fmt.Errorf("user ID is required")
			case "Result":
				return fmt.Errorf("scan result is required")
			case "Vulnerabilities":
				return fmt.Errorf("vulnerabilities count is required")
			case "Status":
				return fmt.Errorf("scan status is required")
			}
		}
	}

	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		// check if the repo id is valid
		checkRepoId, err := h.RepoRepo.FindByParam(tx, newScanData.UserId, newScanData.RepositoryId, "", "")
		if err != nil {
			return fmt.Errorf("failed to find repository: %w", err), http.StatusInternalServerError

		}

		if checkRepoId.Id == 0 {
			return fmt.Errorf("repository with id '%d' in your user data not found", newScanData.RepositoryId), http.StatusNotFound
		}

		newScanData.ScanTime = time.Now().Format("2006-01-02 15:04:05")

		// create new data
		if err := h.ScanRepo.Create(tx, newScanData); err != nil {
			return fmt.Errorf("failed to create scan: %w", err), http.StatusInternalServerError
		}

		return nil, http.StatusOK

	}); err != nil {
		return err
	}

	return nil
}

func (h *ScanHandler) DeleteScan(repoId, scanId int) error {

	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		// check if data scan available
		scanData, err := h.ScanRepo.FindById(tx, repoId, scanId)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		if scanData.Scan.Id == 0 {
			return fmt.Errorf("scan not found, delete failed"), http.StatusNotFound
		}

		// delete repo
		if err := h.ScanRepo.Delete(tx, scanData.Scan.RepositoryId, scanData.Scan.Id); err != nil {
			return fmt.Errorf("failed to delete scan: %w", err), http.StatusInternalServerError
		}

		return nil, http.StatusOK

	}); err != nil {
		return err
	}

	return nil
}
