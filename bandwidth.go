package main

// Bandwidth struct
type Bandwidth struct {
	Upload   uint64
	Download uint64
}

// IncreaseUpload for increase uploaded size
func (bandwidth *Bandwidth) IncreaseUpload(size uint64) {
	bandwidth.Upload += size
}

// IncreaseDownload for increase downloaded size
func (bandwidth *Bandwidth) IncreaseDownload(size uint64) {
	bandwidth.Download += size
}

// Reset for reset uploaded and downloaded size
func (bandwidth *Bandwidth) Reset() {
	bandwidth.Upload = 0
	bandwidth.Download = 0
}

func newBandwidth() *Bandwidth {
	bandwidth := Bandwidth{}

	return &bandwidth
}
