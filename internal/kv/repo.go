package kv

import (
	"bufio"
	"fmt"
	"github.com/oxio/kv/internal/fileop"
	"github.com/oxio/kv/internal/parser"
)

type Repo interface {
	Find(key string, defaultValue *string) (*parser.Item, error)
	Get(key string) (*parser.Item, error)
	Set(item *parser.Item) error
}

var _ Repo = &RepoImpl{}

type RepoImpl struct {
	adapter fileop.FileAdapter
	parser  *parser.LineParser
}

func NewKvRepo(filePath string) *RepoImpl {
	return &RepoImpl{
		adapter: fileop.NewFileAdapter(filePath),
		parser:  parser.NewLineParser(),
	}
}

func (r *RepoImpl) Find(key string, defaultValue *string) (*parser.Item, error) {
	if "" == key {
		return nil, fmt.Errorf("key is empty")
	}

	var collection = parser.NewItemCollection()
	err := r.adapter.ReadByLine(r.makeReader(collection))
	if err != nil {
		return nil, err
	}

	for _, item := range *collection.Items {
		if item.Key == key {
			return item, nil
		}
	}

	return parser.NewItem(key, *defaultValue)
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
		if item.Key == key {
			return item, nil
		}
	}

	return nil, fmt.Errorf("key doesn't exist")
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
		for k, item := range *collection.Items {
			if item.Key == incoming.Key {
				(*collection.Items)[k] = incoming
				return nil
			}
		}
		*collection.Items = append(*collection.Items, incoming)
		return nil
	}
}

func (r *RepoImpl) makeWriter(collection *parser.ItemCollection) fileop.WriterFunc {
	return func(writer *bufio.Writer) error {
		for _, item := range *collection.Items {
			var _, err = writer.WriteString(item.ToLine())
			if err != nil {
				return err
			}
		}
		return nil
	}
}
