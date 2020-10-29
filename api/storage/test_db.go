package storage

import (
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

// StartTestDB runs or starts a local mongodb container if one is not already present.
func StartTestDB() Client {
	cmd := exec.Command("docker", "run", "-d", "-p", "27017:27017", "-e", "MONGO_INITDB_ROOT_USERNAME=admin", "-e", "MONGO_INITDB_ROOT_PASSWORD=password", "--name", "db", "mongo")
	err := cmd.Run()
	if err != nil {
		log.Printf("TestDB: Error running docker image for mongodb. Attempting to start %v", err)
		cmd := exec.Command("docker", "start", "db")
		err := cmd.Run()
		if err != nil {
			log.Printf("TestDB: Error starting docker image for mongodb %v", err)
		}
	}

	os.Setenv("DATABASE_CONNECTION_STRING", "mongodb://admin:password@localhost:27017/")

	databaseName := uuid.New().String()
	os.Setenv("DATABASE_NAME", databaseName)

	client := new(databaseName)

	log.Printf("TestDB: Seed Database Name %v", databaseName)

	return client
}
