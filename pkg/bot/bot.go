package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotClient struct {
	TgBot     *tgbotapi.BotAPI
	BotToken  string
	ChannelID string
}

func NewBotClient(tgBot *tgbotapi.BotAPI, botToken, channelID string) *BotClient {
	return &BotClient{TgBot: tgBot, BotToken: botToken, ChannelID: channelID}
}

func (b *BotClient) SendMessageToChannel(message string) error {
	msg := tgbotapi.NewMessageToChannel(b.ChannelID, message)
	_, err := b.TgBot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
