package processor

import (
	"github.com/go-redis/redis"
	"log"
	"net/smtp"
	"user_mailer/model"
)
var err error

func SubRedis(redisClient *redis.Client,conf model.Config)  error {

	pubsub := redisClient.Subscribe(conf.RedisChannel)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err = pubsub.Receive()
	if err != nil {
		return err
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()

	// Consume messages.Infinite loop for listeninig chanel
	for msg := range ch {

		err = sendMail(msg.Payload,conf)
		if err != nil {
			return  err
		}
	}

	return nil
}

func sendMail(body string,conf model.Config) error {

	from := conf.MailUsername
	pass := conf.MailPassword
	to := "feriz2013@hotmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Confirmation email\n\n" +
		body

	//send message
	err = smtp.SendMail(conf.MailHost+":"+conf.MailPort,
		smtp.PlainAuth("", from, pass, conf.MailHost),
		from, []string{to}, []byte(msg))

	if err != nil {
		return err
	}

	log.Print("Successfully sent email!")
	return nil
}

//https://godoc.org/github.com/go-redis/redis#example-PubSub
//https://gist.github.com/jpillora/cb46d183eca0710d909a