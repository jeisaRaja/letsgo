package mysql

import (
	"database/sql"

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
	return 0, nil
}
func (u *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
