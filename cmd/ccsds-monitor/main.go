package main

import (
	"fmt"
	"os"

	"github.com/alejandrolaguna20/ccsds-rx/pkg/ccsds"
	"github.com/alejandrolaguna20/ccsds-rx/pkg/ingest"
	"github.com/joho/godotenv"
)

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
	satellites := []struct {
		ID   string
		Name string
	}{
		{"25544", "ISS (Zarya)"},
		{"57172", "UmKA-1"},
	}
	for _, sat := range satellites {
		fmt.Printf("\nTargeting satellite %s (NORAD %s)...\n", sat.Name, sat.ID)
		packets, rawFrames, err := ingest.FetchLiveTelemetry(token, sat.ID, 5)
		if err != nil {
			fmt.Printf("[WARN] Skipping %s: %v\n", sat.Name, err)
			continue
		}
		if len(packets) == 0 {
			fmt.Printf("[INFO] No CCSDS packets identified for %s. Analyzing radio protocol:\n", sat.Name)
			for _, frame := range rawFrames {
				ccsds.DecodeRawFrame(frame)
			}
			continue
		}
		fmt.Printf("[SUCCESS] Extracted %d structured packets from %s stream.\n", len(packets), sat.Name)
		for i, p := range packets {
			fmt.Printf("\n[Packet %d]\n", i)
			fmt.Printf("  APID:      0x%03X (%d)\n", p.APID(), p.APID())
			fmt.Printf("  Seq Count: %d\n", p.SequenceCount())
			fmt.Printf("  Size:      %d bytes\n", p.TotalLength())
			ccsds.DecodeMissionData(sat.ID, p)
		}
	}
	fmt.Println("\nAnalysis complete.")
}
