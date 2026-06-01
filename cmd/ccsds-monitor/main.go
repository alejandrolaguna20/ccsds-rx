package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alejandrolaguna20/ccsds-rx/pkg/ccsds"
	"github.com/alejandrolaguna20/ccsds-rx/pkg/ingest"
	"github.com/alejandrolaguna20/ccsds-rx/pkg/state"
	"github.com/joho/godotenv"
)

type Config struct {
	SatNOGSToken string
}

func getSatelliteData(satNoradID string, config Config, packetChan chan<- ccsds.SpacePacket) {
	fmt.Printf("\nTargeting satellite (NORAD %s)...\n", satNoradID)
	packets, rawFrames, err := ingest.FetchLiveTelemetry(config.SatNOGSToken, satNoradID, 5)
	if err != nil {
		fmt.Printf("[WARN] Skipping %s: %v\n", satNoradID, err)
		return
	}
	if len(packets) == 0 {
		fmt.Printf("[INFO] No CCSDS packets identified for %s. Analyzing radio protocol:\n", satNoradID)
		for _, frame := range rawFrames {
			ccsds.DecodeRawFrame(frame)
		}
		return
	}
	fmt.Printf("[SUCCESS] Extracted %d structured packets from %s stream.\n", len(packets), satNoradID)
	for i, p := range packets {
		fmt.Printf("\n[Packet %d]\n", i)
		fmt.Printf("  APID:      0x%03X (%d)\n", p.APID(), p.APID())
		fmt.Printf("  Seq Count: %d\n", p.SequenceCount())
		fmt.Printf("  Size:      %d bytes\n", p.TotalLength())
		ccsds.DecodeMissionData(satNoradID, p)

		// Pipe packet to State Tracker
		packetChan <- p
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("[WARN] No .env file found, relying on system environment variables")
	}
	token := os.Getenv("SATNOGS_TOKEN")
	if token == "" {
		fmt.Println("[ERROR] SATNOGS_TOKEN not set in environment")
		return
	}
	config := Config{
		SatNOGSToken: token,
	}

	tracker := state.NewTracker()
	packetChan := make(chan ccsds.SpacePacket, 100)

	go func() {
		for p := range packetChan {
			tracker.Update(p)
		}
	}()

	id := flag.String("noradID", "0", "Satellite's NORAD ID")
	flag.Parse()

	getSatelliteData(*id, config, packetChan)

	s := tracker.GetState()
	fmt.Printf("\n[STATE TRACKER] Packets Processed: %d | Last Update: %v\n", s.PacketCount, s.LastUpdate.Format("15:04:05"))

	fmt.Println("\nAnalysis complete.")
}
