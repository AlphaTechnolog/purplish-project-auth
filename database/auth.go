package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// this will be the value for company_id when the user registers without a company
const GUEST_STRING = "<GUEST>"

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	LocalScopes string `json:"-"`
	CompanyID   string `json:"-"`
	Password    string `json:"-"`
}

func (u *User) IsGuest() bool {
	return u.CompanyID == GUEST_STRING
}

type CreateUserPayload struct {
	Name           string
	Surname        string
	Email          string
	HashedPassword string
	CompanyID      string
}

func (cu *CreateUserPayload) IsGuest() bool {
	return cu.CompanyID == GUEST_STRING
}

func GetUsers(d *sql.DB) ([]User, error) {
	var users []User

	sql := `
        SELECT u.id, u.name, u.surname, u.email, u.local_scopes, u.company_id, u.password
        FROM users u;
    `

	rows, err := d.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.LocalScopes,
			&user.CompanyID,
			&user.Password,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(d *sql.DB, ID string) (User, error) {
	var user User

	sql := `
        SELECT u.id, u.name, u.surname, u.email, u.local_scopes, u.company_id, u.password
        FROM users u
        WHERE u.id = ?;
    `

	row := d.QueryRow(sql, ID)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.LocalScopes,
		&user.CompanyID,
		&user.Password,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUserByEmail(d *sql.DB, email string) (User, error) {
	var user User

	sql := `
        SELECT u.id, u.name, u.surname, u.email, u.local_scopes, u.company_id, u.password
        FROM users u
        WHERE u.email = ?;
    `

	row := d.QueryRow(sql, email)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.LocalScopes,
		&user.CompanyID,
		&user.Password,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func CreateUser(d *sql.DB, createPayload CreateUserPayload) error {
	sql := `
        INSERT INTO users (id, name, surname, email, password, local_scopes, company_id)
        VALUES
            (?, ?, ?, ?, ?, ?, ?);
    `

	localScopes := ""

	// special:guest will help identify this user as an user which has no
	// company associated (yet).
	if createPayload.IsGuest() {
		localScopes = "special:guest *:items"
	}

	_, err := d.Exec(
		sql,
		uuid.New().String(),
		createPayload.Name,
		createPayload.Surname,
		createPayload.Email,
		createPayload.HashedPassword,
		localScopes,
		createPayload.CompanyID, // if user is guest, this **should** be sent like `GUEST_STRING`.
	)

	if err != nil {
		return fmt.Errorf("Unable to register user %s %s: %w", createPayload.Name, createPayload.Surname, err)
	}

	return nil
}
