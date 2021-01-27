package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type conf struct {
	DbType     string `yaml:"dbType"`
	DbProtocol string `yaml:"dbProtocol"`
	DbHost     string `yaml:"dbHost"`
	DBPort     uint32 `yaml:"dbPort"`
	DbName     string `yaml:"dbName"`
	DbUser     string `yaml:"userName"`
	DbPassword string `yaml:"userPassword"`
}

// Reads DBConfig
func (c *conf) GetConfig() {

	yamlFile, err := ioutil.ReadFile("dbconfig.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func dbConnect() (db *sql.DB) {

	var c conf
	c.GetConfig()

	db, err := sql.Open(c.DbType, c.DbUser+":"+c.DbPassword+"@"+c.DbProtocol+"("+c.DbHost+":"+fmt.Sprint(c.DBPort)+")/"+c.DbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
