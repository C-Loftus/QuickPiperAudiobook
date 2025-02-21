package epub

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	book, err := Open("testdata/dubliners_epub3.epub")
	require.NoError(t, err)
	require.Equal(t, book.Opf.Manifest[0].Properties, "cover-image")
	require.Equal(t, book.Opf.Metadata.Title, []string{"Dubliners"})
}
