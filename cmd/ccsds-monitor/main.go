package main

import (
	"fmt"
	"os"

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
		{"57172", "UmKA-1"},
		{"39090", "STRAND-1"},
		{"40014", "BUGSAT-1"},
	}
	for _, sat := range satellites {
		fmt.Printf("\nTargeting satellite %s (NORAD %s)...\n", sat.Name, sat.ID)
		packets, err := ingest.FetchLiveTelemetry(token, sat.ID, 10)
		if err != nil {
			fmt.Printf("[WARN] Skipping %s: %v\n", sat.Name, err)
			continue
		}
		if len(packets) == 0 {
			fmt.Printf("[INFO] No CCSDS packets identified in the latest radio frames for %s.\n", sat.Name)
			continue
		}
		fmt.Printf("[SUCCESS] Extracted %d packets from %s stream.\n", len(packets), sat.Name)
		for i, p := range packets {
			fmt.Printf("\n[Packet %d]\n", i)
			fmt.Printf("  APID:      0x%03X (%d)\n", p.APID(), p.APID())
			fmt.Printf("  Seq Count: %d\n", p.SequenceCount())
			fmt.Printf("  Size:      %d bytes\n", p.TotalLength())
			if len(p.UserData()) > 0 {
				fmt.Printf("  Payload:   %x\n", p.UserData())
			}
		}
	}
	fmt.Println("\n[FAILURE] Analysis complete: No valid CCSDS packets identified in active feeds.")
}
