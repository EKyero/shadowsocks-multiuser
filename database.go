package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Database struct
type Database struct {
	Connection *sql.DB
	DBHost     string
	DBPort     int
	DBUser     string
	DBPass     string
	DBName     string
}

// User struct
type User struct {
	Upload         uint64
	Download       uint64
	Port           int
	Method         string
	Password       string
	Enable         int
	TransferEnable uint64
}

// Open database connection
func (database *Database) Open() error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", database.DBUser, database.DBPass, database.DBHost, database.DBPort, database.DBName))
	if err != nil {
		return err
	}

	database.Connection = db
	return nil
}

// Close database connection
func (database *Database) Close() error {
	if database.Connection != nil {
		return database.Connection.Close()
	}

	return nil
}

// GetUser RT.
func (database *Database) GetUser() ([]User, error) {
	results, err := database.Connection.Query("SELECT u, d, port, method, passwd, enable, transfer_enable FROM user WHERE enable=1")
	if err != nil {
		return nil, err
	}

	users := make([]User, 65535)
	count := 0
	for results.Next() {
		var user User

		err = results.Scan(&user.Upload, &user.Download, &user.Port, &user.Method, &user.Password, &user.Enable, &user.TransferEnable)
		if err != nil {
			return nil, err
		}

		users[count] = user
		count++
	}

	return users[:count], nil
}

// UpdateBandwidth RT.
func (database *Database) UpdateBandwidth(port int, upload, download uint64) error {
	log.Printf("Reporting %d uploaded %d downloaded %d to database", port, upload, download)

	results, err := database.Connection.Query("SELECT u, d FROM user")
	if err != nil {
		return err
	}

	var cloudUpload uint64
	var cloudDownload uint64

	if results.Next() {
		err = results.Scan(&cloudUpload, &cloudDownload)
		if err != nil {
			return err
		}
	}

	cloudUpload += upload
	cloudDownload += download

	_, err = database.Connection.Query(fmt.Sprintf("UPDATE user SET u=%d, d=%d, t=%d WHERE port=%d", cloudUpload, cloudDownload, time.Now().Unix(), port))
	return err
}

func newDatabase(host string, port int, user, pass, name string) *Database {
	database := Database{}
	database.DBHost = host
	database.DBPort = port
	database.DBUser = user
	database.DBPass = pass
	database.DBName = name

	return &database
}
