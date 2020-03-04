package gosearch

import (
	"reflect"
	"strings"
	"testing"
)

func TestUpdate(t *testing.T) {
	collection := []string{
		"Are you Shima?",
		"Hello World!",
		"My name is Shima!",
	}
	indexer := NewIndexer(NewTokenizer())

	for i, doc := range collection {
		indexer.update(DocumentID(i), strings.NewReader(doc))
	}

	actual := indexer.index
	expected := &Index{
		Dictionary: map[string]PostingsList{
			"are": NewPostingsList(
				NewPosting(0, 0)),
			"hello": NewPostingsList(
				NewPosting(1, 0)),
			"is": NewPostingsList(
				NewPosting(2, 2)),
			"my": NewPostingsList(
				NewPosting(2, 0)),
			"name": NewPostingsList(
				NewPosting(2, 1)),
			"shima": NewPostingsList(
				NewPosting(0, 2),
				NewPosting(2, 3)),
			"world": NewPostingsList(
				NewPosting(1, 1)),
			"you": NewPostingsList(
				NewPosting(0, 1)),
		},
		TotalDocsCount: 3,
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("index error\n")
	}
}
