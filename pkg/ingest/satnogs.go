package ingest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
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

func FetchLiveTelemetry(token string, noradID string, limit int) ([]ccsds.SpacePacket, error) {
	url := fmt.Sprintf("https://db.satnogs.org/api/telemetry/?satellite=%s&limit=%d", noradID, limit)

	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SatNOGS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("satnogs api returned status: %s (check your token)", resp.Status)
	}

	var satResp SatNOGSResponse
	if err := json.NewDecoder(resp.Body).Decode(&satResp); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	var packets []ccsds.SpacePacket
	for _, f := range satResp.Results {
		data, err := hex.DecodeString(f.Frame)
		if err != nil {
			continue
		}

		// TODO: Implement AX.25/FX.25 link-layer header parsing to extract 
		// metadata such as source and destination callsigns.
		packet, err := SynchronizePacket(data)
		if err == nil {
		    packets = append(packets, packet)
		}
		}

		return packets, nil
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
		return nil, fmt.Errorf("no valid CCSDS packet found in frame")
		}
func GenerateMockStream(count int) [][]byte {
    stream := make([][]byte, count)
    for i := 0; i < count; i++ {
        noiseSize := rand.Intn(12) + 8
        frame := make([]byte, noiseSize+6+10)
        rand.Read(frame[:noiseSize])

        apid := uint16(rand.Intn(2048))
        frame[noiseSize] = 0x08 | uint8(apid>>8)
        frame[noiseSize+1] = uint8(apid & 0xFF)

        seq := uint16(i % 16384)
        frame[noiseSize+2] = 0xC0 | uint8(seq>>8)
        frame[noiseSize+3] = uint8(seq & 0xFF)

        frame[noiseSize+4] = 0x00
        frame[noiseSize+5] = 0x09

        payload := make([]byte, 10)
        if apid == 0x10 {
            payload[0] = uint8(180 + rand.Intn(40))
            payload[1] = uint8(rand.Intn(40) - 10)
            payload[2], payload[3] = uint8(rand.Intn(256)), uint8(rand.Intn(256))
            payload[4], payload[5] = uint8(rand.Intn(256)), uint8(rand.Intn(256))
            payload[6], payload[7] = uint8(rand.Intn(256)), uint8(rand.Intn(256))
        } else {
            copy(payload, []byte("TELEMETRY!"))
        }
        copy(frame[noiseSize+6:], payload)

        stream[i] = frame
    }
    return stream
}
