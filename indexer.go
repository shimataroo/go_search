package gosearch

import (
	"bufio"
	"io"
)

// Indexer インデクサー
type Indexer struct {
	index     *Index
	tokenizer *Tokenizer
}

// NewIndexer インデクサーの作成
func NewIndexer(tokenizer *Tokenizer) *Indexer {
	return &Indexer{
		index:     NewIndex(),
		tokenizer: tokenizer,
	}
}

func (idxr *Indexer) update(docID DocumentID, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(idxr.tokenizer.SplitFunc)
	var position int

	for scanner.Scan() {
		term := scanner.Text()
		if postingList, ok := idxr.index.Dictionary[term]; !ok {
			idxr.index.Dictionary[term] =
				NewPostingsList(NewPosting(docID, position))
		} else {
			postingList.Add(NewPosting(docID, position))
		}
		position++
	}
	idxr.index.TotalDocsCount++
}
