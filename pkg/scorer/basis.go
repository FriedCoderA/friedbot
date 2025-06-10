package score

import (
	"math/rand"

	"friedbot/pkg/config"
	"friedbot/pkg/models/schema"
)

type randomScorer struct {
	minScore int
	maxScore int
}

func (t *randomScorer) score(session *schema.Session, score int) int {
	return rand.Intn(t.maxScore-t.minScore) + t.minScore
}

type temperatureTrigger struct{}

func (t *temperatureTrigger) score(session *schema.Session, score int) int {
	temperature := config.GetTriggerSettings().Temperature
	return int(temperature * float64(score))
}
