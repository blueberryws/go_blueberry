package go_blueberry

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/mailgun/mailgun-go/v4"
)

type Mail struct {
    Sender string
    SenderDomain string
    ApiKey string
    Recipient string
}

func (m Mail) SendMail(subject, body string) {
    // Create an instance of the Mailgun Client
    mg := mailgun.NewMailgun(m.SenderDomain, m.ApiKey)

    // The message object allows you to add attachments and Bcc recipients
    message := mg.NewMessage(m.Sender, subject, body, m.Recipient)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()

    // Send the message with a 10 second timeout
    resp, id, err := mg.Send(ctx, message)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
