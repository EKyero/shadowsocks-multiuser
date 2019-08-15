package main

import (
	"sync"
	"time"
)

// Bandwidth struct
type Bandwidth struct {
	sync.Mutex
	Upload   uint64
	Download uint64
	Last     int64
}

// IncreaseUpload for increase uploaded size
func (bandwidth *Bandwidth) IncreaseUpload(size uint64) {
	bandwidth.Lock()
	defer bandwidth.Unlock()
	bandwidth.Upload += size
	bandwidth.Last = time.Now().Unix()
}

// IncreaseDownload for increase downloaded size
func (bandwidth *Bandwidth) IncreaseDownload(size uint64) {
	bandwidth.Lock()
	defer bandwidth.Unlock()
	bandwidth.Download += size
	bandwidth.Last = time.Now().Unix()
}

// Reset for reset uploaded and downloaded size
func (bandwidth *Bandwidth) Reset() {
	bandwidth.Lock()
	defer bandwidth.Unlock()
	bandwidth.Upload = 0
	bandwidth.Download = 0
	bandwidth.Last = time.Now().Unix()
}

func newBandwidth() *Bandwidth {
	bandwidth := Bandwidth{}
	bandwidth.Last = 0

	return &bandwidth
}
