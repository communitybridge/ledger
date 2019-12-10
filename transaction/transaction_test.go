package transaction

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"gopkg.in/testfixtures.v2"
)

var (
	fixtures *testfixtures.Context
)

func TestMain(m *testing.M) {

	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	fixtures, err = testfixtures.NewFolder(conn, &testfixtures.PostgreSQL{}, "../testdata/fixtures")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}
