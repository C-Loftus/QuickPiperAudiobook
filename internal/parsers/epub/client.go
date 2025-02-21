package epub

import (
	"fmt"
	"io"
)

type EpubSplitter struct {
	filepath string
	book     *Book
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

func (p *EpubSplitter) GetCoverImage() (io.Reader, error) {
	for _, manifestItem := range p.book.Opf.Manifest {
		if manifestItem.Properties == "cover-image" {
			file, err := p.book.OpenInternalBookFile(manifestItem.Href)
			if err != nil {
				return nil, err
			}
			return file, nil
		}
	}
	return nil, fmt.Errorf("no cover image found in the epub")
}

type SectionData struct {
	Filename string
	Text     io.Reader
}

// Split a book into individual io.Readers for each chapter
// which can be used to process the book in parallel
func (p *EpubSplitter) SplitBySection() ([]SectionData, error) {
	spineItemsInOrder := p.book.Opf.Spine.Items
	idToFile := make(map[string]string)

	for _, manifestItem := range p.book.Opf.Manifest {
		idToFile[manifestItem.ID] = manifestItem.Href
	}

	var sections []SectionData
	for _, item := range spineItemsInOrder {
		filepath := idToFile[item.IDref]
		reader, err := p.book.OpenInternalBookFile(filepath)
		if err != nil {
			return nil, err
		}
		sections = append(sections, SectionData{
			Filename: filepath,
			Text:     reader,
		})
	}
	return sections, nil
}

// Return a list of filenames for each section in the book
// While this is maximally semantically useful for an audiobook
// given the fact that toc.ncx is used for talking books, it can
// point to individual parts of the inner XML and thus is harder to
// parse of make use of. It is here mainly for completeness
func (p *EpubSplitter) GetSectionNamesViaToc() ([]string, error) {

	var sections []string
	for _, navPoint := range p.book.Ncx.NavPoints {
		sections = append(sections, navPoint.NavLabel.Content.Src)
	}
	if len(sections) == 0 {
		return nil, fmt.Errorf("no sections found in the epub")
	}

	return sections, nil
}

func (p *EpubSplitter) Close() {
	p.book.Close()
}
