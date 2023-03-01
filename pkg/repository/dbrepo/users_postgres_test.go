package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host       = "localhost"
	user       = "postgres"
	password   = "postgres"
	dbName     = "users_test"
	port       = "5435"
	dsn        = `host%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5`
	authMethod = "trust"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB

func TestMain(m *testing.M) {
	//connect to docker
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker.")
	}
	pool = p

	//setup our docker options, specify image and so forth
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_HOST_AUTH_METHOD=" + authMethod,
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	//get a resource (docker image)
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("Could not start resource.")

	}
	//start the image and wait until it's ready

	retryDSN := fmt.Sprintf(dsn, host, port, user, password, dbName)
	fmt.Println(retryDSN)

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", retryDSN)
		if err != nil {
			log.Println("Error: ", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		fmt.Println(err)
		_ = pool.Purge(resource)
		log.Fatal("Could not connect to database.", err)
	}

	//populate the database with empty tables
	err = createTables()
	if err != nil {
		log.Fatalf("Error creating tables")
	}

	//run tests
	code := m.Run()

	//clean up
	if err := pool.Purge(resource); err != nil {
		log.Fatal("Could not purge  resource ", err)
	}

	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/users.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("Can't ping db.")
	}
}
