package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/NetchX/shadowsocks-multiuser/core"
)

var flags struct {
	ListCipher   bool
	DBHost       string
	DBPort       int
	DBUser       string
	DBPass       string
	DBName       string
	NodeID       int
	UDPEnabled   bool
	SyncInterval int
}

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
		if time.Now().Unix()-instance.Bandwidth.Last > 10 {
			if err := database.UpdateBandwidth(instance); err == nil {
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

	if instance.Started {
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
	flag.BoolVar(&flags.ListCipher, "listcipher", false, "List all cipher")
	flag.StringVar(&flags.DBHost, "dbhost", "localhost", "Database hostname")
	flag.IntVar(&flags.DBPort, "dbport", 3306, "Database port")
	flag.StringVar(&flags.DBUser, "dbuser", "root", "Database username")
	flag.StringVar(&flags.DBPass, "dbpass", "123456", "Database password")
	flag.StringVar(&flags.DBName, "dbname", "sspanel", "Database name")
	flag.IntVar(&flags.NodeID, "nodeid", -1, "Node ID")
	flag.IntVar(&flags.SyncInterval, "syncinterval", 30, "Sync interval")
	flag.BoolVar(&flags.UDPEnabled, "udp", false, "UDP forward")
	flag.Parse()

	if flags.ListCipher {
		for _, v := range core.ListCipher() {
			fmt.Println(v)
		}

		return
	}

	log.Println("Starting shadowsocks-multiuser")
	log.Println("Version: 1.0.0")

	if flags.NodeID == -1 {
		log.Println("Node id must be specified")
		return
	}

	instanceList := make(map[int]*Instance, 65535)
	first := true

	log.Println("Started")
	for {
		if !first {
			log.Printf("Wait %d seconds for sync users", flags.SyncInterval)
			time.Sleep(time.Second * time.Duration(flags.SyncInterval))
		} else {
			first = false
		}

		log.Println("Start syncing")

		log.Println("Opening database connection")
		database := newDatabase(flags.DBHost, flags.DBPort, flags.DBUser, flags.DBPass, flags.DBName, flags.NodeID)
		if err := database.Open(); err != nil {
			log.Println(err)
			continue
		}

		defer database.Close()

		log.Println("Update heartbeat")
		database.UpdateHeartbeat()

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
				log.Printf("Starting new instance for %d", user.ID)
				instance := newInstance(user.ID, user.Port, user.Method, user.Password)
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
