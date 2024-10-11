package telega

import (
	"fmt"
	"time"

	"github.com/pttrulez/invest_telega/pkg/protogen"

	"github.com/pttrulez/invest_telega/pkg/logger"

	tele "gopkg.in/telebot.v3"
)

func (s *Service) SendMsg(msgInfo *protogen.MessageInfo) error {
	s.logger.Debug(fmt.Sprintf("ChatID: %s", msgInfo.GetChatId()))
	s.logger.Debug(fmt.Sprintf("Sent message text: %s", msgInfo.GetText()))

	u := &User{ID: msgInfo.GetChatId()}
	_, err := s.bot.Send(u, msgInfo.GetText())
	if err != nil {
		return fmt.Errorf("TelegaService.SendMsg: %w", err)
	}

	return nil
}

func New(botToken string, logger *logger.Logger) (*Service, error) {
	b, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return &Service{}, err
	}

	b.Handle("/getid", func(c tele.Context) error {
		return c.Send(fmt.Sprintf("Ваш ID: %d", c.Chat().ID))
	})

	go func() {
		b.Start()
	}()

	return &Service{bot: b, logger: logger}, nil
}

type Service struct {
	bot    *tele.Bot
	logger *logger.Logger
}

type User struct {
	ID string
}

func (u *User) Recipient() string {
	return u.ID
}
