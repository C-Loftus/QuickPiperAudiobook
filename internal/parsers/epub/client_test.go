package epub

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseEnglish(t *testing.T) {

	englishFiles := []string{"dubliners_epub2.epub", "dubliners_epub3.epub"}

	for _, f := range englishFiles {
		client, err := NewEpubSplitter(filepath.Join("testdata", f))
		require.NoError(t, err)
		readers, err := client.SplitByChapter()
		require.NoError(t, err)
		const sectionsInDublinersEbook = 18
		require.Len(t, readers, sectionsInDublinersEbook)
	}
}

func TestParse漢語(t *testing.T) {

	chineseFiles := []string{"呻吟語_epub3.epub", "呻吟語_epub2.epub"}

	for _, f := range chineseFiles {
		client, err := NewEpubSplitter(filepath.Join("testdata", f))
		require.NoError(t, err)
		_, err = client.SplitByChapter()
		require.NoError(t, err)
	}

}
