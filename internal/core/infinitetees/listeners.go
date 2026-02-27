package infinitetees

import (
	"log"

	"github.com/brentyates/squaregolf-connector/internal/core"
)

func (it *Integration) registerStateListeners() {
	it.stateManager.RegisterBallReadyCallback(it.onBallReadyChanged)
	it.stateManager.RegisterLastBallMetricsCallback(it.onLastBallMetricsChanged)
	it.stateManager.RegisterLastClubMetricsCallback(it.onLastClubMetricsChanged)
}

func (it *Integration) onBallReadyChanged(oldValue, newValue bool) {
	if oldValue == newValue {
		return
	}

	if !it.Base.Connected || it.Base.Socket == nil {
		return
	}

	emptyShotData := ShotData{
		DeviceID:   "CustomLaunchMonitor",
		Units:      "Yards",
		APIversion: "1",
		ShotNumber: it.lastShotNumber,
		ShotDataOptions: ShotOptions{
			ContainsBallData:          false,
			ContainsClubData:          false,
			LaunchMonitorIsReady:      newValue,
			LaunchMonitorBallDetected: newValue,
		},
	}

	if err := it.sendData(emptyShotData); err != nil {
		log.Printf("[%s] Error sending empty shot data: %v", it.Name(), err)
	}
}

func (it *Integration) onLastBallMetricsChanged(oldValue, newValue *core.BallMetrics) {
	if oldValue == newValue {
		return
	}

	if !it.Base.Connected || it.Base.Socket == nil {
		return
	}

	if newValue == nil {
		return
	}

	shotData := it.convertToShotFormat(*newValue, true)
	log.Printf("[%s] Sending ball data: speed=%.1f mph, VLA=%.1f°, HLA=%.1f°, totalSpin=%d, backSpin=%d, sideSpin=%d, spinAxis=%.1f°",
		it.Name(), shotData.BallData.Speed, shotData.BallData.VLA, shotData.BallData.HLA,
		shotData.BallData.TotalSpin, shotData.BallData.BackSpin, shotData.BallData.SideSpin, shotData.BallData.SpinAxis)
	if err := it.sendData(shotData); err != nil {
		log.Printf("[%s] Error sending shot data: %v", it.Name(), err)
	}
}

func (it *Integration) onLastClubMetricsChanged(oldValue, newValue *core.ClubMetrics) {
	if oldValue == newValue {
		return
	}

	if !it.Base.Connected || it.Base.Socket == nil {
		return
	}

	if newValue == nil {
		zeroedClubData := &ClubData{
			Speed:                0,
			AngleOfAttack:        0,
			FaceToTarget:         0,
			Lie:                  0,
			Loft:                 0,
			Path:                 0,
			SpeedAtImpact:        0,
			VerticalFaceImpact:   0,
			HorizontalFaceImpact: 0,
			ClosureRate:          0,
		}

		shotData := it.convertToShotFormat(core.BallMetrics{}, false)
		shotData.ShotDataOptions.ContainsBallData = false
		shotData.ShotDataOptions.ContainsClubData = true
		shotData.ClubData = zeroedClubData
		if err := it.sendData(shotData); err != nil {
			log.Printf("[%s] Error sending zeroed club data: %v", it.Name(), err)
		}
		return
	}

	clubData := it.convertClubData(*newValue)
	shotData := it.convertToShotFormat(core.BallMetrics{}, false)
	shotData.ShotDataOptions.ContainsBallData = false
	shotData.ShotDataOptions.ContainsClubData = true
	shotData.ClubData = clubData
	log.Printf("[%s] Sending club data: speed=%.1f mph, attack=%.1f°, path=%.1f°, face=%.1f°, loft=%.1f°, lie=%.1f°, closure=%.1f, vImpact=%.1f, hImpact=%.1f",
		it.Name(), clubData.Speed, clubData.AngleOfAttack, clubData.Path, clubData.FaceToTarget,
		clubData.Loft, clubData.Lie, clubData.ClosureRate, clubData.VerticalFaceImpact, clubData.HorizontalFaceImpact)
	if err := it.sendData(shotData); err != nil {
		log.Printf("[%s] Error sending club data: %v", it.Name(), err)
	}
}
