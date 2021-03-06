package storage

import (
	"context"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

// StartTestDB runs or starts a local mongodb container if one is not already present.
func StartTestDB() Client {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Printf("TestDB: Error creating docker client %v", err)
	}

	hostBinding := nat.PortBinding{
		HostIP:   "127.0.0.1",
		HostPort: "27017",
	}
	containerPort, err := nat.NewPort("tcp", "27017")
	if err != nil {
		log.Printf("TestDB: Error creating docker port %v", err)
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "mongo",
			Env: []string{
				"MONGO_INITDB_ROOT_USERNAME=admin",
				"MONGO_INITDB_ROOT_PASSWORD=password",
			},
		},
		&container.HostConfig{
			PortBindings: portBinding,
		}, nil, "mongo")
	if err != nil {
		log.Printf("TestDB: Error creating docker for mongodb %v", err)
	}

	cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})

	os.Setenv("DATABASE_CONNECTION_STRING", "mongodb://admin:password@127.0.0.1:27017/")

	databaseName := uuid.New().String()
	os.Setenv("DATABASE_NAME", databaseName)

	client := new(databaseName)

	log.Printf("TestDB: Seed Database Name %v", databaseName)

	return client
}
