package telegram

type WebhookRequest struct {
	Message *WebhookRequestMessage `json:"message"`
}

type WebhookRequestMessage struct {
	MessageId int64                        `json:"message_id"`
	From      *WebhookRequestMessageSender `json:"from"`
	Chat      *WebhookRequestMessageChat   `json:"chat"`
	Text      string                       `json:"text"`
}

type WebhookRequestMessageSender struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type WebhookRequestMessageChat struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}
