package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/momokii/go-codecheck/backend/core/models"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Find(tx *sql.Tx, userId int, paginationFilter models.PaginationFiltering) (*[]models.Repository, int, error) {

	var repositories []models.Repository

	offset := (paginationFilter.Page - 1) * paginationFilter.PerPage
	total := 0
	paramData := []interface{}{}

	total_query := "SELECT COUNT(*) FROM repositories WHERE 1=1"
	query := "SELECT id, user_id, name, description, path, created_at, updated_at FROM repositories WHERE 1=1"

	if userId <= 0 {
		return nil, total, fmt.Errorf("invalid user ID") // No user ID provided
	} else {
		add_query := " AND user_id = ?"
		total_query += add_query
		query += add_query
		paramData = append(paramData, userId)
	}

	if paginationFilter.Search != "" {
		search := "%" + paginationFilter.Search + "%"
		add_query := " AND (name LIKE ? OR description LIKE ? OR path LIKE ?)"
		total_query += add_query
		query += add_query
		paramData = append(paramData, search, search, search)
	}

	// count total data
	if err := tx.QueryRow(total_query, paramData...).Scan(&total); err != nil {
		return nil, total, err
	}

	// ordering data
	if paginationFilter.SortBy == models.SortByDesc {
		query += " ORDER BY id DESC"
	} else {
		query += " ORDER BY id ASC"
	}
	query += " LIMIT ? OFFSET ?"
	paramData = append(paramData, paginationFilter.PerPage, offset)

	rows, err := tx.Query(query, paramData...)
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()

	for rows.Next() {
		var repo models.Repository
		if err := rows.Scan(&repo.Id, &repo.UserId, &repo.Name, &repo.Description, &repo.Path, &repo.CreateAt, &repo.UpdateAt); err != nil {
			return nil, total, err
		}
		repositories = append(repositories, repo)
	}

	return &repositories, total, nil
}

func (r *Repository) FindByParam(tx *sql.Tx, user_id, repo_id int, name_repo, path string) (*models.Repository, error) {

	var repo models.Repository

	query := "SELECT id, user_id, name, description, path, created_at, updated_at FROM repositories WHERE 1=1 AND user_id = ?"
	base_param := []interface{}{user_id}

	if name_repo != "" || path != "" {
		query += " AND ("

		var conditions []string
		if name_repo != "" {
			conditions = append(conditions, "name = ?")
			base_param = append(base_param, name_repo)
		}

		if path != "" {
			if len(conditions) > 0 {
				query += strings.Join(conditions, " OR ") + " OR path = ?"
			} else {
				query += "path = ?"
			}
			base_param = append(base_param, path)
		} else if len(conditions) > 0 {
			query += strings.Join(conditions, " OR ")
		}

		query += ")"
	}

	if repo_id != 0 {
		if name_repo != "" || path != "" {
			query += " OR (user_id = ? AND id = ?)"
			base_param = append(base_param, user_id, repo_id)
		} else {
			query += " AND id = ?"
			base_param = append(base_param, repo_id)
		}
	}

	row := tx.QueryRow(query, base_param...)
	if err := row.Scan(&repo.Id, &repo.UserId, &repo.Name, &repo.Description, &repo.Path, &repo.CreateAt, &repo.UpdateAt); err != nil {
		if err == sql.ErrNoRows {
			return &repo, nil // No repository found with the given ID
		}
		return nil, err // An error occurred while scanning the row
	}

	return &repo, nil
}

func (r *Repository) Create(tx *sql.Tx, repo *models.RepositoryCreate) (*models.Repository, error) {

	query := "INSERT INTO repositories (user_id, name, description, path) VALUES (?, ?, ?, ?)"

	row, err := tx.Exec(query, repo.UserId, repo.Name, repo.Description, repo.Path)
	if err != nil {
		return nil, err
	}

	id_new, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Repository{
		Id: int(id_new),
	}, nil
}

func (r *Repository) Update(tx *sql.Tx, repo *models.RepositoryUpdate) error {

	query := "UPDATE repositories SET name = ?, description = ?, path = ? WHERE id = ?"
	_, err := tx.Exec(query, repo.Name, repo.Description, repo.Path, repo.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(tx *sql.Tx, id, user_id int) error {

	_, err := tx.Exec("DELETE FROM scans WHERE repository_id = ?", id, user_id)
	if err != nil {
		return err
	}

	query := "DELETE FROM repositories WHERE id = ? AND user_id = ?"
	_, err = tx.Exec(query, id, user_id)
	if err != nil {
		return err
	}

	return nil
}
