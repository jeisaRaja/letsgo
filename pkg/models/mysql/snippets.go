package mysql

import (
	"database/sql"
	"log"

	"jeisaRaja.git/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id,title,content,created,expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id=?`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	stmt2 := "SELECT id,title,content,created,expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10"
	test, err := m.DB.Query(stmt2)
	if err != nil {
		return nil, nil
	}
	defer test.Close()
	out := []*models.Snippet{}

	for test.Next() {
		s2 := &models.Snippet{}
		err := test.Scan(&s2.ID, &s2.Title, &s2.Content, &s2.Created, &s2.Expires)
		if err != nil {
			return nil, nil
		}
		out = append(out, s2)
	}
	log.Print(out)
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmts := "SELECT id,title,content,created,expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10"
	rows, err := m.DB.Query(stmts)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()
	s := []*models.Snippet{}

	for rows.Next() {
		tmp := &models.Snippet{}
		err = rows.Scan(&tmp.ID, &tmp.Title, &tmp.Content, &tmp.Created, &tmp.Expires)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		s = append(s, tmp)
	}
	log.Print(s)
	return s, nil
}
