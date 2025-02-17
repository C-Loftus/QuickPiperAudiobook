package testcontainers

// import (
// 	"context"
// 	"path/filepath"
// 	"testing"

// 	"github.com/stretchr/testify/require"
// 	"github.com/testcontainers/testcontainers-go"
// )

// func TestUbuntuE2E(t *testing.T) {
// 	files := []testcontainers.ContainerFile{
// 		{
// 			HostFilePath:      filepath.Join("testdata", "titlepage_and_2_chapters.epub"),
// 			ContainerFilePath: "/app/examples/titlepage_and_2_chapters.epub",
// 		},
// 	}
// 	container, err := makeContainer("Dockerfile", []string{"--help"}, files)
// 	require.NoError(t, err)
// 	defer container.Terminate(context.Background())

// 	exitCode, reader, err := container.Exec(context.Background(), []string{"QuickPiperAudiobook", "/app/examples/titlepage_and_2_chapters.epub"})
// 	require.NoError(t, err)
// 	require.Equal(t, 0, exitCode)
// 	require.NotNil(t, reader)
// }
