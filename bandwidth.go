package main

import "sync"

// Bandwidth struct
type Bandwidth struct {
	sync.Mutex
	Upload   uint64
	Download uint64
}

// IncreaseUpload for increase uploaded size
func (bandwidth *Bandwidth) IncreaseUpload(size uint64) {
	bandwidth.Lock()
	defer bandwidth.Unlock()
	bandwidth.Upload += size
}

// IncreaseDownload for increase downloaded size
func (bandwidth *Bandwidth) IncreaseDownload(size uint64) {
	bandwidth.Lock()
	defer bandwidth.Unlock()
	bandwidth.Download += size
}

// Reset for reset uploaded and downloaded size
func (bandwidth *Bandwidth) Reset() {
	bandwidth.Lock()
	defer bandwidth.Unlock()
	bandwidth.Upload = 0
	bandwidth.Download = 0
}

func newBandwidth() *Bandwidth {
	bandwidth := Bandwidth{}

	return &bandwidth
}
