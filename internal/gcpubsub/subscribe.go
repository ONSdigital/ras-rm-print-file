package gcpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Subscriber struct {
	Printer pkg.Printer
}

func (s Subscriber) Start() {
	log.Debug("starting worker process")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, viper.GetString("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	log.Debug("about to subscribe")
	s.subscribe(ctx, client)
}

func (s Subscriber) subscribe(ctx context.Context, client *pubsub.Client) {
	subId := viper.GetString("PUBSUB_SUB_ID")
	log.WithField("subId", subId).Info("subscribing to subscription")
	sub := client.Subscription(subId)
	cctx, cancel := context.WithCancel(ctx)
	log.Debug("waiting to receive")
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Info("print file received - processing")
		if msg.DeliveryAttempt != nil {
			log.WithField("delivery attempts", *msg.DeliveryAttempt).Info("Message delivery attempted")
		}

		dataFileName := string(msg.Data)
		attribute := msg.Attributes
		printFilename, ok := attribute["printFilename"]

		if ok {
			log.WithField("printFilename", printFilename).Info("about to process print file")
			err := s.Printer.Process(printFilename, dataFileName)
			if err != nil {
				log.WithError(err).Error("error processing printfile - nacking message")
				//after x number of nacks message will be DLQ
				msg.Nack()
			} else {
				log.Info("print file processed - acking message")
				msg.Ack()
			}
		} else {
			log.Error("missing printFilename - sending to DLQ")
			err := deadLetter(ctx, client, msg)
			if err != nil {
				msg.Nack()
			}
		}
	})

	if err != nil {
		log.WithError(err).Error("error subscribing")
		cancel()
	}
}

// send message to DLQ immediately
func deadLetter(ctx context.Context, client *pubsub.Client, msg *pubsub.Message) error {
	//DLQ are always named TOPIC + -dead-letter in our terraform scripts
	deadLetterTopic := viper.GetString("PUB_SUB_TOPIC") + "-dead-letter"
	dlq := client.Topic(deadLetterTopic)
	id, err := dlq.Publish(ctx, msg).Get(ctx)
	if err != nil {
		log.WithField("msg", string(msg.Data)).WithError(err).Error("unable to forward to dead letter topic")
		return err
	}
	log.WithField("id", id).Info("published to dead letter topic")
	return nil
}
