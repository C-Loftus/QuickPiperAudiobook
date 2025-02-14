package epub

import (
	"QuickPiperAudiobook/internal/parsers"
	"io"
)

type EpubSplitter struct {
	filepath string
	book     *Book
	parsers.TextSplitter
}

// Create a client for splitting an epub book into individual chapters
func NewEpubSplitter(filepath string) (*EpubSplitter, error) {

	book, err := Open(filepath)
	if err != nil {
		return nil, err
	}

	return &EpubSplitter{
		filepath: filepath,
		book:     book,
	}, nil
}

// Split a book into individual io.Readers for each chapter
// which can be used to process the book in parallel
func (p *EpubSplitter) SplitByChapter() ([]io.Reader, error) {
	spineItemsInOrder := p.book.Opf.Spine.Items
	idToFile := make(map[string]string)

	for _, manifestItem := range p.book.Opf.Manifest {
		idToFile[manifestItem.ID] = manifestItem.Href
	}

	var readers []io.Reader
	for _, item := range spineItemsInOrder {
		filepath := idToFile[item.IDref]
		reader, err := p.book.OpenInternalBookFile(filepath)
		defer p.book.Close()
		if err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}
	return readers, nil
}
