package triggers

import (
	"math/rand"

	"friedbot/pkg/config"
	"friedbot/pkg/models/schema"
)

type RandomTrigger struct {
	MinScore int
	MaxScore int
}

func (t *RandomTrigger) score(session *schema.Session, score int) int {
	return rand.Intn(t.MaxScore-t.MinScore) + t.MinScore
}

type TemperatureTrigger struct{}

func (t *TemperatureTrigger) score(session *schema.Session, score int) int {
	temperature := config.GetTriggerSettings().Temperature
	return int(temperature * float64(score))
}
