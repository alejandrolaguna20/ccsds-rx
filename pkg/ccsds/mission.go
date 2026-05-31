package ccsds

import "fmt"

type UmKa1Housekeeping struct {
	BatteryVoltage int
	BatteryCurrent int
	SolarVoltage   int
	Temperature    int8
	RSSI           int8
}

func (u *UmKa1Housekeeping) Decode(data []byte) error {
	if len(data) < 8 {
		return fmt.Errorf("USP payload too short")
	}

	u.BatteryVoltage = int(data[0]) | int(data[1])<<8
	u.BatteryCurrent = int(int16(data[2]) | int16(data[3])<<8)
	u.SolarVoltage = int(data[4]) | int(data[5])<<8
	u.Temperature = int8(data[6])
	u.RSSI = int8(data[7])

	return nil
}

func (u *UmKa1Housekeeping) String() string {
	return fmt.Sprintf("BAT: %dmV | CUR: %dmA | SOL: %dmV | TEMP: %d°C | RSSI: %ddBm",
		u.BatteryVoltage, u.BatteryCurrent, u.SolarVoltage, u.Temperature, u.RSSI)
}

func DecodeMissionData(noradID string, packet SpacePacket) {
	switch noradID {
	case "57172":
		if packet.APID() == 0x001 {
			var hk UmKa1Housekeeping
			if err := hk.Decode(packet.UserData()); err == nil {
				fmt.Printf("  [UmKA-1 USP HK] %s\n", hk.String())
			}
		} else {
			fmt.Printf("  [UmKA-1] APID 0x%03X payload: %x\n", packet.APID(), packet.UserData())
		}

	case "39090":
		// TODO: Implement STRAND-1 specific payload decoding
		fmt.Printf("  [STRAND-1] APID 0x%03X payload: %x (Parsing pending)\n", packet.APID(), packet.UserData())

	case "40014":
		// TODO: Implement BUGSAT-1 specific payload decoding
		fmt.Printf("  [BUGSAT-1] APID 0x%03X payload: %x (Parsing pending)\n", packet.APID(), packet.UserData())

	default:
		fmt.Printf("  [Unknown Satellite] APID: 0x%03X | Payload: %x\n", packet.APID(), packet.UserData())
	}
}
