package protee

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/brentyates/squaregolf-connector/internal/core"
)

const (
	simShotInterval = 20 * time.Second
)

// Simulator writes fake ShotData.json files to the watch directory
type Simulator struct {
	watchPath    string
	stateManager *core.StateManager
	stopChan     chan struct{}
	rng          *rand.Rand
}

// NewSimulator creates a new ProTee shot simulator
func NewSimulator(watchPath string, stateManager *core.StateManager) *Simulator {
	return &Simulator{
		watchPath:    watchPath,
		stateManager: stateManager,
		stopChan:     make(chan struct{}),
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Start begins generating fake shots at regular intervals
func (s *Simulator) Start() {
	log.Printf("ProTee Simulator: generating shots every %v in %s", simShotInterval, s.watchPath)

	// Ensure the watch directory exists
	os.MkdirAll(s.watchPath, 0755)

	go s.loop()
}

// Stop stops the simulator
func (s *Simulator) Stop() {
	close(s.stopChan)
}

func (s *Simulator) loop() {
	// Wait a bit before first shot so the watcher is ready
	time.Sleep(3 * time.Second)

	for {
		if s.stateManager.GetInfiniteTeesStatus() == core.InfiniteTeesStatusConnected {
			s.writeShot()
		} else {
			log.Println("ProTee Simulator: waiting for Infinite Tees connection...")
		}

		select {
		case <-s.stopChan:
			return
		case <-time.After(simShotInterval):
		}
	}
}

func (s *Simulator) writeShot() {
	// Create timestamped directory matching ProTee format
	dirName := time.Now().Format("2006-01-02-150405") + fmt.Sprintf("%03d", s.rng.Intn(1000))
	dirPath := filepath.Join(s.watchPath, dirName)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Printf("ProTee Simulator: failed to create dir: %v", err)
		return
	}

	shot := s.generateShot()
	data, err := json.MarshalIndent(shot, "", "  ")
	if err != nil {
		log.Printf("ProTee Simulator: failed to marshal: %v", err)
		return
	}

	filePath := filepath.Join(dirPath, "ShotData.json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		log.Printf("ProTee Simulator: failed to write: %v", err)
		return
	}

	log.Printf("ProTee Simulator: wrote shot to %s (ball speed: %s)", dirName, shot.BallData.Speed)
}

func (s *Simulator) generateShot() ShotData {
	ballSpeed := 100.0 + s.rng.Float64()*60.0     // 100-160 mph
	launchAngle := 8.0 + s.rng.Float64()*12.0      // 8-20 degrees
	launchDir := (s.rng.Float64() - 0.5) * 10.0    // -5 to +5 degrees
	totalSpin := 2000.0 + s.rng.Float64()*3000.0    // 2000-5000 RPM
	backSpin := totalSpin * (0.85 + s.rng.Float64()*0.1)
	sideSpin := (s.rng.Float64() - 0.5) * 1500.0    // -750 to +750 RPM
	spinAxis := (s.rng.Float64() - 0.5) * 30.0      // -15 to +15 degrees

	clubSpeed := ballSpeed * (0.62 + s.rng.Float64()*0.08) // ~65-70% of ball speed
	swingPath := (s.rng.Float64() - 0.5) * 10.0
	faceAngle := (s.rng.Float64() - 0.5) * 8.0
	attackAngle := -4.0 + s.rng.Float64()*10.0 // -4 to +6
	loft := 10.0 + s.rng.Float64()*15.0
	lie := 4.0 + s.rng.Float64()*6.0
	closureRate := 2000.0 + s.rng.Float64()*3000.0
	impactX := (s.rng.Float64() - 0.5) * 1.5
	impactY := (s.rng.Float64() - 0.5) * 1.0

	return ShotData{
		IsRealShot: true,
		ClubData: ClubData{
			Speed:        fmt.Sprintf("%.1f mph", clubSpeed),
			SwingPath:    fmt.Sprintf("%.1f°", swingPath),
			FaceAngle:    fmt.Sprintf("%.1f°", faceAngle),
			AttackAngle:  fmt.Sprintf("%.1f°", attackAngle),
			Loft:         fmt.Sprintf("%.1f°", loft),
			Lie:          fmt.Sprintf("%.1f°", lie),
			ClosureRate:  fmt.Sprintf("%.0f °/s", closureRate),
			ImpactPointX: fmt.Sprintf("%.1f inch", impactX),
			ImpactPointY: fmt.Sprintf("%.1f inch", impactY),
		},
		BallData: BallData{
			Speed:           fmt.Sprintf("%.1f mph", ballSpeed),
			LaunchAngle:     fmt.Sprintf("%.1f°", launchAngle),
			LaunchDirection: fmt.Sprintf("%.1f°", launchDir),
			TotalSpin:       fmt.Sprintf("%.0f RPM", totalSpin),
			BackSpin:        fmt.Sprintf("%.0f RPM", backSpin),
			SideSpin:        fmt.Sprintf("%.0f RPM", sideSpin),
			SpinAxis:        fmt.Sprintf("%.1f°", spinAxis),
		},
		FlightData:      json.RawMessage(`{}`),
		PhysicsSettings: json.RawMessage(`{}`),
	}
}
