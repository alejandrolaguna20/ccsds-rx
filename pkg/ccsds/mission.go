package ccsds

import (
	"fmt"
	"unicode"
)

func DecodeMissionData(noradID string, packet SpacePacket) {
	switch noradID {
	case "57172": // UmKA-1
		fmt.Printf("  [UmKA-1] APID 0x%03X payload: %x\n", packet.APID(), packet.UserData())
	default:
		fmt.Printf("  [Generic] APID: 0x%03X | Payload: %x\n", packet.APID(), packet.UserData())
	}
}

func DecodeRawFrame(frame []byte) {
	if isAX25(frame) {
		dest, src := parseAX25Address(frame)
		payload := ""
		if len(frame) > 16 {
			payload = ExtractPrintableText(frame[16:])
		}
		fmt.Printf("  [AX.25 Frame] From: %s | To: %s | Msg: %s\n", src, dest, payload)
		return
	}
	fmt.Printf("  [Unknown Protocol] Size: %d bytes | Hex: %x\n", len(frame), frame)
}

func isAX25(data []byte) bool {
	if len(data) < 14 {
		return false
	}
	for i := 0; i < 7; i++ {
		char := rune(data[i] >> 1)
		if !unicode.IsPrint(char) && !unicode.IsSpace(char) {
			return false
		}
	}
	return true
}

func parseAX25Address(data []byte) (string, string) {
	dest := ""
	for i := 0; i < 6; i++ {
		dest += string(rune(data[i] >> 1))
	}
	src := ""
	for i := 7; i < 13; i++ {
		src += string(rune(data[i] >> 1))
	}
	return dest, src
}

func ExtractPrintableText(data []byte) string {
	var result []rune
	for _, b := range data {
		if unicode.IsPrint(rune(b)) {
			result = append(result, rune(b))
		} else if b == 0x0D || b == 0x0A {
			result = append(result, ' ')
		}
	}
	return string(result)
}
