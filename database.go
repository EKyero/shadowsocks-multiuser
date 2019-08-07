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
	NodeID     int
	NodeRate   float64
}

// User struct
type User struct {
	ID             int
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

	rate, err := database.GetRate()
	if err != nil {
		return err
	}

	database.NodeRate = rate
	return nil
}

// Close database connection
func (database *Database) Close() error {
	if database.Connection != nil {
		return database.Connection.Close()
	}

	return nil
}

// GetRate RT.
func (database *Database) GetRate() (float64, error) {
	results, err := database.Connection.Query(fmt.Sprintf("SELECT traffic_rate FROM ss_node WHERE id=%d", database.NodeID))
	if err != nil {
		return -1, err
	}

	rate := float64(-1)
	if results.Next() {
		err = results.Scan(&rate)
		if err != nil {
			return -1, err
		}
	}

	return rate, nil
}

// GetUser RT.
func (database *Database) GetUser() ([]User, error) {
	results, err := database.Connection.Query("SELECT id, u, d, port, method, passwd, enable, transfer_enable FROM user WHERE enable=1")
	if err != nil {
		return nil, err
	}

	users := make([]User, 65535)
	count := 0
	for results.Next() {
		var user User

		err = results.Scan(&user.ID, &user.Upload, &user.Download, &user.Port, &user.Method, &user.Password, &user.Enable, &user.TransferEnable)
		if err != nil {
			return nil, err
		}

		users[count] = user
		count++
	}

	return users[:count], nil
}

// UpdateBandwidth RT.
func (database *Database) UpdateBandwidth(instance *Instance) error {
	log.Printf("Reporting %d uploaded %d downloaded %d to database", instance.UserID, instance.Bandwidth.Upload, instance.Bandwidth.Download)

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

	userUpload := uint64(float64(instance.Bandwidth.Upload) * database.NodeRate)
	userDownload := uint64(float64(instance.Bandwidth.Download) * database.NodeRate)

	_, err = database.Connection.Query(fmt.Sprintf("INSERT INTO user_traffic_log (user_id, u, d, node_id, rate, traffic, log_time) VALUES (%d, %d, %d, %d, %f, %d, %d)", instance.UserID, userUpload, userDownload, database.NodeID, database.NodeRate, userUpload+userDownload, time.Now().Unix()))
	if err != nil {
		return err
	}

	cloudUpload += userUpload
	cloudDownload += userDownload
	_, err = database.Connection.Query(fmt.Sprintf("UPDATE user SET u=%d, d=%d, t=%d WHERE id=%d", cloudUpload, cloudDownload, time.Now().Unix(), instance.UserID))
	return err
}

func newDatabase(host string, port int, user, pass, name string, id int) *Database {
	database := Database{}
	database.DBHost = host
	database.DBPort = port
	database.DBUser = user
	database.DBPass = pass
	database.DBName = name
	database.NodeID = id

	return &database
}
