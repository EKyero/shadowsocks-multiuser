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

		if !contains {
			instance.Stop()
		}
	}
}

func report(instance *Instance, database *Database) {
	if instance.Bandwidth.Upload != 0 || instance.Bandwidth.Download != 0 {
		log.Println(instance.Bandwidth)
		if err := database.UpdateBandwidth(instance.Port, instance.Bandwidth.Upload, instance.Bandwidth.Download); err == nil {
			instance.Bandwidth.Reset()
		} else {
			log.Println(err)
		}
	}
}

func update(instance *Instance, method, password string) {
	if instance.Method != method || instance.Password != password {
		instance.Method = method
		instance.Password = password

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
	log.Println("shadowsocks-multiuser 1.0.0")

	var flags struct {
		ListCipher bool
		DbHost     string
		DbPort     int
		DbUser     string
		DbPass     string
		DbName     string
	}

	flag.BoolVar(&flags.ListCipher, "listcipher", false, "list cipher")
	flag.StringVar(&flags.DbHost, "dbhost", "localhost", "database host")
	flag.IntVar(&flags.DbPort, "dbport", 3306, "database port")
	flag.StringVar(&flags.DbUser, "dbuser", "root", "database user")
	flag.StringVar(&flags.DbPass, "dbpass", "123456", "database pass")
	flag.StringVar(&flags.DbName, "dbname", "sspanel", "database name")
	flag.Parse()

	if flags.ListCipher {
		for _, v := range core.ListCipher() {
			fmt.Println(v)
		}

		return
	}

	instanceList := make(map[int]*Instance, 65535)

	for {
		time.Sleep(10 * time.Second)

		database := newDatabase(flags.DbHost, flags.DbPort, flags.DbUser, flags.DbPass, flags.DbName)
		if err := database.Open(); err != nil {
			log.Println(err)
			continue
		}

		users, err := database.GetUser()
		if err != nil {
			log.Println(err)
			continue
		}

		purge(instanceList, users)

		for _, user := range users {
			if instance, ok := instanceList[user.Port]; ok {
				if user.Enable == 1 {
					update(instance, user.Method, user.Password)
				} else {
					if instance.Started {
						instance.Stop()
					}

					report(instance, database)
					delete(instanceList, user.Port)
				}
			} else {
				if user.Enable == 1 {
					instance := newInstance(user.Port, user.Method, user.Password)
					instance.Start()

					instanceList[user.Port] = instance
				}
			}
		}

		for _, instance := range instanceList {
			report(instance, database)
		}

		log.Println("done")
	}
}
