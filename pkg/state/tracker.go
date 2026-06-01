package state

import (
	"sync"
	"time"

	"github.com/alejandrolaguna20/ccsds-rx/pkg/ccsds"
)

type SatelliteState struct {
	LastUpdate     time.Time
	PacketCount    uint64
	BatteryMV      int
	TempC          int8
	RSSI           int8
	LastRawPayload []byte
}

type Tracker struct {
	mu    sync.RWMutex
	state SatelliteState
}

func NewTracker() *Tracker {
	return &Tracker{}
}

func (t *Tracker) Update(p ccsds.SpacePacket) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.state.LastUpdate = time.Now()
	t.state.PacketCount++

	// we must copy the payload because the core parser's
	// buffer will be reused in the hot path
	// we minimize allocation by reusing the existing capacity of the slice
	if cap(t.state.LastRawPayload) < len(p) {
		t.state.LastRawPayload = make([]byte, len(p))
	} else {
		t.state.LastRawPayload = t.state.LastRawPayload[:len(p)]
	}
	copy(t.state.LastRawPayload, p)
}

func (t *Tracker) GetState() SatelliteState {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.state
}
