# VX Connector

An unofficial connector for ProTee VX launch monitors with Infinite Tees integration.

## Overview

VX Connector is a Go-based application that reads shot data from the ProTee VX launch monitor and forwards it to Infinite Tees (or any GSPro-compatible simulator). It provides:

- **ProTee VX integration** — watches for shot data files from ProTee Labs software
- **Infinite Tees connectivity** — forwards ball and club metrics via the GSPro connector protocol
- **Web-based UI** for monitoring and control
- **Full data forwarding** — all 7 ball fields and 9 club fields are sent

## Quick Start

1. Download `vx-connector.exe` from the [latest release](https://github.com/zonkey-acoustic/itgolf-vx-connector/releases/latest)
2. Run `vx-connector.exe` — the web UI opens automatically at `http://localhost:8080`
3. Click the **Connect** button in the Infinite Tees section of the web UI to establish the connection
4. Hit shots on the ProTee VX

## Data Flow

```
ProTee Labs software → ShotData.json files → VX Connector → Infinite Tees
```

The ProTee Labs application writes shot data to `%AppData%\ProTeeUnited\Shots\{timestamp}\ShotData.json`. VX Connector watches this directory, parses each shot, and sends the metrics to Infinite Tees in real-time.

### Ball Data (7 fields)

| ProTee Field | Sent As |
|---|---|
| BallData.Speed | Speed (mph) |
| BallData.LaunchAngle | VLA |
| BallData.LaunchDirection | HLA |
| BallData.TotalSpin | TotalSpin |
| BallData.BackSpin | BackSpin |
| BallData.SideSpin | SideSpin |
| BallData.SpinAxis | SpinAxis |

### Club Data (9 fields)

| ProTee Field | Sent As |
|---|---|
| ClubData.Speed | Speed / SpeedAtImpact |
| ClubData.AttackAngle | AngleOfAttack |
| ClubData.FaceAngle | FaceToTarget |
| ClubData.Loft | Loft |
| ClubData.SwingPath | Path |
| ClubData.Lie | Lie |
| ClubData.ClosureRate | ClosureRate |
| ClubData.ImpactPointX | HorizontalFaceImpact |
| ClubData.ImpactPointY | VerticalFaceImpact |

## Requirements

- Windows PC
- ProTee VX launch monitor with ProTee Labs software installed

## Building from Source

Requires Go 1.23 or later.

```bash
git clone https://github.com/zonkey-acoustic/itgolf-vx-connector.git
cd itgolf-vx-connector

go mod download
go build -o vx-connector.exe main.go
```

## Usage

The web interface will be available at `http://localhost:8080`. The app automatically:
- Starts watching the ProTee Shots directory
- Opens the web browser

### Configuration via Settings

Infinite Tees connection settings (IP, port, auto-connect) can be configured in the Settings page of the web UI.

### Simulation Mode

```bash
# Generate fake shot data every 20 seconds (for testing without hardware)
./vx-connector --mock=simulate
```

The simulator waits for Infinite Tees to connect before sending shots.

### Custom ProTee Path

```bash
# Override the default ProTee shots directory
./vx-connector --protee-path="C:\path\to\shots"
```

## Command-Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `--mock` | "" | Mock mode: 'stub' or 'simulate' |
| `--web-port` | 8080 | Port for web server |
| `--it-ip` | 127.0.0.1 | Infinite Tees server IP |
| `--it-port` | 999 | Infinite Tees server port |
| `--enable-protee` | true | Enable ProTee VX integration |
| `--protee-path` | auto-detected | Path to ProTee Shots directory |
| `--headless` | false | Run in CLI mode without web UI |

## Configuration

Settings are automatically saved and loaded from `~/.vx-connector/config.json`.

## Project Structure

```
.
├── main.go                          # Application entry point
├── internal/
│   ├── core/                        # Core business logic
│   │   ├── state_manager.go         # Application state with reactive callbacks
│   │   ├── parse_notifications.go   # BallMetrics and ClubMetrics structs
│   │   ├── protee/                  # ProTee VX integration
│   │   │   ├── models.go            # ProTee JSON schema
│   │   │   ├── data_conversion.go   # ProTee → internal metrics conversion
│   │   │   ├── integration.go       # File watcher manager
│   │   │   └── simulator.go         # Fake shot generator for testing
│   │   └── infinitetees/            # Infinite Tees integration
│   ├── config/                      # Configuration management
│   ├── web/                         # Web server and API
│   └── logging/                     # Logging utilities
├── web/                             # Frontend assets
│   ├── index.html
│   └── static/
│       ├── css/
│       └── js/
└── docs/
    └── protee-data-mapping.md       # Full data mapping reference
```

## How It Works

1. **File Watching**: Polls `%AppData%\ProTeeUnited\Shots\` every 500ms for new shot directories
2. **Data Parsing**: Reads `ShotData.json`, parses string values with embedded units (e.g., `"148.5 mph"`, `"11.6°"`)
3. **Conversion**: Converts ball speed from mph to m/s, maps all fields to internal metrics
4. **State Management**: Pushes metrics to the StateManager, which triggers callbacks
5. **Forwarding**: Infinite Tees integration picks up metrics via callbacks and sends over TCP

## Troubleshooting

### ProTee shots not detected

- Verify ProTee Labs is writing to `%AppData%\ProTeeUnited\Shots\`
- Check the Watch Path shown in the web UI matches your ProTee configuration
- Ensure the watcher status shows "Watching"

### Infinite Tees not receiving data

- Confirm Infinite Tees is running and listening on the configured port (default 999)
- Check the IT connection status in the web UI
- Verify firewall settings allow the connection

### Web UI not accessible

- Confirm port 8080 is not already in use
- Try a different port with `--web-port`

## Acknowledgements

This project is based on [squaregolf-connector](https://github.com/brentyates/squaregolf-connector) by Brent Yates, which provides Bluetooth connectivity for SquareGolf launch monitors with GSPro integration. The core architecture, state management, and web UI were adapted from that project.

## License

See [LICENSE](LICENSE) file for details.

## Disclaimer

This is an unofficial, community-developed connector and is not affiliated with or endorsed by ProTee or Infinite Tees.
