package mail

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/mailgun/mailgun-go/v4"
)

func SendMail(sender, subject, body, recipient, senderDomain, apiKey string) {
    // Create an instance of the Mailgun Client
    mg := mailgun.NewMailgun(senderDomain, apiKey)

    // The message object allows you to add attachments and Bcc recipients
    message := mg.NewMessage(sender, subject, body, recipient)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()

    // Send the message with a 10 second timeout
    resp, id, err := mg.Send(ctx, message)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
