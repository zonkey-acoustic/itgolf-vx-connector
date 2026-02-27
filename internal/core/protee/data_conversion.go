package protee

import (
	"log"

	"github.com/brentyates/squaregolf-connector/internal/core"
)

const mphToMps = 1.0 / 2.23694

// convertBallMetrics converts ProTee shot data to internal BallMetrics
func convertBallMetrics(shot *ShotData) *core.BallMetrics {
	metrics := &core.BallMetrics{
		ShotType: core.ShotTypeFull,
	}

	// Ball speed: mph → m/s
	if speed, err := parseFloat(shot.BallData.Speed); err == nil {
		metrics.BallSpeedMPS = speed * mphToMps
	} else {
		log.Printf("ProTee: failed to parse ball speed %q: %v", shot.BallData.Speed, err)
	}

	// Launch angle → VerticalAngle
	if angle, err := parseFloat(shot.BallData.LaunchAngle); err == nil {
		metrics.VerticalAngle = angle
	} else {
		log.Printf("ProTee: failed to parse launch angle %q: %v", shot.BallData.LaunchAngle, err)
	}

	// Launch direction → HorizontalAngle
	if angle, err := parseFloat(shot.BallData.LaunchDirection); err == nil {
		metrics.HorizontalAngle = angle
	} else {
		log.Printf("ProTee: failed to parse launch direction %q: %v", shot.BallData.LaunchDirection, err)
	}

	// Total spin
	if spin, err := parseInt(shot.BallData.TotalSpin); err == nil {
		metrics.TotalspinRPM = spin
	} else {
		log.Printf("ProTee: failed to parse total spin %q: %v", shot.BallData.TotalSpin, err)
	}

	// Spin axis (negate: ProTee uses GSPro convention, internal uses BLE convention)
	if axis, err := parseFloat(shot.BallData.SpinAxis); err == nil {
		metrics.SpinAxis = -axis
	} else {
		log.Printf("ProTee: failed to parse spin axis %q: %v", shot.BallData.SpinAxis, err)
	}

	// Back spin
	if spin, err := parseInt(shot.BallData.BackSpin); err == nil {
		metrics.BackspinRPM = spin
	} else {
		log.Printf("ProTee: failed to parse back spin %q: %v", shot.BallData.BackSpin, err)
	}

	// Side spin (negate: ProTee uses GSPro convention, internal uses BLE convention)
	if spin, err := parseInt(shot.BallData.SideSpin); err == nil {
		metrics.SidespinRPM = -spin
	} else {
		log.Printf("ProTee: failed to parse side spin %q: %v", shot.BallData.SideSpin, err)
	}

	return metrics
}

// convertClubMetrics converts ProTee shot data to internal ClubMetrics
func convertClubMetrics(shot *ShotData) *core.ClubMetrics {
	metrics := &core.ClubMetrics{}

	// Swing path → PathAngle
	if angle, err := parseFloat(shot.ClubData.SwingPath); err == nil {
		metrics.PathAngle = angle
	} else {
		log.Printf("ProTee: failed to parse swing path %q: %v", shot.ClubData.SwingPath, err)
	}

	// Face angle
	if angle, err := parseFloat(shot.ClubData.FaceAngle); err == nil {
		metrics.FaceAngle = angle
	} else {
		log.Printf("ProTee: failed to parse face angle %q: %v", shot.ClubData.FaceAngle, err)
	}

	// Attack angle
	if angle, err := parseFloat(shot.ClubData.AttackAngle); err == nil {
		metrics.AttackAngle = angle
	} else {
		log.Printf("ProTee: failed to parse attack angle %q: %v", shot.ClubData.AttackAngle, err)
	}

	// Loft → DynamicLoftAngle
	if angle, err := parseFloat(shot.ClubData.Loft); err == nil {
		metrics.DynamicLoftAngle = angle
	} else {
		log.Printf("ProTee: failed to parse loft %q: %v", shot.ClubData.Loft, err)
	}

	// Club speed (mph, stored as-is for downstream conversion)
	if speed, err := parseFloat(shot.ClubData.Speed); err == nil {
		metrics.ClubSpeed = speed
	} else {
		log.Printf("ProTee: failed to parse club speed %q: %v", shot.ClubData.Speed, err)
	}

	// Lie
	if angle, err := parseFloat(shot.ClubData.Lie); err == nil {
		metrics.Lie = angle
	} else {
		log.Printf("ProTee: failed to parse lie %q: %v", shot.ClubData.Lie, err)
	}

	// Closure rate
	if rate, err := parseFloat(shot.ClubData.ClosureRate); err == nil {
		metrics.ClosureRate = rate
	} else {
		log.Printf("ProTee: failed to parse closure rate %q: %v", shot.ClubData.ClosureRate, err)
	}

	// Impact point X
	if x, err := parseFloat(shot.ClubData.ImpactPointX); err == nil {
		metrics.ImpactPointX = x
	} else {
		log.Printf("ProTee: failed to parse impact point X %q: %v", shot.ClubData.ImpactPointX, err)
	}

	// Impact point Y
	if y, err := parseFloat(shot.ClubData.ImpactPointY); err == nil {
		metrics.ImpactPointY = y
	} else {
		log.Printf("ProTee: failed to parse impact point Y %q: %v", shot.ClubData.ImpactPointY, err)
	}

	return metrics
}
