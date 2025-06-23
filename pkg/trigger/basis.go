package trigger

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strings"

	"friedbot/pkg/config"
	"friedbot/pkg/events"
	"friedbot/pkg/kinds"
	"friedbot/pkg/models/dao"
)

type randomScorer struct {
	minScore int
	maxScore int
}

func (t *randomScorer) score(event *events.MessageEvent, score int) int {
	return score + rand.Intn(t.maxScore-t.minScore) + t.minScore
}

type temperatureTrigger struct{}

func (t *temperatureTrigger) score(event *events.MessageEvent, score int) int {
	temperature := config.GetTriggerSettings().Temperature
	return int(temperature * float64(score))
}

type atScorer struct{}

func (t *atScorer) score(event *events.MessageEvent, score int) int {
	msgManager := dao.NewMessageManager(event.Session.ID)
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

func (t *privateScorer) score(event *events.MessageEvent, score int) int {
	if event.Session.MessageType == kinds.MessageTypePrivate {
		score += 100
	}
	return score
}

type lengthScorer struct{}

func (t *lengthScorer) score(event *events.MessageEvent, score int) int {
	length := len(event.Message.Content)
	switch {
	case length < 2:
		score -= 100
	case 15 < length && length < 30:
		score += 10
	}
	return score
}
