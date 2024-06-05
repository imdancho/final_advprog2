package mysql

import (
	"database/sql"
	"errors" // New import
	"fmt"
	"strconv"

	// New import
	"aitunews.kz/snippetbox/pkg/models" // New import
	"golang.org/x/crypto/bcrypt"        // New import
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(full_name, email, phone, password, role string, c_code int, confirmation bool) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (full_name, email, phone, hashed_password, created, role, c_code, confirmation)
	VALUES($1, $2, $3, $4, NOW(), $5, $6, $7)`

	_, err = m.DB.Exec(stmt, full_name, email, phone, string(hashedPassword), role, c_code, confirmation)
	if err != nil {
		// Handle PostgreSQL error codes and messages here if needed
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, string, error) {
	// Retrieve the id and hashed password associated with the given email. If no
	// matching email exists, or the user is not active, we return the
	// ErrInvalidCredentials error.
	var id int
	var role string
	var hashedPassword []byte
	stmt := "SELECT id, role, hashed_password FROM users WHERE email = $1 AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &role, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", models.ErrInvalidCredentials
		} else {
			return 0, "", err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, "", models.ErrInvalidCredentials
		} else {
			return 0, "", err
		}
	}

	// stmt = fmt.Sprintf(`SELECT c_code FROM services WHERE id= %s`, "")
	stmt = fmt.Sprintf(`UPDATE users SET c_code = 0, confirmation = true WHERE id = %s`, "11")
	result, err := m.DB.Exec(stmt)

	if err != nil || result == nil {
		return id, role, nil
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, role, nil
}

// We'll use the Get method to fetch details for a specific user based
// on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}

func (m *UserModel) Confirm(email, c_code_inputed string) error {
	var id int
	var c_code string
	stmt := "SELECT id, c_code FROM users WHERE email = $1"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &c_code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}

	if c_code_inputed != c_code {
		return errors.New("not correct confirmation code")
	}

	stmt = fmt.Sprintf(`UPDATE users SET c_code = 0, confirmation = true WHERE id = %s`, strconv.Itoa(id))
	result, err := m.DB.Exec(stmt)

	if err != nil || result == nil {
		return nil
	}
	return nil
}

func (m *UserModel) ConfirmToo(email, c_code_inputed string) error {
	var id int
	var c_code string
	stmt := "SELECT id, c_code FROM users WHERE email = $1"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &c_code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}

	if c_code_inputed != c_code {
		return errors.New("not correct confirmation code")
	}

	stmt = fmt.Sprintf(`UPDATE users SET c_code = 0, confirmation = true WHERE id = %s`, strconv.Itoa(id))
	result, err := m.DB.Exec(stmt)

	if err != nil || result == nil {
		return nil
	}
	return nil
}

func (m *UserModel) ConfirmCheck(email string) bool {
	var confirmation bool
	stmt := "SELECT confirmation FROM users WHERE email = $1"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&confirmation)

	if err == nil {
		return confirmation
	}

	// log.Println(confirmation)
	return confirmation

}

func (m *UserModel) CodeCheck(email string) string {
	var c_code string
	stmt := "SELECT c_code FROM users WHERE email = $1"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&c_code)

	if err == nil {
		return c_code
	}

	// log.Println(confirmation)
	return c_code
}

func (m *UserModel) AddCode(email, c_code string) error {
	stmt := `UPDATE users SET c_code = $1 WHERE email = $2`
	result, err := m.DB.Exec(stmt, c_code, email)

	if err != nil || result == nil {
		return nil
	}
	return nil
}
