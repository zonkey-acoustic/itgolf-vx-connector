# ProTee VX Data Mapping

## Ball Data — all 7 fields sent

| ProTee | Internal | IT Field | Status |
|---|---|---|---|
| BallData.Speed | BallSpeedMPS (converted mph→m/s) | Speed (converted back to mph) | Sent |
| BallData.LaunchAngle | VerticalAngle | VLA | Sent |
| BallData.LaunchDirection | HorizontalAngle | HLA | Sent |
| BallData.TotalSpin | TotalspinRPM | TotalSpin | Sent |
| BallData.BackSpin | BackspinRPM | BackSpin | Sent |
| BallData.SideSpin | SidespinRPM | SideSpin | Sent |
| BallData.SpinAxis | SpinAxis | SpinAxis | Sent |

## Club Data — all 9 fields sent

| ProTee | Internal | IT Field | Status |
|---|---|---|---|
| ClubData.Speed | ClubSpeed | Speed | Sent |
| ClubData.AttackAngle | AttackAngle | AngleOfAttack | Sent |
| ClubData.FaceAngle | FaceAngle | FaceToTarget | Sent |
| ClubData.Loft | DynamicLoftAngle | Loft | Sent |
| ClubData.SwingPath | PathAngle | Path | Sent |
| ClubData.Speed | ClubSpeed | SpeedAtImpact | Sent |
| ClubData.Lie | Lie | Lie | Sent |
| ClubData.ClosureRate | ClosureRate | ClosureRate | Sent |
| ClubData.ImpactPointX | ImpactPointX | HorizontalFaceImpact | Sent |
| ClubData.ImpactPointY | ImpactPointY | VerticalFaceImpact | Sent |

## Data Not Forwarded

ProTee also provides `FlightData` (carry, total distance, offline, max height, flight time) and `PhysicsSettings` (altitude, temperature, humidity, air density). These are not part of the GSPro connector protocol that Infinite Tees uses — the simulator calculates its own flight from the ball/club metrics.
