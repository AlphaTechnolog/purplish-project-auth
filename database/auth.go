package database

import (
	"database/sql"
	"fmt"
	// "strings"

	"github.com/alphatechnolog/purplish-auth/api"
	"github.com/alphatechnolog/purplish-auth/lib"
	"github.com/google/uuid"
)

// TODO: Extract these to envvars.
// this will be the value for company_id when the user registers without a company
const GUEST_STRING = "<GUEST>"
const MEMBERSHIP_MICRO_BASE = "http://localhost:8006/memberships"

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	LocalScopes string `json:"-"`
	CompanyID   string `json:"-"`
	Password    string `json:"-"`
}

type CompanyMembership struct {
	CompanyID    string `json:"company_id"`
	MembershipID string `json:"membership_id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	Scopes       string `json:"scopes"`
}

func (u *User) IsGuest() bool {
	return u.CompanyID == GUEST_STRING
}

func (u *User) ResolveScopes() (string, error) {
	companyMembership, err := api.GetCompanyMembership(u.CompanyID)
	if err != nil {
		return "", err
	}

	scopes := lib.ExpandScopes(companyMembership.Scopes + " " + u.LocalScopes)

	return scopes, nil

	// companyScopes := lib.ExpandScopes(companyMembership.Scopes)
	// userScopes := lib.ExpandScopes(u.LocalScopes)
	// scopesArray := strings.Split(companyScopes+" "+userScopes, " ")
	// uniqueScopes := make(map[string]struct{})
	// for _, scope := range scopesArray {
	// 	if scope != "" {
	// 		uniqueScopes[scope] = struct{}{}
	// 	}
	// }

	// var scopes string = ""
	// for scope := range uniqueScopes {
	// 	if scopes == "" {
	// 		scopes = scope
	// 	} else {
	// 		scopes += " " + scope
	// 	}
	// }

	// return scopes, nil

	// resp, err := http.Get(fmt.Sprintf(
	// 	"%s/company-membership/%s",
	// 	MEMBERSHIP_MICRO_BASE,
	// 	u.CompanyID,
	// ))

	// if err != nil {
	// 	return "", fmt.Errorf("Unable to tryna obtain company membership: %w", err)
	// }

	// defer resp.Body.Close()

	// responseBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", fmt.Errorf("Cannot read company membership stream: %w", err)
	// }

	// type companyMembershipResponse struct {
	// 	companyMembership CompanyMembership `json:"company_membership"`
	// }

	// var companyMembershipRes companyMembershipResponse
	// if err := json.Unmarshal(responseBytes, &companyMembershipRes); err != nil {
	// 	return "", fmt.Errorf("Unable to decode company membership data: %w", err)
	// }

	// marshalized, _ := json.Marshal(companyMembershipRes)
	// fmt.Println("responseString:\n", string(responseBytes))
	// fmt.Println("marshalized response:\n", string(marshalized))
	// fmt.Println("parsed response:\n", companyMembershipRes)

	// companyMembership := companyMembershipRes.companyMembership
	// scopesArray := strings.Split(companyMembership.Scopes+" "+u.LocalScopes, " ")
	// uniqueScopes := make(map[string]struct{})
	// for _, scope := range scopesArray {
	// 	if scope != "" {
	// 		uniqueScopes[scope] = struct{}{}
	// 	}
	// }

	// var scopes string = ""
	// for scope := range uniqueScopes {
	// 	if scopes == "" {
	// 		scopes = scope
	// 		continue
	// 	}

	// 	scopes += " " + scope
	// }

	// return scopes, nil
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
