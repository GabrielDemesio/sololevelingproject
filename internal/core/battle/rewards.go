package battle

func RankMultiplier(rank string) float64 {
	switch rank {
	case "S":
		return 2.2
	case "A":
		return 1.8
	case "B":
		return 1.5
	case "C":
		return 1.3
	case "D":
		return 1.1
	case "E":
		fallthrough
	default:
		return 1.0
	}
}
func ComputeRewards(minutes int, dungeonMultiplier float64, questWeight int, quality float64) (xp, gold int64) {
	if minutes < 0 {
		minutes = 0
	}
	base := 10.0
	qw := questWeight
	if qw <= 0 {
		qw = 1
	}
	if quality < 0.5 {
		quality = 0.5
	}
	if quality > 1.5 {
		quality = 1.5
	}
	xpF := base * float64(minutes) * dungeonMultiplier * float64(qw) * quality
	if xpF < 0 {
		xpF = 0
	}
	return int64(xpF), int64(xpF / 2.0)
}
