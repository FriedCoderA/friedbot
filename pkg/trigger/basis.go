package trigger

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strings"

	"friedbot/pkg/config"
	"friedbot/pkg/kinds"
	"friedbot/pkg/models/dao"
	"friedbot/pkg/models/schema"
)

type randomScorer struct {
	minScore int
	maxScore int
}

func (t *randomScorer) score(session *schema.Session, score int) int {
	return score + rand.Intn(t.maxScore-t.minScore) + t.minScore
}

type temperatureTrigger struct{}

func (t *temperatureTrigger) score(session *schema.Session, score int) int {
	temperature := config.GetTriggerSettings().Temperature
	return int(temperature * float64(score))
}

type atScorer struct{}

func (t *atScorer) score(session *schema.Session, score int) int {
	msgManager := dao.NewMessageManager(session.ID)
	messages, err := msgManager.TopN(5)
	if err != nil {
		slog.Error("at score get user messages error", "err", err)
		return score
	}
	selfID := config.GetBotSettings().QQ
	for _, message := range messages {
		if strings.Contains(message.Content, fmt.Sprintf("qq=%d", selfID)) {
			score += 100
			break
		}
	}
	return score
}

type privateScorer struct{}

func (t *privateScorer) score(session *schema.Session, score int) int {
	if session.MessageType == kinds.MessageTypePrivate {
		score += 100
	}
	return score
}
