package mysql

import (
	"database/sql"
	"errors"

	"aitunews.kz/snippetbox/pkg/models"
)

type ServiceModel struct {
	DB *sql.DB
}

func (m *ServiceModel) Insert(title, content, master string, price int) (int, error) {
	stmt := `INSERT INTO services (title, content, master, price) VALUES($1, $2, $3, $4) RETURNING id`
	var id int
	err := m.DB.QueryRow(stmt, title, content, master, price).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *ServiceModel) Update(title string, price int) error {
	stmt := `UPDATE services SET price = $1 WHERE title = $2`
	_, err := m.DB.Exec(stmt, price, title)
	if err != nil {
		return err
	}
	return nil
}

func (m *ServiceModel) Delete(title string) error {
	stmt := `DELETE FROM services WHERE title = $1`
	_, err := m.DB.Exec(stmt, title)
	if err != nil {
		return err
	}
	return nil
}

func (m *ServiceModel) Get(id int) (*models.Service, error) {
	stmt := `SELECT id, title, content, master, price FROM services WHERE id = $1`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Service{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Master, &s.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

func (m *ServiceModel) Latest(name_for, sort, sort_type string) ([]*models.Service, error) {
	stmt := `SELECT id, title, content, master, price FROM services`
	if name_for != "" {
		stmt += " WHERE title = $1"
	}
	if sort != "" {
		stmt += " ORDER BY " + sort
		if sort_type != "" {
			stmt += " " + sort_type
		}
	}

	rows, err := m.DB.Query(stmt, name_for)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	services := []*models.Service{}
	for rows.Next() {
		s := &models.Service{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Master, &s.Price)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return services, nil
}
