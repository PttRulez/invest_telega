package telega

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pttrulez/invest_telega/internal/grpctransport"
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

func New(botToken string, investorEndpoint string, logger *logger.Logger) (*Service, error) {
	b, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return &Service{}, err
	}

	investorGrpcClient, err := grpctransport.NewInvestorGRPCClient(investorEndpoint)
	if err != nil {
		return &Service{}, err
	}

	b.Handle("/getid", func(c tele.Context) error {
		return c.Send(fmt.Sprintf("Ваш ID: %d", c.Chat().ID))
	})

	b.Handle("/list", func(c tele.Context) error {
		fmt.Println("LIST")
		chatId := c.Chat().ID
		portfolios, err := investorGrpcClient.GetPortfolioList(context.Background(),
			strconv.FormatInt(chatId, 10))
		if err != nil {
			return c.Send(fmt.Sprintf(`Не удалось получить список портфолио.\n
			Возможно вы не привязали на сайте ваш чат id.\n
			Ваш чат id: %d`, c.Chat().ID))
		}

		if len(portfolios) == 0 {
			return c.Send(" У вас нет ни одного портфолио")
		}

		buttons := make([]tele.Btn, len(portfolios))
		r := b.NewMarkup()

		for _, p := range portfolios {
			fmt.Println("p.GetId()", strconv.FormatInt(p.GetId(), 10))
			b := tele.Btn{
				Unique: "portofliosummary",
				Text:   p.GetName(),
				Data:   string(strconv.FormatInt(p.GetId(), 10)),
			}
			buttons = append(buttons, b)
		}

		r.Inline(r.Row(buttons...))

		return c.Send("Ваши портфолио", r)
	})

	b.Handle(&tele.Btn{Unique: "portofliosummary"}, func(c tele.Context) error {
		portfolioID, err := strconv.Atoi(c.Data())
		if err != nil {
			return c.Send("Не удалось получить портфолио")
		}
		msg, err := investorGrpcClient.GetPortfolioSummaryMessage(context.Background(), portfolioID,
			strconv.FormatInt(c.Chat().ID, 10))
		if err != nil {
			return c.Send("Не удалось получить портфолио")
		}

		fmt.Println("portfolioID", portfolioID)
		fmt.Println("chatID", strconv.FormatInt(c.Chat().ID, 10))

		return c.Send(msg)
	})

	go func() {
		b.Start()
	}()

	return &Service{bot: b, logger: logger}, nil
}

func (s *Service) Close() {
	s.bot.Stop()
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
