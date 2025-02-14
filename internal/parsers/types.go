package parsers

import "io"

// RawTextSplitter is an interface for splitting a text into chapters
// It should be implemented by any parser that wants to be used
// to eventually generate an audiobook with chapters in an mp3 file
type TextSplitter interface {
	SplitByChapter() ([]io.Reader, error)
}
