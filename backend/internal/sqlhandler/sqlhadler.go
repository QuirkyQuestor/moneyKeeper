package sqlhandler

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	// _ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	PGErrUniqueViolation = "unique_violation"
)

var (
	SQLConflict = errors.New("the record not unique")
)

type DBConfig struct {
	DbType     string `yaml:"dbType"`
	DbProtocol string `yaml:"dbProtocol"`
	DbSSLMode  string `yaml:"dbSSLMode"`
	DbHost     string `yaml:"dbHost"`
	DBPort     uint32 `yaml:"dbPort"`
	DbName     string `yaml:"dbName"`
	DbUser     string `yaml:"userName"`
	DbPassword string `yaml:"userPassword"`
}

// Reads DBConfig
func (c *DBConfig) GetConfig() {

	yamlFile, err := os.ReadFile("dbconfig.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

// dbConnect opens a connection to the database
func DBConnect() *sql.DB {

	var c DBConfig
	c.GetConfig()

	// db, err := sql.Open(c.DbType, c.DbUser+":"+c.DbPassword+"@"+c.DbProtocol+"("+c.DbHost+":"+fmt.Sprint(c.DBPort)+")/"+c.DbName)

	// connStr := "postgres://postgres:password@localhost/DB_1?sslmode=disable"
	// db, err = sql.Open("postgres", connStr)

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.DbHost, c.DBPort, c.DbUser, c.DbPassword, c.DbName, c.DbSSLMode)

	// open database
	db, err := sql.Open(c.DbType, psqlconn)

	if err != nil {
		log.WithError(err).Error("Could not open connection to the DB")
		panic(err)
	}
	log.Info("Connection to the DB established")

	return db
}
