package kvf

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/oxio/kvf/internal/fileop"
	"github.com/oxio/kvf/internal/parser"
)

type Repo interface {
	Get(key string) (*parser.Item, error)
	Set(item *parser.Item) error
}

var _ Repo = &RepoImpl{}

var (
	ErrItemNotFound = errors.New("item not found")
)

type RepoImpl struct {
	adapter fileop.FileAdapter
	parser  *parser.LineParser
}

func NewRepo(filePath string, noErrorOnInaccessibleFile bool) *RepoImpl {
	return &RepoImpl{
		adapter: fileop.NewFileAdapter(filePath, noErrorOnInaccessibleFile),
		parser:  parser.NewLineParser(),
	}
}

func (r *RepoImpl) FindAll() (*[]*parser.Item, error) {
	var collection = parser.NewItemCollection()
	err := r.adapter.ReadByLine(r.makeReader(collection))
	if err != nil {
		return nil, err
	}
	return collection.Items, nil
}

func (r *RepoImpl) Get(key string) (*parser.Item, error) {
	if "" == key {
		return nil, fmt.Errorf("key is empty")
	}

	var collection = parser.NewItemCollection()
	err := r.adapter.ReadByLine(r.makeReader(collection))
	if err != nil {
		return nil, err
	}

	for _, item := range *collection.Items {
		if item.IsEmpty || item.IsComment {
			continue
		}
		if item.Key == key {
			return item, nil
		}
	}

	return nil, ErrItemNotFound
}

func (r *RepoImpl) Set(item *parser.Item) error {
	var collection = parser.NewItemCollection()
	read := r.makeReader(collection)
	update := r.makeUpdater(collection, item)
	write := r.makeWriter(collection)

	return r.adapter.EnsureUpdate(read, update, write)
}

func (r *RepoImpl) makeReader(collection *parser.ItemCollection) fileop.ReaderFunc {
	return func(line string) error {
		item, err := r.parser.Parse(line)
		if err != nil {
			return err
		}
		*collection.Items = append(*collection.Items, item)
		return nil
	}
}

func (r *RepoImpl) makeUpdater(collection *parser.ItemCollection, incoming *parser.Item) fileop.UpdateFunc {
	return func() error {
		found := false
		for k, item := range *collection.Items {
			if item.Key == incoming.Key {
				item.Val = incoming.Val
				(*collection.Items)[k] = item
				found = true
				break
			}
		}
		if !found {
			*collection.Items = append(*collection.Items, incoming)
		}
		return nil
	}
}

func (r *RepoImpl) makeWriter(collection *parser.ItemCollection) fileop.WriterFunc {
	return func(writer *bufio.Writer) (bytesWritten int64, err error) {
		for _, item := range *collection.Items {
			var n, err = writer.WriteString(item.ToLine())
			bytesWritten += int64(n)
			if err != nil {
				return bytesWritten, err
			}
		}
		return bytesWritten, nil
	}
}
