package gosearch

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var testDB *sql.DB

func setup() *sql.DB {
	db, err := sql.Open("mysql", "root:@/gosearch")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(` TRUNCATE TABLE documents`)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.RemoveAll("_index_data"); err != nil {
		log.Fatal(err)
	}
	if err := os.Mkdir("_index_data", 0777); err != nil {
		log.Fatal(err)
	}
	return db
}

func TestMain(m *testing.M) {
	testDB = setup()
	defer testDB.Close()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreateIndex(t *testing.T) {
	engine := NewSearchEngine(testDB)

	type testDoc struct {
		title string
		body  string
	}
	docs := []testDoc{
		{"test1", "Do you Shimataro?"},
	}
	for _, doc := range docs {
		r := strings.NewReader(doc.body)
		if err := engine.AddDocument(doc.title, r); err != nil {
			t.Fatalf("faild to add doc %s : %v", doc.title, err)
		}
	}
	if err := engine.Flush(); err != nil {
		t.Fatalf("failed to save index to file : %v", err)
	}

	type testCase struct {
		file       string
		postingStr string
	}
	testCases := []testCase{
		{"_index_data/_0.dc", "1"},
		{
			"_index_data/do",
			` [{"DocID":1, "Positions":[0], "TermFrequency":1}]`,
		},
		{
			"_index_data/shimataro",
			` [{"DocID":1, "Positions":[2], "TermFrequency":1}]`,
		},
		{
			"_index_data/you",
			` [{"DocID":1, "Positions":[1], "TermFrequency":1}]`,
		},
	}

	for _, testCase := range testCases {
		func() {
			file, err := os.Open(testCase.file)
			if err != nil {
				t.Fatalf("failed to load index : %v", err)
			}
			defer file.Close()
			b, err := ioutil.ReadAll(file)
			if err != nil {
				t.Fatalf("failed to load index : %v", err)
			}
			got := string(b)
			var buf bytes.Buffer
			_ = json.Compact(&buf, []byte(testCase.postingStr))
			want := buf.String()
			if got != want {
				t.Errorf("got : %v\nwant : %v", got, want)
			}
		}()
	}

}
