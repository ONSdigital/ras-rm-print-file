package gcpubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	logger "github.com/ONSdigital/ras-rm-print-file/logging"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Subscriber struct {
	Printer pkg.Printer
}

func (s Subscriber) Start() {
	logger.Debug("starting worker process")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, viper.GetString("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
		logger.Fatal("failed to start worker process", zap.Error(err))
	}
	defer client.Close()
	logger.Debug("about to subscribe")
	s.subscribe(ctx, client)
}

func (s Subscriber) subscribe(ctx context.Context, client *pubsub.Client) {
	subId := viper.GetString("PUBSUB_SUB_ID")
	logger.Info("subscribing to subscription", zap.String("subId", subId))
	sub := client.Subscription(subId)
	cctx, cancel := context.WithCancel(ctx)
	logger.Debug("waiting to receive")
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		logger.Info("print file received - processing")
		if msg.DeliveryAttempt != nil {
			// TODO: Handle DeliverAttempt *int in Zap logger:
			log.WithField("delivery attempts", *msg.DeliveryAttempt).Info("Message delivery attempted")
		}

		dataFileName := string(msg.Data)
		attribute := msg.Attributes
		printFilename, ok := attribute["printFilename"]

		if ok {
			logger.Info("about to process print file", zap.String("printFilename", printFilename))
			err := s.Printer.Process(printFilename, dataFileName)
			if err != nil {
				logger.Error("error processing printfile - nacking message", zap.Error(err))
				// after x number of nacks message will be DLQ
				msg.Nack()
			} else {
				logger.Info("print file processed - acking message")
				msg.Ack()
			}
		} else {
			logger.Error("missing printFilename - sending to DLQ")
			err := deadLetter(ctx, client, msg)
			if err != nil {
				msg.Nack()
			}
		}
	})

	if err != nil {
		logger.Error("error subscribing", zap.Error(err))
		cancel()
	}
}

// send message to DLQ immediately
func deadLetter(ctx context.Context, client *pubsub.Client, msg *pubsub.Message) error {
	// DLQ are always named TOPIC + -dead-letter in our terraform scripts
	deadLetterTopic := viper.GetString("PUB_SUB_TOPIC") + "-dead-letter"
	dlq := client.Topic(deadLetterTopic)
	id, err := dlq.Publish(ctx, msg).Get(ctx)
	if err != nil {
		// TODO: Handle msg *Message in Zap logger:
		log.WithField("msg", string(msg.Data)).WithError(err).Error("unable to forward to dead letter topic")
		return err
	}
	logger.Info("published to dead letter topic", zap.String("id", id))
	return nil
}
