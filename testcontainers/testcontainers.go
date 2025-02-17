package testcontainers

// import (
// 	"context"

// 	"github.com/testcontainers/testcontainers-go"
// )

// type quickPiperAudiobookContainer struct {
// 	testcontainers.Container
// }

// func makeContainer(dockerfile string, command []string, files []testcontainers.ContainerFile) (*quickPiperAudiobookContainer, error) {
// 	req := testcontainers.ContainerRequest{
// 		FromDockerfile: testcontainers.FromDockerfile{
// 			Context:    "../",
// 			Dockerfile: dockerfile,
// 		},
// 		Files: files,
// 		Cmd:   command,
// 	}
// 	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
// 		ContainerRequest: req,
// 		Started:          true,
// 	})
// 	return &quickPiperAudiobookContainer{Container: container}, err
// }
