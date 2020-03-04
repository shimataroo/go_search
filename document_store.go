package gosearch

import (
	"database/sql"
	"log"
)

// DocumentStore ドキュメントストア構造体
type DocumentStore struct {
	db *sql.DB
}

// NewDocumentStore ドキュメントストア作成
func NewDocumentStore(db *sql.DB) *DocumentStore {
	return &DocumentStore{db: db}
}

func (ds *DocumentStore) save(title string) (DocumentID, error) {
	query := "INSERT INTO documents (document_title) VALUES (?)"
	result, err := ds.db.Exec(query, title)
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	return DocumentID(id), err
}
