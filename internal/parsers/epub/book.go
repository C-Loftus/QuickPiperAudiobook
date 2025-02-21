package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path"
)

// A top level epub book
type Book struct {
	Ncx       Ncx       `json:"ncx"`
	Opf       Opf       `json:"opf"`
	Container Container `json:"-"`
	Mimetype  string    `json:"-"`

	zipReader *zip.ReadCloser
}

// Open an epub file and return a parsed book representation
func Open(fn string) (*Book, error) {
	zipReader, err := zip.OpenReader(fn)
	if err != nil {
		return nil, err
	}

	book := Book{zipReader: zipReader}
	mimetype, err := book.readBytes("mimetype")

	if err == nil {
		book.Mimetype = string(mimetype)
		err = book.readXML("META-INF/container.xml", &book.Container)
	}
	if err == nil {
		err = book.readXML(book.Container.Rootfile.Path, &book.Opf)
	}

	for _, mf := range book.Opf.Manifest {
		if mf.ID == book.Opf.Spine.Toc {
			err = book.readXML(book.relativeFilename(mf.Href), &book.Ncx)
			break
		}
	}

	if err != nil {
		zipReader.Close()
		return nil, err
	}

	return &book, nil
}

// OpenInternalBookFile open resource file
func (p *Book) OpenInternalBookFile(n string) (io.ReadCloser, error) {
	return p.openInternal(p.relativeFilename(n))
}

// Files list resource files
func (p *Book) Files() []string {
	var fns []string
	for _, f := range p.zipReader.File {
		fns = append(fns, f.Name)
	}
	return fns
}

// Close close file reader
func (p *Book) Close() {
	p.zipReader.Close()
}

func (p *Book) relativeFilename(n string) string {
	return path.Join(path.Dir(p.Container.Rootfile.Path), n)
}

func (p *Book) readXML(n string, v interface{}) error {
	fd, err := p.openInternal(n)
	if err != nil {
		return nil
	}
	defer fd.Close()
	dec := xml.NewDecoder(fd)
	return dec.Decode(v)
}

func (p *Book) readBytes(n string) ([]byte, error) {
	fd, err := p.openInternal(n)
	if err != nil {
		return nil, nil
	}
	defer fd.Close()

	return io.ReadAll(fd)
}

func (p *Book) openInternal(n string) (io.ReadCloser, error) {
	for _, f := range p.zipReader.File {
		if f.Name == n {
			return f.Open()
		}
	}
	return nil, fmt.Errorf("file %s does not exist", n)
}
