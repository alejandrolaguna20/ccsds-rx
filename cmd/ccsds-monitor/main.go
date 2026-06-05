package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/alejandrolaguna20/ccsds-rx/internal/tui"
	"github.com/alejandrolaguna20/ccsds-rx/pkg/ccsds"
	"github.com/alejandrolaguna20/ccsds-rx/pkg/ingest"
	"github.com/alejandrolaguna20/ccsds-rx/pkg/state"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

type Config struct {
	SatNOGSToken string
}

func getSatelliteData(satNoradID string, config Config, packetChan chan<- ccsds.SpacePacket) {
	packets, rawFrames, err := ingest.FetchLiveTelemetry(config.SatNOGSToken, satNoradID, 10)
	if err != nil {
		return
	}
	if len(packets) == 0 {
		for _, frame := range rawFrames {
			ccsds.DecodeRawFrame(frame)
		}
		return
	}
	for _, p := range packets {
		packetChan <- p
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		// Silent if no .env
	}
	token := os.Getenv("SATNOGS_TOKEN")
	if token == "" {
		fmt.Println("[ERROR] SATNOGS_TOKEN not set in environment")
		return
	}
	config := Config{
		SatNOGSToken: token,
	}

	id := flag.String("noradID", "57172", "Satellite's NORAD ID (default: UmKA-1)")
	flag.Parse()

	tracker := state.NewTracker()
	packetChan := make(chan ccsds.SpacePacket, 100)

	go func() {
		for p := range packetChan {
			tracker.Update(p)
		}
	}()

	go getSatelliteData(*id, config, packetChan)

	m := tui.NewModel(tracker, *id)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
