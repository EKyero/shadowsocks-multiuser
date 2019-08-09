package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/NetchX/shadowsocks-multiuser/socks"
)

const udpBufferSize = 64 * 1024

func udpRemote(instance *Instance, cipher func(net.PacketConn) net.PacketConn) {
	socket, err := net.ListenPacket("udp", fmt.Sprintf(":%d", instance.Port))
	if err != nil {
		log.Printf("Failed to listen UDP on %d: %v", instance.Port, err)
		return
	}
	defer socket.Close()
	socket = cipher(socket)

	nat := newNAT(1 * time.Minute)
	buffer := make([]byte, udpBufferSize)

	instance.UDPStarted = true

	for instance.Started {
		size, remoteAddress, err := socket.ReadFrom(buffer)
		if err != nil {
			continue
		}

		targetAddress := socks.SplitAddr(buffer[:size])
		if targetAddress == nil {
			continue
		}

		targetUDPAddress, err := net.ResolveUDPAddr("udp", targetAddress.String())
		if err != nil {
			continue
		}

		data := buffer[len(targetAddress):size]
		conn := nat.Get(remoteAddress.String())
		if conn == nil {
			conn, err = net.ListenPacket("udp", "")
			if err != nil {
				continue
			}

			nat.Add(instance, conn, socket, remoteAddress)
		}

		size, err = conn.WriteTo(data, targetUDPAddress)
		if err != nil {
			continue
		}

		instance.Bandwidth.IncreaseUpload(uint64(size))
	}

	instance.UDPStarted = false
}

// NAT struct
type NAT struct {
	sync.RWMutex
	Map     map[string]net.PacketConn
	Timeout time.Duration
}

// Get PacketConn
func (nat *NAT) Get(id string) net.PacketConn {
	nat.RLock()
	defer nat.RUnlock()

	return nat.Map[id]
}

// Set PacketConn
func (nat *NAT) Set(id string, conn net.PacketConn) {
	nat.Lock()
	defer nat.Unlock()

	nat.Map[id] = conn
}

// Add NAT
func (nat *NAT) Add(instance *Instance, src, dst net.PacketConn, peer net.Addr) {
	nat.Set(peer.String(), src)

	go func() {
		udpRelay(instance, src, dst, peer, nat.Timeout)
		if conn := nat.Delete(peer.String()); conn != nil {
			conn.Close()
		}
	}()
}

// Delete NAT
func (nat *NAT) Delete(id string) net.PacketConn {
	nat.Lock()
	defer nat.Unlock()

	conn, ok := nat.Map[id]
	if ok {
		delete(nat.Map, id)

		return conn
	}

	return nil
}

func newNAT(timeout time.Duration) *NAT {
	nat := NAT{}
	nat.Map = make(map[string]net.PacketConn)
	nat.Timeout = timeout

	return &nat
}

func udpRelay(instance *Instance, src, dst net.PacketConn, target net.Addr, timeout time.Duration) error {
	buffer := make([]byte, udpBufferSize)

	for {
		src.SetReadDeadline(time.Now().Add(timeout))
		size, remoteAddress, err := src.ReadFrom(buffer)
		if err != nil {
			return err
		}

		sourceAddress := socks.ParseAddr(remoteAddress.String())
		copy(buffer[len(sourceAddress):], buffer[:size])
		copy(buffer, sourceAddress)

		size, err = dst.WriteTo(buffer[:len(sourceAddress)+size], target)

		if err != nil {
			return err
		}

		instance.Bandwidth.IncreaseDownload(uint64(size))
	}
}
