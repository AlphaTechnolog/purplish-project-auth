package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alphatechnolog/purplish-auth/database"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TODO: This should be set using some envvar method.
const TOKEN_SECRET = "secret"
const BCRYPT_COST = 14

type UserLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserLoginPayload struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func decryptUser(tokenString string) (*database.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(TOKEN_SECRET), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &database.User{
			ID:          claims["user_id"].(string),
			Name:        claims["name"].(string),
			Surname:     claims["surname"].(string),
			Email:       claims["email"].(string),
			LocalScopes: claims["local_scopes"].(string),
			CompanyID:   claims["company_id"].(string),
			Password:    "",
		}, nil
	}

	return nil, errors.New("Unable to decode user")
}

func getUserScopes(d *sql.DB, c *gin.Context) error {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return errors.New("Invalid authorization string")
	}

	tokenString, err := func() (string, error) {
		parts := strings.Split(authorization, " ")
		if len(parts) <= 1 {
			return "", errors.New("Malformed authorization string")
		}

		return parts[1], nil
	}()

	if err != nil {
		return err
	}

	user, err := decryptUser(tokenString)
	if err != nil {
		return err
	}

	scopes, err := user.ResolveScopes()
	if err != nil {
		return fmt.Errorf("Unable to resolve user scopes: %w", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"user":   user,
		"scopes": scopes,
	})

	return nil
}

func loginUser(d *sql.DB, c *gin.Context) error {
	bodyContents, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	var loginPayload UserLoginPayload
	if err = json.Unmarshal(bodyContents, &loginPayload); err != nil {
		return err
	}

	userMatch, err := database.GetUserByEmail(d, loginPayload.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either email or password is incorrect"})
		return nil
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(userMatch.Password),
		[]byte(loginPayload.Password),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either email or password is incorrect"})
		return nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      userMatch.ID,
		"company_id":   userMatch.CompanyID,
		"email":        userMatch.Email,
		"name":         userMatch.Name,
		"surname":      userMatch.Surname,
		"local_scopes": userMatch.LocalScopes,
	})

	tokenString, err := token.SignedString([]byte(TOKEN_SECRET))

	if err != nil {
		return fmt.Errorf("Unable to create user token: %w", err)
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

	return nil
}

func createUser(d *sql.DB, c *gin.Context) error {
	bodyContents, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	defer c.Request.Body.Close()

	var createUserPayload CreateUserLoginPayload
	if err = json.Unmarshal(bodyContents, &createUserPayload); err != nil {
        return fmt.Errorf("Unexpected JSON input: %w", err)
	}

    if _, err := database.GetUserByEmail(d, createUserPayload.Email); err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already used"})
        return nil
    }

    bytes, err := bcrypt.GenerateFromPassword([]byte(createUserPayload.Password), BCRYPT_COST)
    encryptedPassword := string(bytes)

    err = database.CreateUser(d, database.CreateUserPayload{
        Name: createUserPayload.Name,
        Surname: createUserPayload.Surname,
        Email: createUserPayload.Email,
        HashedPassword: encryptedPassword,
        CompanyID: database.GUEST_STRING,
    })

    if err != nil {
        return fmt.Errorf("Unable to create user: %w", err)
    }

    c.JSON(http.StatusCreated, gin.H{"ok": true})

	return nil
}

func CreateAuthRoutes(d *sql.DB, r *gin.RouterGroup) {
	r.GET("/user-scopes/", WrapError(WithDB(d, getUserScopes)))
	r.POST("/login/", WrapError(WithDB(d, loginUser)))
	r.POST("/register/", WrapError(WithDB(d, createUser)))
}
