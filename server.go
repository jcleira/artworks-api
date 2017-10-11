package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	yaml "gopkg.in/yaml.v2"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jcleira/artworks-api/artworks"
)

// config is the configuration struct, it contains all the settings needed to
// run the server.
//
// TODO Important!. This development struct should be removed, should we use
// a config management tool?.
type config struct {
	Development struct {
		Datasource string `yaml:"datasource"`
	}
	Preproduction struct {
		Datasource string `yaml:"datasource"`
	}
}

// getConfiguration reads all the necesary configurations for the server to run.
//
// Returns a configuration struct.
func getConfiguration() *config {
	data, err := ioutil.ReadFile("dbconfig.yml")
	if err != nil {
		log.Fatal(err)
	}
	var config config

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}

// configureRoutes will configure all the REST API routes, it returns a *mux.Router
// with all the core api routes configured.
func configureRoutes(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	artworks.ConfigureHandlers(r, db)

	return r
}

// main would initialize and run the http server.
func main() {
	environment := flag.String("environment", "development", "Running environment")
	flag.Parse()

	config := getConfiguration()

	datasource := config.Development.Datasource

	switch *environment {
	case "preproduction":
		datasource = config.Preproduction.Datasource
		break

	}

	db, err := sql.Open("mysql", datasource)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.ListenAndServe(":3000", configureRoutes(db))
}
