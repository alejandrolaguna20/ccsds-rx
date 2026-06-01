package state

import (
	"testing"

	"github.com/alejandrolaguna20/ccsds-rx/pkg/ccsds"
)

func TestTracker_Update(t *testing.T) {
	tracker := NewTracker()
	packet := ccsds.SpacePacket{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	tracker.Update(packet)

	state := tracker.GetState()
	if state.PacketCount != 1 {
		t.Errorf("expected PacketCount 1, got %d", state.PacketCount)
	}

	if state.LastUpdate.IsZero() {
		t.Error("expected LastUpdate to be set")
	}

	if len(state.LastRawPayload) != len(packet) {
		t.Errorf("expected LastRawPayload length %d, got %d", len(packet), len(state.LastRawPayload))
	}
}

func TestTracker_Concurrency(t *testing.T) {
	tracker := NewTracker()
	packet := ccsds.SpacePacket{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	count := 1000

	done := make(chan bool)
	for i := 0; i < count; i++ {
		go func() {
			tracker.Update(packet)
			done <- true
		}()
	}

	for i := 0; i < count; i++ {
		<-done
	}

	state := tracker.GetState()
	if state.PacketCount != uint64(count) {
		t.Errorf("expected PacketCount %d, got %d", count, state.PacketCount)
	}
}
