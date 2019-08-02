package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/NetchX/shadowsocks-multiuser/core"
)

func purge(instanceList map[int]*Instance, users []User) {
	for _, instance := range instanceList {
		contains := false

		for _, v := range users {
			if instance.Port == v.Port {
				contains = v.TransferEnable > v.Upload+v.Download
				break
			}
		}

		if !contains && instance.Started {
			instance.Stop()
		}
	}
}

func report(instance *Instance, database *Database) {
	if instance.Bandwidth.Upload != 0 || instance.Bandwidth.Download != 0 {
		if instance.Bandwidth.Last-time.Now().Unix() < 10 {
			if err := database.UpdateBandwidth(instance.Port, instance.Bandwidth.Upload, instance.Bandwidth.Download); err == nil {
				instance.Bandwidth.Reset()
			} else {
				log.Println(err)
			}
		}
	}
}

func update(instance *Instance, method, password string) {
	if instance.Method != method || instance.Password != password {
		instance.Method = method
		instance.Password = password

		restart(instance)
	}

	if instance.Started && (!instance.TCPStarted || !instance.UDPStarted) {
		restart(instance)
	}
}

func restart(instance *Instance) {
	if instance.Started {
		instance.Stop()
	}

	instance.Start()
}

func main() {
	log.Println("Starting shadowsocks-multiuser")
	log.Println("Version: 1.0.0")

	var flags struct {
		ListCipher bool
		DBHost     string
		DBPort     int
		DBUser     string
		DBPass     string
		DBName     string
	}

	flag.BoolVar(&flags.ListCipher, "listcipher", false, "list cipher")
	flag.StringVar(&flags.DBHost, "dbhost", "localhost", "database host")
	flag.IntVar(&flags.DBPort, "dbport", 3306, "database port")
	flag.StringVar(&flags.DBUser, "dbuser", "root", "database user")
	flag.StringVar(&flags.DBPass, "dbpass", "123456", "database pass")
	flag.StringVar(&flags.DBName, "dbname", "sspanel", "database name")
	flag.Parse()

	if flags.ListCipher {
		for _, v := range core.ListCipher() {
			fmt.Println(v)
		}

		return
	}

	instanceList := make(map[int]*Instance, 65535)

	log.Println("Started")
	for {
		log.Println("Wait 10 second for sync users")
		time.Sleep(10 * time.Second)

		log.Println("Start syncing")

		log.Println("Opening database connection")
		database := newDatabase(flags.DBHost, flags.DBPort, flags.DBUser, flags.DBPass, flags.DBName)
		if err := database.Open(); err != nil {
			log.Println(err)
			continue
		}

		defer database.Close()

		log.Println("Get database users")
		users, err := database.GetUser()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("Purge server users")
		purge(instanceList, users)

		for _, user := range users {
			log.Println(user)
			if instance, ok := instanceList[user.Port]; ok {
				if user.TransferEnable > user.Upload+user.Download {
					update(instance, user.Method, user.Password)
				} else {
					if instance.Started {
						instance.Stop()
					}

					report(instance, database)
					delete(instanceList, user.Port)
				}
			} else if user.TransferEnable > user.Upload+user.Download {
				log.Printf("Starting new instance for %d", user.Port)
				instance := newInstance(user.Port, user.Method, user.Password)
				instance.Start()

				instanceList[user.Port] = instance
			}
		}

		for _, instance := range instanceList {
			report(instance, database)
		}

		log.Println("Sync done")
	}
}
