# `ccsds-rx`: CCSDS Downlink Decoder

A high-performance, zero-copy telemetry engine designed to synchronize and decode CCSDS Space Packet Protocol frames from live satellite observations.

## Prerequisites
- **Go**: Version 1.26 or higher.
- **SatNOGS API Token**: Required for live data ingestion. You can obtain one by creating an account at [db.satnogs.org](https://db.satnogs.org/).

## Setup & Configuration

1. **Clone the repository**:
   ```bash
   git clone https://github.com/alejandrolaguna20/ccsds-rx
   cd ccsds-rx
   ```

2. **Configure Environment**:
   Create a `.env` file in the root directory and add your SatNOGS API token:
   ```env
   SATNOGS_TOKEN=your_token_here
   ```

3. **Install Dependencies**:
   ```bash
   go mod download
   ```
## Usage

### Running CCSDS Monitor
The CCSDS Monitor application targets active satellites (like UmKA-1) and attempts to extract real-time telemetry:

```bash
go run cmd/ccsds-monitor/main.go
```

## License
MIT
