package ingest

import (
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/alejandrolaguna20/ccsds-rx/pkg/ccsds"
)

type SatNOGSResponse struct {
	Results []SatNOGSFrame `json:"results"`
}

type SatNOGSFrame struct {
	Timestamp string `json:"timestamp"`
	Frame     string `json:"frame"`
}

func FetchLiveTelemetry(token string, noradID string, limit int) ([]ccsds.SpacePacket, [][]byte, error) {
	url := fmt.Sprintf("https://db.satnogs.org/api/telemetry/?satellite=%s&limit=%d", noradID, limit)

	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to SatNOGS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("satnogs api returned status: %s", resp.Status)
	}

	var satResp SatNOGSResponse
	if err := json.NewDecoder(resp.Body).Decode(&satResp); err != nil {
		return nil, nil, fmt.Errorf("failed to decode json: %w", err)
	}

	var packets []ccsds.SpacePacket
	var rawFrames [][]byte
	for _, f := range satResp.Results {
		data, err := hex.DecodeString(f.Frame)
		if err != nil {
			continue
		}
		rawFrames = append(rawFrames, data)
		packet, err := SynchronizePacket(data)
		if err == nil {
			packets = append(packets, packet)
		}
	}

	return packets, rawFrames, nil
}

func SynchronizePacket(frame []byte) (ccsds.SpacePacket, error) {
	for i := 0; i <= len(frame)-ccsds.PrimaryHeaderSize; i++ {
		if (frame[i] >> 5) == 0 {
			packet, err := ccsds.NewSpacePacket(frame[i:])
			if err == nil {
				return packet, nil
			}
		}
	}
	return nil, fmt.Errorf("no recognizable CCSDS packet found")
}

func GenerateMockStream(count int) [][]byte {
	stream := make([][]byte, count)
	for i := range count {
		noiseSize := rand.IntN(12) + 8
		frame := make([]byte, noiseSize+ccsds.PrimaryHeaderSize+10)
		crand.Read(frame[:noiseSize])

		apid := uint16(rand.IntN(2048))
		frame[noiseSize] = 0x08 | uint8(apid>>8)
		frame[noiseSize+1] = uint8(apid & 0xFF)
		frame[noiseSize+2] = 0xC0
		frame[noiseSize+3] = uint8(i % 256)
		frame[noiseSize+4] = 0x00
		frame[noiseSize+5] = 0x09

		copy(frame[noiseSize+6:], []byte("TELEMETRY!"))

		stream[i] = frame
	}
	return stream
}
