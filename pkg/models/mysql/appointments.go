package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"aitunews.kz/snippetbox/pkg/models"
)

type AppointmentModel struct {
	DB *sql.DB
}

func (m *AppointmentModel) Insert(user_id int, service_id, time string) (int, error) {
	stmt := `INSERT INTO appointments (user_id, service_id, time) VALUES($1, $2, $3) RETURNING id`
	var id int
	err := m.DB.QueryRow(stmt, user_id, service_id, time).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *AppointmentModel) Update(id int, time string) error {
	stmt := `UPDATE appointments SET time = $1 WHERE id = $2`
	_, err := m.DB.Exec(stmt, time, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *AppointmentModel) Delete(id int) error {
	stmt := `DELETE FROM appointments WHERE id = $1`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *AppointmentModel) Get(id int) (*models.Appointment, error) {
	stmt := `SELECT id, user_id, service_id, time FROM appointments WHERE id = $1`

	row := m.DB.QueryRow(stmt, id)
	s := &models.Appointment{}
	err := row.Scan(&s.ID, &s.User_id, &s.Service_id, &s.Time)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

// func (m *AppointmentModel) Latest(name_for string) ([]*models.Appointment, error) {
// 	stmt := `SELECT * FROM appointments`
// 	if name_for == "" {
// 		stmt = `SELECT * FROM appointments`
// 	}
// 	rows, err := m.DB.Query(stmt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	appointments := []*models.Appointment{}
// 	for rows.Next() {
// 		s := &models.Appointment{}
// 		err = rows.Scan(&s.ID, &s.User_id, &s.Service_id, &s.Time)
// 		if err != nil {
// 			return nil, err
// 		}
// 		appointments = append(appointments, s)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return appointments, nil
// }

func (m *AppointmentModel) Latest(name_for, limit, offset int) ([]*models.Appointment, error) {
	var stmt string
	if name_for == 1 {
		stmt = `SELECT id, user_id, service_id, time FROM appointments`
	} else {
		stmt = `SELECT id, user_id, service_id, time FROM appointments WHERE user_id = $1`
	}

	stmt += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	var rows *sql.Rows
	var err error
	if name_for == 5 {
		rows, err = m.DB.Query(stmt)
	} else {
		rows, err = m.DB.Query(stmt, name_for)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appointments := []*models.Appointment{}
	for rows.Next() {
		s := &models.Appointment{}
		err = rows.Scan(&s.ID, &s.User_id, &s.Service_id, &s.Time)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}
