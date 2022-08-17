package brain

type Protector interface {
	Check(chatId int64, text string) (forbidden bool, response string)
}
