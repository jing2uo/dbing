package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	dbEngine := flag.String("e", "postgresql", "Your database engine: postgresql or mysql")
	dbHost := flag.String("h", "localhost", "Your database host")
	dbPort := flag.String("p", "", "Your database port")
	dbUser := flag.String("u", "postgres", "Your database user")
	dbPassword := flag.String("k", "alauda", "Your database password")
	dbName := flag.String("d", "", "Database name you want")
	flag.Parse()

	if strings.EqualFold(string(*dbPort), "") && strings.EqualFold(string(*dbEngine), "postgresql") {
		*dbPort = "5432"
	}

	if strings.EqualFold(string(*dbPort), "") && strings.EqualFold(string(*dbEngine), "mysql") {
		*dbPort = "3306"
	}

	if strings.EqualFold(string(*dbEngine), "postgresql") {
		sqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", *dbHost, *dbPort, *dbUser, *dbPassword)
		var err error
		db, err = sql.Open("postgres", sqlconn)
		checkErr(err)
		err = db.Ping()
		checkErr(err)
		println("Succese connect to database.")
		_, err = db.Exec("DROP DATABASE IF EXISTS alaudatestdb;")
		_, err = db.Exec("CREATE DATABASE alaudatestdb;")
		checkErr(err)
		println("Success write to database.")
		if string(*dbName) != "" {
			_, err = db.Exec("CREATE DATABASE " + *dbName)
			checkErr(err)
			println("Success create database: ", *dbName)
		}
	} else if strings.EqualFold(string(*dbEngine), "mysql") {
		sqlconn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", *dbUser, *dbPassword, *dbHost, *dbPort)
		println(sqlconn)
		var err error
		db, err = sql.Open("mysql", sqlconn)
		checkErr(err)
		err = db.Ping()
		checkErr(err)
		println("Succese connect to database.")
		_, err = db.Exec("DROP DATABASE IF EXISTS alaudatestdb;")
		_, err = db.Exec("CREATE DATABASE IF NOT EXISTS alaudatestdb;")
		checkErr(err)
		println("Success write to database.")
		if string(*dbName) != "" {
			_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + *dbName)
			checkErr(err)
			println("Success create database: ", *dbName)
		}
	} else {
		println("Sorry, -e only support postgresql or mysql.")
		os.Exit(1)
	}
}
