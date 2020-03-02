package gosearch

import (
	"container/list"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// DocumentID ドキュメント番号
type DocumentID int64

// Index 転置インデックス
type Index struct {
	Dictionary     map[string]PostingsList
	TotalDocsCount int
}

// Posting is structポスティング
type Posting struct {
	DocID         DocumentID
	Positions     []int
	TermFrequency int
}

// PostingsList ポスティングリスト
// list
type PostingsList struct {
	*list.List
}

// NewIndex インデックス作成
func NewIndex() *Index {
	dict := make(map[string]PostingsList)
	return &Index{
		Dictionary:     dict,
		TotalDocsCount: 0,
	}
}

// NewPosting ポスティング作成
func NewPosting(docID DocumentID, positions ...int) *Posting {
	return &Posting{
		docID,
		positions,
		len(positions),
	}
}

// NewPostingsList ポスティングリストの作成
// 双方向性リストにより実装
func NewPostingsList(postings ...*Posting) PostingsList {
	l := list.New()
	for _, posting := range postings {
		l.PushBack(posting)
	}
	return PostingsList{l}
}

func (pl PostingsList) add(p *Posting) {
	pl.PushBack(p)
}

func (pl PostingsList) last() *Posting {
	e := pl.List.Back()
	if e == nil {
		return nil
	}
	return e.Value.(*Posting)
}

// Add ポスティング追加メソッド
func (pl PostingsList) Add(new *Posting) {
	last := pl.last()
	// 最後尾がない or 追加のドキュメントIDと追加が違ったら新規に追加
	if last == nil || last.DocID != new.DocID {
		pl.add(new)
		return
	}
	last.Positions = append(last.Positions, new.Positions...)
	// 出現回数のインクリメント
	last.TermFrequency++
}

// インデックス表示メソッド
func (idx Index) String() string {
	var padding int
	keys := make([]string, 0, len(idx.Dictionary))
	for k := range idx.Dictionary {
		l := utf8.RuneCountInString(k)
		if padding < l {
			padding = l
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	strs := make([]string, len(keys))
	format := " [%-" + strconv.Itoa(padding) + "s] -> %s"
	for i, k := range keys {
		if postingList, ok := idx.Dictionary[k]; ok {
			strs[i] = fmt.Sprintf(format, k, postingList.String())
		}
	}
	return fmt.Sprintf("total documents : %v\ndictionary:\n%v\n",
		idx.TotalDocsCount, strings.Join(strs, "\n"))
}

func (pl PostingsList) String() string {
	str := make([]string, 0, pl.Len())
	for e := pl.Front(); e != nil; e = e.Next() {
		str = append(str, e.Value.(*Posting).String())
	}
	return strings.Join(str, "=>")
}

func (p Posting) String() string {
	return fmt.Sprintf("(%v,%v,%v)",
		p.DocID, p.TermFrequency, p.Positions)
}
