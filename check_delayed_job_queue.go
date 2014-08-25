// Sensu/Nagios check to monitor the delayed_job queue depth
//
// Copyright 2014 David Radcliffe <radcliffe.david@gmail.com>

package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
)

var exitCode int = 3
var count int

func codeToString() string {
	switch exitCode {
	case 0:
		return "OK"
	case 1:
		return "WARNING"
	case 2:
		return "CRITICAL"
	}
	return "UNKNOWN"
}

func done() {
	final(strconv.Itoa(count) + " queued jobs")
}

func final(outString string) {
	fmt.Println(codeToString() + ": " + outString)
	os.Exit(exitCode)
}

func main() {

	// Setup vars
	var dbname string
	var dbuser string
	var dbpassword string
	var dbhost string
	var dbport string
	var debug bool
	var warning int
	var critical int

	// Parse command line options
	flag.StringVar(&dbname, "dbname", "", "database name")
	flag.StringVar(&dbuser, "dbuser", "", "database username")
	flag.StringVar(&dbpassword, "dbpassword", "", "database password")
	flag.StringVar(&dbhost, "dbhost", "", "database host")
	flag.StringVar(&dbport, "dbport", "3306", "database port")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.IntVar(&warning, "warning", 5, "warning threshold")
	flag.IntVar(&critical, "critical", 3, "critical threshold")
	flag.Parse()

	// Print debug info
	if debug == true {
		fmt.Println("==============")
		fmt.Println("DEBUG INFO")
		fmt.Println("dbname:", dbname)
		fmt.Println("dbuser:", dbuser)
		fmt.Println("dbpassword:", dbpassword)
		fmt.Println("dbhost:", dbhost)
		fmt.Println("dbport:", dbport)
		fmt.Println("warning:", warning)
		fmt.Println("critical:", critical)
		fmt.Println("==============")
	}

	// Run Query
	if dbpassword != "" {
		dbpassword = ":" + dbpassword
	}
	if dbhost != "" {
		dbhost = "@tcp(" + dbhost + ":" + dbport + ")"
	}
	dsn := dbuser + dbpassword + dbhost + "/" + dbname
	if debug == true {
		fmt.Println("dsn: ", dsn)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		final(err.Error())
	}
	defer db.Close()

	err = db.QueryRow("SELECT COUNT(*) as count FROM delayed_jobs WHERE run_at < NOW()").Scan(&count)
	db.Close()
	if err != nil {
		final(err.Error())
	}

	db.Close()

	// Output logic
	if count > critical {
		exitCode = 2
		done()
	}

	if count > warning {
		exitCode = 1
		done()
	}

	exitCode = 0
	done()

}
