package bot

type TgBot struct {
	token string
	proxy string
}

const tgBotName = "telegram bot"

func (b *TgBot) SendMsg(msg string) error {
	return nil
}

func (b *TgBot) GetBotName() string {
	return tgBotName
}
