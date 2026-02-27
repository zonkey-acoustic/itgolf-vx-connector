package protee

import (
	"encoding/json"
	"strconv"
	"strings"
)

// ShotData represents the top-level ProTee ShotData.json structure
type ShotData struct {
	ClubData        ClubData        `json:"ClubData"`
	BallData        BallData        `json:"BallData"`
	FlightData      json.RawMessage `json:"FlightData"`
	PhysicsSettings json.RawMessage `json:"PhysicsSettings"`
	IsRealShot      bool            `json:"IsRealShot"`
}

// ClubData represents club metrics from ProTee
type ClubData struct {
	Speed        string `json:"Speed"`
	SwingPath    string `json:"SwingPath"`
	FaceAngle    string `json:"FaceAngle"`
	AttackAngle  string `json:"AttackAngle"`
	Loft         string `json:"Loft"`
	Lie          string `json:"Lie"`
	ClosureRate  string `json:"ClosureRate"`
	ImpactPointX string `json:"ImpactPointX"`
	ImpactPointY string `json:"ImpactPointY"`
}

// BallData represents ball metrics from ProTee
type BallData struct {
	Speed           string `json:"Speed"`
	LaunchAngle     string `json:"LaunchAngle"`
	LaunchDirection string `json:"LaunchDirection"`
	TotalSpin       string `json:"TotalSpin"`
	BackSpin        string `json:"BackSpin"`
	SideSpin        string `json:"SideSpin"`
	SpinAxis        string `json:"SpinAxis"`
}


// parseFloat extracts a numeric value from a ProTee string with unit suffix.
// Examples: "148.5 mph" → 148.5, "11.6°" → 11.6, "2696 RPM" → 2696.0, "-2.3°" → -2.3
func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	// Strip known unit suffixes
	s = strings.TrimRight(s, "°")
	s = strings.TrimSuffix(s, " mph")
	s = strings.TrimSuffix(s, " RPM")
	s = strings.TrimSuffix(s, " °/s")
	s = strings.TrimSuffix(s, " inch")
	s = strings.TrimSuffix(s, " yards")
	s = strings.TrimSuffix(s, " ft")
	s = strings.TrimSuffix(s, " sec")
	s = strings.TrimSuffix(s, " %")
	s = strings.TrimSpace(s)

	return strconv.ParseFloat(s, 64)
}

// parseInt extracts an integer value from a ProTee string with unit suffix.
// Examples: "2696 RPM" → 2696, "545 RPM" → 545
func parseInt(s string) (int16, error) {
	f, err := parseFloat(s)
	if err != nil {
		return 0, err
	}
	return int16(f), nil
}
