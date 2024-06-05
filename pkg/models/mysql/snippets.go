package mysql

import (
	"database/sql"
	"errors"

	"aitunews.kz/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(user_id int, content, created string) (int, error) {
	stmt := `INSERT INTO snippets (user_id, content, created)
	VALUES($1, $2, $3) RETURNING id`
	var id int
	err := m.DB.QueryRow(stmt, user_id, content, created).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, user_id, content, created FROM snippets WHERE id = $1`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}
	err := row.Scan(&s.ID, &s.User_id, &s.Content, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, user_id, content, created FROM snippets`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.User_id, &s.Content, &s.Created)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

func (m *SnippetModel) Latestt(name_for int, filter string) ([]*models.Snippet, error) {
	stmt := `SELECT id, user_id, content, created FROM snippets WHERE user_id = $1`
	rows, err := m.DB.Query(stmt, name_for)

	if filter == "all" {
		stmt = `SELECT id, user_id, content, created FROM snippets`
		rows, err = m.DB.Query(stmt)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.User_id, &s.Content, &s.Created)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

func (m *SnippetModel) Delete(id int) error {
	stmt := `DELETE FROM snippets WHERE id = $1`
	result, err := m.DB.Exec(stmt, id)

	if err != nil || result == nil {
		return err
	}
	return nil
}
