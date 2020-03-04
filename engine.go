package gosearch

import (
	"database/sql"
	"io"
	"os"
	"path/filepath"
)

// Engine 検索エンジン
type Engine struct {
	tokenizer     *Tokenizer
	indexer       *Indexer
	documentStore *DocumentStore
	indexDir      string
}

// NewSearchEngine 検索エンジン作成
func NewSearchEngine(db *sql.DB) *Engine {
	tokenizer := NewTokenizer()
	indexer := NewIndexer(tokenizer)
	documentStore := NewDocumentStore(db)

	path, ok := os.LookupEnv("INDEX_DIR_PATH")
	if !ok {
		current, _ := os.Getwd()
		path = filepath.Join(current, "_index_data")
	}
	return &Engine{
		tokenizer:     tokenizer,
		indexer:       indexer,
		documentStore: documentStore,
		indexDir:      path,
	}
}

// AddDocument インデックスの構築
func (e *Engine) AddDocument(title string, reader io.Reader) error {
	id, err := e.documentStore.save(title)
	if err != nil {
		return err
	}
	e.indexer.update(id, reader)
	return nil
}

// Flush ファイル書き出し
func (e *Engine) Flush() error {
	writer := NewIndexWriter(e.indexDir)
	return writer.Flush(e.indexer.index)
}
