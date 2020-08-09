package mysql

import (
	"database/sql"

	"github.com/eze8789/snippetbox/pkg/models"
)

// SnippetModel wraps the db connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert function create snippets
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 1, err
	}
	return int(id), nil
}

// Get retrieve an specific snippet
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// get pointer to sql row of the snippet
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a snippet struct with value 0
	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Latest return the latest 10 snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	// Get all the records matching the stmt query
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create snippets slice to store records
	snippets := []*models.Snippet{}

	// Iterate through all the retrieved rows and save it into snippet
	for rows.Next() {
		// Declare a pointer and save parsed records on it
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	// Manage possible errors in the previous iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
