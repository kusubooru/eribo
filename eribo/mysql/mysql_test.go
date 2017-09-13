package mysql

import (
	"flag"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	user   = flag.String("user", "eribo", "username to use to run tests on MySQL")
	pass   = flag.String("pass", "", "password to use to run tests on MySQL")
	host   = flag.String("host", "localhost", "host for connecting to MySQL and run the tests")
	port   = flag.String("port", "3306", "port for connecting to MySQL and run the tests")
	dbname = flag.String("dbname", "eribo_test", "test database to use to run the tests")
)

func init() {
	flag.Parse()
}

func setup(t *testing.T) *EriboStore {
	if *pass == "" {
		t.Errorf("No password provided for user %q to connect to MySQL and run the tests.", *user)
		t.Errorf("These tests need a MySQL account %q that has access to test database %q.", *user, *dbname)
		t.Fatal("Use: go test -tags=db -pass '<db password>'")
	}
	datasource := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", *user, *pass, *host, *port, *dbname)
	s, err := NewEriboStore(datasource)
	if err != nil {
		t.Fatalf("NewEriboStore failed for datasource %q: %v", datasource, err)
	}
	return s
}

func teardown(t *testing.T, s *EriboStore) {
	if err := s.dropSchema(); err != nil {
		t.Error(err)
	}
}
