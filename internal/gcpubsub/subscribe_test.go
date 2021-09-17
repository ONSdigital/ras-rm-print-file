package gcpubsub

import (
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"context"
	"fmt"
	"github.com/ONSdigital/ras-rm-print-file/internal/config"
	mocks "github.com/ONSdigital/ras-rm-print-file/mocks/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"os"
	"testing"
	"time"
)

var (
	printFile = "[{\n  \"sampleUnitRef\":\"10001\",\n  \"iac\":\"ai9bt497r7bn\",\n  \"caseGroupStatus\":\"NOTSTARTED\",\n  \"enrolmentStatus\":\"\",\n  \"respondentStatus\":\"\",\n  \"contact\":{\n    \"forename\":\"Jon\",\n    \"surname\":\"Snow\",\n    \"emailAddress\":\"jon.snow@example.com\"\n  },\n  \"region\":\"HH\"\n}]"
	client    *pubsub.Client
	ctx       context.Context
)

func TestMain(m *testing.M) {
	config.SetDefaults()
	//create a fake Pub Sub serer
	ctx = context.Background()
	// Start a fake server running locally.
	srv := pstest.NewServer()
	defer srv.Close()
	// Connect to the server without using TLS.
	conn, _ := grpc.Dial(srv.Addr, grpc.WithInsecure())

	defer conn.Close()
	// Use the connection when creating a pubsub client.
	client, _ = pubsub.NewClient(ctx, "rm-ras-sandbox", option.WithGRPCConn(conn))
	defer client.Close()

	os.Exit(m.Run())
}

func TestSubscribe(t *testing.T) {
	assert := assert.New(t)

	topic, err := createTopic(assert)
	defer topic.Delete(ctx)

	sub := createSubscription(err, topic, assert)
	defer sub.Delete(ctx)

	printFilename := "test.csv"

	printer := new(mocks.Printer)
	printer.On("Process", printFilename, mock.Anything).Return(nil)

	subscriber := Subscriber{
		Printer: printer,
	}

	msg := &pubsub.Message{
		Data: []byte(printFile),
		Attributes: map[string]string{
			"printFilename": printFilename,
		},
	}
	// now publish the message
	_, err = topic.Publish(ctx, msg).Get(ctx)
	assert.Nil(err)

	go subscriber.subscribe(ctx, client)

	//sleep a second for the test to complete, then allow everything to shut down
	time.Sleep(1 * time.Second)

	printer.AssertCalled(t, "Process", printFilename, mock.Anything)
}

func TestSubscribeFailsMissingFilename(t *testing.T) {
	assert := assert.New(t)

	topic, err := createTopic(assert)
	defer topic.Delete(ctx)

	sub := createSubscription(err, topic, assert)
	defer sub.Delete(ctx)

	printFilename := "test.csv"

	printer := new(mocks.Printer)
	printer.On("Process", printFilename, mock.Anything).Return(nil)

	subscriber := Subscriber{
		Printer: printer,
	}

	msg := &pubsub.Message{
		Data: []byte(printFile),
	}
	// now publish the message
	_, err = topic.Publish(ctx, msg).Get(ctx)
	assert.Nil(err)

	go subscriber.subscribe(ctx, client)

	//sleep a second for the test to complete, then allow everything to shut down
	time.Sleep(1 * time.Second)

	printer.AssertNotCalled(t, "Process", printFilename, mock.Anything)
}

func createSubscription(err error, topic *pubsub.Topic, assert *assert.Assertions) *pubsub.Subscription {
	sub, err := client.CreateSubscription(ctx, "print-file", pubsub.SubscriptionConfig{
		Topic: topic,
	})
	assert.Nil(err)
	assert.NotNil(sub)
	return sub
}

func createTopic(assert *assert.Assertions) (*pubsub.Topic, error) {
	topic, err := client.CreateTopic(ctx, "print-file")
	assert.Nil(err)
	assert.NotNil(topic)
	fmt.Println(topic)
	return topic, err
}
