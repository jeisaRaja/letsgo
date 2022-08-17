package mysql

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"jeisaRaja.git/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name,email,hashed_password,created) VALUES(?,?,?, UTC_TIMESTAMP()) `
	_, err = u.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return models.ErrDuplicateEmail
			}
		}
	}
	return nil
}
func (u *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashed_password []byte
	row := u.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashed_password)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashed_password, []byte(password))
	fmt.Println("ini errornya:", err)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (u *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
