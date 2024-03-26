//go:build integration

package dbrepo

import (
	"24-Testing-Simple-Web-App/pkg/data"
	"24-Testing-Simple-Web-App/pkg/repository"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "users_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB
var testRepo repository.DataBaseRepo

func TestMain(m *testing.M) {
	// connect to docker; fail if docker not running
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	pool = p

	// set up our docker options, specifying the image and so forth
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// get a resource (docker image)
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// start the image and wait until it's ready
	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database: %s", err)
	}

	// populate the database with empty tables
	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	testRepo = &PostgresDBRepo{DB: testDB}

	// run tests
	code := m.Run()

	// clean up
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
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
		t.Error("can't ping database")
	}
}

func TestPostgresDBRepo_InsertUser(t *testing.T) {
	testUser1 := data.User{
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
		Password:  "secret",
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	testUser2 := data.User{
		FirstName: "Jack",
		LastName:  "User",
		Email:     "jack@example.com",
		Password:  "secret",
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id1, err := testRepo.InsertUser(testUser1)
	id2, err := testRepo.InsertUser(testUser2)
	if err != nil {
		t.Errorf("insert user returned %s", err.Error())
	}

	if id1 != 1 || id2 != 2 {
		t.Errorf("inserted user returned wrong id expected 1 but got %d", id1)
	}
}

func TestPostgresDBRepo_AllUsers(t *testing.T) {
	users, err := testRepo.AllUsers()
	if err != nil {
		t.Errorf("%s: failed with error %s", "AllUsers", err)
	}

	if len(users) < 1 {
		t.Errorf("expected 1 but got %d", len(users))
	}
}
func TestPostgresDBRepo_GetUser(t *testing.T) {
	testUser, err := testRepo.GetUser(1)
	if err != nil {
		t.Errorf("user with id %d, not found, %s", 1, err)
	}

	if testUser.ID != 1 {
		t.Errorf("user with id 1 not found")
	}
}
func TestPostgresDBRepo_GetUserByEmail(t *testing.T) {
	testUser, err := testRepo.GetUserByEmail("admin@example.com")
	if err != nil {
		t.Errorf("error by GetUserByEmail")
	}

	if testUser.Email != "admin@example.com" {
		t.Errorf("test for GetUserByEmail failed %s", err)
	}
}
func TestPostgresDBRepo_UpdateUser(t *testing.T) {
	testUser, _ := testRepo.GetUser(2)

	testUser.FirstName = "ivan"
	testUser.Email = "ivan@ivan.com"

	err := testRepo.UpdateUser(*testUser)
	if err != nil {
		t.Errorf("test for UpdateUser failed %s", err)
	}

	testUser, _ = testRepo.GetUser(2)
	if testUser.FirstName != "ivan" || testUser.Email != "ivan@ivan.com" {
		t.Errorf("test failed for  UpdateUser %s", err)
	}
}

func TestPostgresDBRepo_DeleteUser(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "Delete user",
			args:    args{id: 2},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := testRepo.DeleteUser(tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			testUser, _ := testRepo.GetUser(tt.args.id)

			if testUser != nil {
				t.Errorf("DeleteUser() failed want %v, but got %v", nil, testUser)
			}
		})
	}
}

func TestPostgresDBRepo_ResetPassword(t *testing.T) {
	type args struct {
		id       int
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "ResetPassword",
			args: args{
				id:       1,
				password: "password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := testRepo.ResetPassword(tt.args.id, tt.args.password); !errors.Is(err, tt.wantErr) {
				t.Errorf("ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
			}

			user, _ := testRepo.GetUser(tt.args.id)
			matches, err := user.PasswordMatches(tt.args.password)
			if err != nil {
				t.Error(err)
			}
			if !matches {
				t.Errorf("paswword do not match %s", err)
			}
		})
	}
}
func TestPostgresDBRepo_InsertUserImage(t *testing.T) {
	image := data.UserImage{
		UserID:    1,
		FileName:  "test.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newId, err := testRepo.InsertUserImage(image)
	if err != nil {
		t.Errorf("inserting image failed %s", err)
	}

	if newId < 1 {
		t.Errorf("Test failed %s", err)
	}
}
