// Copyright 2025 Colton Loftus
// SPDX-License-Identifier: AGPL-3.0-only

package epub

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseEnglish(t *testing.T) {
	englishFiles := []string{"dubliners_epub2.epub", "dubliners_epub3.epub"}

	t.Run("Split using opf", func(t *testing.T) {
		for _, f := range englishFiles {
			client, err := NewEpubSplitter(filepath.Join("testdata", f))
			require.NoError(t, err)
			readers, err := client.SplitBySection()
			require.NoError(t, err)
			const sectionsInDublinersEbook = 18
			require.Len(t, readers, sectionsInDublinersEbook)
			cover, err := client.GetCoverImage()
			if strings.HasSuffix(f, "epub3.epub") {
				require.NoError(t, err)
				require.NotNil(t, cover)
			} else {
				require.Error(t, err)
				require.Nil(t, cover)
			}

		}
	})

	t.Run("Split using toc", func(t *testing.T) {
		for _, f := range englishFiles {
			client, err := NewEpubSplitter(filepath.Join("testdata", f))
			require.NoError(t, err)
			_, err = client.GetSectionNamesViaToc()
			require.NoError(t, err)
		}
	})

}

func TestParseChinese(t *testing.T) {

	chineseFiles := []string{"呻吟語_epub3.epub", "呻吟語_epub2.epub"}

	for _, f := range chineseFiles {
		client, err := NewEpubSplitter(filepath.Join("testdata", f))
		require.NoError(t, err)
		_, err = client.SplitBySection()
		require.NoError(t, err)
	}

}
