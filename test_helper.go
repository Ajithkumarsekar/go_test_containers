package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

type pgContainer struct {
	container  testcontainers.Container
	connString string
}

var singleInstanceContainer *pgContainer

func getPgContainerInstance() *pgContainer {
	if singleInstanceContainer == nil {
		fmt.Println("Creating a new container")
		singleInstanceContainer = startPgContainer()
	}
	fmt.Println("returning container instance")
	return singleInstanceContainer
}

func startPgContainer() *pgContainer {
	ctx := context.Background()

	dbCredentials := struct {
		user     string
		password string
		dbName   string
	}{
		user:     "postgres",
		password: "password",
		dbName:   "postgres",
	}

	var env = map[string]string{
		"POSTGRES_PASSWORD": dbCredentials.password,
		"POSTGRES_USER":     dbCredentials.user,
		"POSTGRES_DB":       dbCredentials.dbName,
	}

	var port = "5432/tcp"
	dbURL := func(port nat.Port) string {
		return fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", port.Port(), dbCredentials.dbName)
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{port},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
			Env:          env,
			WaitingFor:   wait.ForSQL(nat.Port(port), "postgres", dbURL).Timeout(time.Second * 5),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		fmt.Println(err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		fmt.Println(err)
	}

	dbConnStr := fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", mappedPort.Port(), dbCredentials.dbName)
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		fmt.Println(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.ExecContext(ctx, "CREATE SCHEMA IF NOT EXISTS netstat")
	if err != nil {
		fmt.Println("error while creating new schema")
	}

	return &pgContainer{
		container:  container,
		connString: dbConnStr,
	}
}
