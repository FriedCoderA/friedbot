package triggers

import (
	"friedbot/pkg/models/schema"
)

const (
	targetScore  = 60
	defaultScore = 0
	deadScore    = -100
)

type Scorer interface {
	score(session *schema.Session, score int) int
}

var InstalledScorers = []Scorer{
	&RandomScorer{
		MaxScore: 20,
		MinScore: -20,
	},
	&TemperatureTrigger{},
	&AIScorer{},
}

func Trigger(session *schema.Session) bool {
	score := defaultScore
	for _, scorer := range InstalledScorers {
		score = scorer.score(session, score)
		if score >= targetScore || score <= deadScore {
			break
		}
	}
	return score >= targetScore
}
