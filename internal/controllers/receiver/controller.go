package receiver

import (
	"encoding/json"
	"log/slog"

	"friedbot/pkg/models/dao"
	"friedbot/pkg/models/schema"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

const (
	MaxMessageSize = 1024 * 1024 * 1024
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Router(router *gin.RouterGroup) {
	router.GET("/receive", HandelMessage)
}

func HandelMessage(c *gin.Context) {
	m := melody.New()
	m.Config.MaxMessageSize = MaxMessageSize
	err := m.HandleRequest(c.Writer, c.Request)
	if err != nil {
		slog.Error("HandleRequest error", "error", err)
	}
	m.HandleConnect(func(s *melody.Session) {
		slog.Info("connect onebot success")
	})
	m.HandleClose(func(s *melody.Session, i int, s2 string) error {
		slog.Warn("onebot connection closed")
		return nil
	})
	m.HandleError(func(s *melody.Session, err error) {
		slog.Error("onebot connection panic", "error", err)
		return
	})
	m.HandleMessage(func(s *melody.Session, bytes []byte) {
		var msg schema.Message
		err := json.Unmarshal(bytes, &msg)
		if err != nil {
			slog.Error("handle message error", "error", err)
			return
		}
		if !msg.IsAccess() {
			return
		}
		session, err := dao.NewSessionManager().GetOrCreate(&msg)
		if err != nil {
			slog.Error("get or create session error", err)
			return
		}
		messages := dao.NewMessageManager(session.ID)
		err = messages.Create(&msg)
		if err != nil {
			slog.Error("create message error", err)
			return
		}
	})
}
