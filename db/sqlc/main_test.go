package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/WanCodeBase/GinModule/util"

	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	conf, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatalln("load config failed:", err)
		return
	}
	testDB, err = sql.Open(conf.DBDriver, conf.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
		return
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
