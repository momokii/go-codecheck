package scan

import (
	"database/sql"

	"github.com/momokii/go-codecheck/backend/core/models"
)

type ScanRepository struct{}

func NewScanRepository() *ScanRepository {
	return &ScanRepository{}
}

func (r *ScanRepository) Find(tx *sql.Tx, repoId int, paginationFilter models.PaginationFiltering) (*[]models.ScanFull, int, error) {

	var scans []models.ScanFull

	offset := (paginationFilter.Page - 1) * paginationFilter.PerPage
	total := 0
	paramData := []interface{}{}

	total_query := "SELECT COUNT(s.id) FROM scans s LEFT JOIN repositories r ON s.repository_id = r.id WHERE 1=1"
	query := `
		SELECT 
			s.id, 
			s.repository_id, 
			COALESCE(r.name, '') as name,
			COALESCE(r.description, '') as description,
			COALESCE(r.path, '') as path,
			s.scan_time, 
			s.result, 
			s.vulnerabilities,
			s.status,
			s.created_at, 
			s.updated_at 
		FROM 
			scans s 
		LEFT JOIN 
			repositories r 
		ON 
			s.repository_id = r.id 
		WHERE 
			1 = 1
	`

	if repoId > 0 {
		addQuery := " AND s.repository_id = ?"
		total_query += addQuery
		query += addQuery
		paramData = append(paramData, repoId)
	}

	if paginationFilter.Search != "" {
		search := "%" + paginationFilter.Search + "%"
		addQuery := " AND (r.name LIKE ? OR r.description LIKE ? OR r.path LIKE ?)"
		total_query += addQuery
		query += addQuery
		paramData = append(paramData, search, search)
	}

	// count data
	if err := tx.QueryRow(total_query, paramData...).Scan(&total); err != nil {
		return nil, total, err
	}

	// ordering data
	if paginationFilter.SortBy == models.SortByDesc {
		query += " ORDER BY s.id DESC"
	} else {
		query += " ORDER BY s.id ASC"
	}
	query += " LIMIT ? OFFSET ?"
	paramData = append(paramData, paginationFilter.PerPage, offset)

	rows, err := tx.Query(query, paramData...)
	if err != nil {
		return &scans, total, nil
	}
	defer rows.Close()

	for rows.Next() {
		var scan models.ScanFull
		if err := rows.Scan(&scan.Scan.Id, &scan.Scan.RepositoryId, &scan.RepositoryName, &scan.RepositoryDescription, &scan.RepositoryPath, &scan.Scan.ScanTime, &scan.Scan.Result, &scan.Scan.Vulnerabilities, &scan.Scan.Status, &scan.Scan.CreatedAt, &scan.Scan.UpdatedAt); err != nil {
			return &scans, total, err
		}
		scans = append(scans, scan)
	}

	return &scans, total, nil
}

func (r *ScanRepository) FindById(tx *sql.Tx, repoId, id int) (*models.ScanFull, error) {

	var scan models.ScanFull
	query := `
		SELECT 
			s.id, 
			s.repository_id, 
			r.name,
			r.description,
			r.path,
			s.scan_time, 
			s.result, 
			s.vulnerabilities,
			s.status,
			s.created_at, 
			s.updated_at 
		FROM 
			scans s 
		LEFT JOIN 
			repositories r 
		ON 
			s.repository_id = r.id 
		WHERE 
			s.id = ?
			AND s.repository_id = ?
	`
	err := tx.QueryRow(query, id, repoId).Scan(&scan.Scan.Id, &scan.Scan.RepositoryId, &scan.RepositoryName, &scan.RepositoryDescription, &scan.RepositoryPath, &scan.Scan.ScanTime, &scan.Scan.Result, &scan.Scan.Vulnerabilities, &scan.Scan.Status, &scan.Scan.CreatedAt, &scan.Scan.UpdatedAt)
	if err != nil {

		if err == sql.ErrNoRows {
			return &scan, nil
		}

		return nil, err
	}

	return &scan, nil
}

func (r *ScanRepository) Create(tx *sql.Tx, scan *models.ScanCreate) error {

	query := "INSERT INTO scans (repository_id, scan_time, result, vulnerabilities, status) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.Exec(query, scan.RepositoryId, scan.ScanTime, scan.Result, scan.Vulnerabilities, scan.Status)
	if err != nil {
		return err
	}

	return nil
}

func (r *ScanRepository) Delete(tx *sql.Tx, repoId, id int) error {

	query := "DELETE FROM scans WHERE repository_id = ? AND id = ?"
	_, err := tx.Exec(query, repoId, id)
	if err != nil {
		return err
	}

	return nil
}
