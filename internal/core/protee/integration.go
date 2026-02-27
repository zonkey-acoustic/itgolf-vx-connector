package protee

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/brentyates/squaregolf-connector/internal/core"
)

const (
	pollInterval = 500 * time.Millisecond
	writeDelay   = 200 * time.Millisecond
)

// Manager watches a ProTee Shots directory for new ShotData.json files
type Manager struct {
	stateManager    *core.StateManager
	watchPath       string
	processedDirs   map[string]bool
	lastNewestDir   string
	stopChan        chan struct{}
	running         bool
	mu              sync.RWMutex
}

var (
	instance *Manager
	once     sync.Once
)

// GetInstance returns the singleton Manager instance
func GetInstance(stateManager *core.StateManager) *Manager {
	once.Do(func() {
		instance = &Manager{
			stateManager:  stateManager,
			processedDirs: make(map[string]bool),
		}
	})
	return instance
}

// Start begins watching the ProTee Shots directory
func (m *Manager) Start(watchPath string) error {
	m.mu.Lock()

	if m.running {
		m.mu.Unlock()
		return fmt.Errorf("ProTee watcher already running")
	}

	// Validate directory exists
	info, err := os.Stat(watchPath)
	if err != nil {
		m.mu.Unlock()
		m.stateManager.SetProTeeStatus(core.ProTeeStatusError)
		m.stateManager.SetProTeeError(fmt.Errorf("shots directory not found: %s", watchPath))
		return fmt.Errorf("shots directory not found: %s", watchPath)
	}
	if !info.IsDir() {
		m.mu.Unlock()
		m.stateManager.SetProTeeStatus(core.ProTeeStatusError)
		m.stateManager.SetProTeeError(fmt.Errorf("path is not a directory: %s", watchPath))
		return fmt.Errorf("path is not a directory: %s", watchPath)
	}

	m.watchPath = watchPath
	m.stopChan = make(chan struct{})
	m.running = true

	// Snapshot existing directories so we only process new ones
	m.snapshotExistingDirs()

	log.Printf("ProTee: watching for shots in %s", watchPath)

	go m.pollLoop()
	m.mu.Unlock()

	// Set status outside the lock to avoid deadlock with callbacks
	m.stateManager.SetProTeeStatus(core.ProTeeStatusWatching)
	m.stateManager.SetProTeeError(nil)

	return nil
}

// Stop stops watching
func (m *Manager) Stop() {
	m.mu.Lock()

	if !m.running {
		m.mu.Unlock()
		return
	}

	close(m.stopChan)
	m.running = false
	m.mu.Unlock()

	// Set status outside the lock to avoid deadlock with callbacks
	m.stateManager.SetProTeeStatus(core.ProTeeStatusDisabled)
	m.stateManager.SetProTeeError(nil)
	log.Println("ProTee: watcher stopped")
}

// IsRunning returns whether the watcher is currently active
func (m *Manager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// GetWatchPath returns the current watch path
func (m *Manager) GetWatchPath() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.watchPath
}

// snapshotExistingDirs records all current subdirectories so they're not reprocessed
func (m *Manager) snapshotExistingDirs() {
	entries, err := os.ReadDir(m.watchPath)
	if err != nil {
		log.Printf("ProTee: failed to snapshot existing dirs: %v", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			m.processedDirs[entry.Name()] = true
		}
	}
	log.Printf("ProTee: snapshotted %d existing shot directories", len(m.processedDirs))
}

// pollLoop polls the shots directory for new subdirectories
func (m *Manager) pollLoop() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			m.checkForNewShots()
		}
	}
}

// checkForNewShots looks for new shot directories
func (m *Manager) checkForNewShots() {
	entries, err := os.ReadDir(m.watchPath)
	if err != nil {
		log.Printf("ProTee: failed to read shots directory: %v", err)
		m.stateManager.SetProTeeStatus(core.ProTeeStatusError)
		m.stateManager.SetProTeeError(fmt.Errorf("failed to read shots directory: %v", err))
		return
	}

	// Find new directories
	var newDirs []string
	for _, entry := range entries {
		if entry.IsDir() && !m.processedDirs[entry.Name()] {
			newDirs = append(newDirs, entry.Name())
		}
	}

	if len(newDirs) == 0 {
		return
	}

	// Sort to process in order (timestamps sort naturally)
	sort.Strings(newDirs)

	for _, dirName := range newDirs {
		m.processedDirs[dirName] = true
		shotFile := filepath.Join(m.watchPath, dirName, "ShotData.json")

		// Wait for ProTee to finish writing the file
		time.Sleep(writeDelay)

		m.processShotFile(shotFile, dirName)
	}
}

// processShotFile reads and processes a single ShotData.json file
func (m *Manager) processShotFile(path string, dirName string) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Retry once after a short delay (file might still be writing)
		time.Sleep(writeDelay)
		data, err = os.ReadFile(path)
		if err != nil {
			log.Printf("ProTee: failed to read %s: %v", path, err)
			return
		}
	}

	var shot ShotData
	if err := json.Unmarshal(data, &shot); err != nil {
		log.Printf("ProTee: failed to parse %s: %v", path, err)
		return
	}

	// Skip non-real shots
	if !shot.IsRealShot {
		log.Printf("ProTee: skipping non-real shot in %s", dirName)
		return
	}

	log.Printf("ProTee: processing shot from %s", dirName)

	// Convert and push metrics to state manager
	ballMetrics := convertBallMetrics(&shot)
	clubMetrics := convertClubMetrics(&shot)

	log.Printf("ProTee: BALL speed=%.1f m/s (%.1f mph), launch=%.1f°, dir=%.1f°, totalSpin=%d, backSpin=%d, sideSpin=%d, axis=%.1f°",
		ballMetrics.BallSpeedMPS, ballMetrics.BallSpeedMPS*2.23694,
		ballMetrics.VerticalAngle, ballMetrics.HorizontalAngle,
		ballMetrics.TotalspinRPM, ballMetrics.BackspinRPM, ballMetrics.SidespinRPM,
		ballMetrics.SpinAxis)
	log.Printf("ProTee: CLUB speed=%.1f mph, attack=%.1f°, path=%.1f°, face=%.1f°, loft=%.1f°, lie=%.1f°, closure=%.0f °/s, impactX=%.1f, impactY=%.1f",
		clubMetrics.ClubSpeed, clubMetrics.AttackAngle, clubMetrics.PathAngle,
		clubMetrics.FaceAngle, clubMetrics.DynamicLoftAngle, clubMetrics.Lie,
		clubMetrics.ClosureRate, clubMetrics.ImpactPointX, clubMetrics.ImpactPointY)

	m.stateManager.SetLastClubMetrics(clubMetrics)
	m.stateManager.SetLastBallMetrics(ballMetrics)
}
