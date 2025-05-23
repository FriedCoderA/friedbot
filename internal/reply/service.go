package reply

import (
	"encoding/json"
	"errors"
	"log/slog"

	"friedbot/pkg/models/dao"
	"friedbot/pkg/models/schema"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/spf13/viper"
)

const (
	MaxMessageSize = 1024 * 1024
)

type Service struct {
	engine *gin.Engine
}

func NewService() *Service {
	s := &Service{
		engine: gin.Default(),
	}
	s.Router()
	return s
}

func (s *Service) Start() error {
	addr := viper.GetString("server.address")
	if addr == "" {
		return errors.New("server.address is empty")
	}
	return s.engine.Run(addr)
}

func (s *Service) Router() {
	m := melody.New()
	m.Config.MaxMessageSize = MaxMessageSize
	s.engine.GET("/message", func(c *gin.Context) {
		err := m.HandleRequest(c.Writer, c.Request)
		if err != nil {
			slog.Error("HandleRequest error", "error", err)
		}
	})
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
	HandelMessage(m)
}

func HandelMessage(m *melody.Melody) {
	sessions := dao.NewSessionManager()
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
		session, err := sessions.GetOrCreate(&msg)
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
