package models

import (
	"encoding/json"

	"github.com/otwdev/galaxylib"
	nsq "github.com/segmentio/nsq-go"
)

type Notice struct {
	Body string
}

type MailInfo struct {
	Type string `json:"Type"`
	Data struct {
		Title string `json:"Title"`
		Body  string `json:"Body"`
		To    string `json:"To"`
	} `json:"Data"`
}

func (n *Notice) Send() {
	mail := galaxylib.GalaxyCfgFile.MustValue("data", "mail")

	nsqAddr := galaxylib.GalaxyCfgFile.MustValue("data", "nsq")

	topic := galaxylib.GalaxyCfgFile.MustValue("data", "topic")

	producer, _ := nsq.StartProducer(nsq.ProducerConfig{
		Topic:   topic,
		Address: nsqAddr,
	})

	mailInfo := &MailInfo{
		Type: "email",
	}

	mailInfo.Data.Body = n.Body
	mailInfo.Data.Title = "汇率更新提醒"
	mailInfo.Data.To = mail

	buf, _ := json.Marshal(mailInfo)

	// Publishes a message to the topic that this producer is configured for,
	// the method returns when the operation completes, potentially returning an
	// error if something went wrong.
	producer.Publish(buf)

	// Stops the producer, all in-flight requests will be canceled and no more
	// messages can be published through this producer.
	producer.Stop()

}
