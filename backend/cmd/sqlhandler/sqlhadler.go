package sqlhandler

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	DbType     string `yaml:"dbType"`
	DbProtocol string `yaml:"dbProtocol"`
	DbHost     string `yaml:"dbHost"`
	DBPort     uint32 `yaml:"dbPort"`
	DbName     string `yaml:"dbName"`
	DbUser     string `yaml:"userName"`
	DbPassword string `yaml:"userPassword"`
}

// Reads DBConfig
func (c *DBConfig) GetConfig() {

	yamlFile, err := ioutil.ReadFile("dbconfig.yaml")
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

	db, err := sql.Open(c.DbType, c.DbUser+":"+c.DbPassword+"@"+c.DbProtocol+"("+c.DbHost+":"+fmt.Sprint(c.DBPort)+")/"+c.DbName)
	if err != nil {
		log.WithError(err).Error("Could not open connection to the DB")
		panic(err)
	}
	return db
}
