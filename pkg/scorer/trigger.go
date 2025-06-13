package score

import (
	"log/slog"
	"time"

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
	&randomScorer{
		maxScore: 20,
		minScore: 0,
	},
	&atScorer{},
	&privateScorer{},
	&aiScorer{
		msgLoadCount:      8,
		msgExpireDuration: time.Minute * 5,
	},
	&temperatureTrigger{},
}

func Trigger(session *schema.Session) bool {
	score := defaultScore
	for _, scorer := range InstalledScorers {
		score = scorer.score(session, score)
		if score >= targetScore || score <= deadScore {
			break
		}
	}
	slog.Debug("trigger", "score", score)
	return score >= targetScore
}
