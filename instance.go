package main

import (
	"fmt"
	"log"
	"net"

	"github.com/NetchX/shadowsocks-multiuser/core"
)

// Instance struct
type Instance struct {
	UserID     int
	Port       int
	Method     string
	Password   string
	Bandwidth  *Bandwidth
	Started    bool
	TCPStarted bool
	UDPStarted bool
}

// Start instance
func (instance *Instance) Start() {
	cipher, err := core.PickCipher(instance.Method, make([]byte, 0), instance.Password)
	if err != nil {
		log.Println(err)
		return
	}

	instance.Started = true

	go tcpRemote(instance, cipher.StreamConn)
	go udpRemote(instance, cipher.PacketConn)
}

// Stop instance
func (instance *Instance) Stop() {
	instance.Started = false

	tcp, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", instance.Port))
	if err == nil {
		tcp.Close()
	}

	udp, err := net.Dial("udp", fmt.Sprintf("127.0.0.1:%d", instance.Port))
	if err == nil {
		fmt.Fprint(udp, "NMSL")

		udp.Close()
	}
}

func newInstance(id int, port int, method, password string) *Instance {
	instance := Instance{}
	instance.UserID = id
	instance.Port = port
	instance.Method = method
	instance.Password = password
	instance.Bandwidth = newBandwidth()
	instance.Started = false
	instance.TCPStarted = false
	instance.UDPStarted = false

	return &instance
}
