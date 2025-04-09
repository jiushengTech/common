package aliyun

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"

	"github.com/jiushengTech/common/broker"
	rocketmqOption "github.com/jiushengTech/common/broker/rocketmq/option"
	api "github.com/jiushengTech/common/testing/api/manual"
	"github.com/jiushengTech/common/tracing"
)

const (
	testBroker    = "127.0.0.1:9876"
	testTopic     = "test_topic"
	testGroupName = "CID_ONSAPI_OWNER"
)

func handleHygrothermograph(_ context.Context, topic string, headers broker.Headers, msg *api.Hygrothermograph) error {
	log.Infof("Topic %s, Headers: %+v, Payload: %+v\n", topic, headers, msg)
	return nil
}

func createTracerProvider(exporterName, serviceName string) broker.Option {
	switch exporterName {
	case "otlp-grpc":
		return broker.WithTracerProvider(tracing.NewTracerProvider(exporterName,
			"localhost:4317",
			serviceName,
			"",
			"1.0.0",
			1.0,
		),
			"rocketmq-tracer",
		)
	case "zipkin":
		return broker.WithTracerProvider(tracing.NewTracerProvider(exporterName,
			"http://localhost:9411/api/v2/spans",
			serviceName,
			"test",
			"1.0.0",
			1.0,
		),
			"rocketmq-tracer",
		)
	}

	return nil
}

func Test_Aliyun_Publish(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx := context.Background()

	endpoint := ""
	accessKey := ""
	secretKey := ""
	instanceId := ""
	topicName := ""

	b := NewBroker(
		broker.WithCodec("json"),
		rocketmqOption.WithEnableTrace(),
		rocketmqOption.WithNameServerDomain(endpoint),
		rocketmqOption.WithAccessKey(accessKey),
		rocketmqOption.WithSecretKey(secretKey),
		rocketmqOption.WithInstanceName(instanceId),
	)

	_ = b.Init()

	if err := b.Connect(); err != nil {
		t.Logf("cant connect to broker, skip: %v", err)
		t.Skip()
	}
	defer b.Disconnect()

	var msg api.Hygrothermograph
	const count = 10
	for i := 0; i < count; i++ {
		startTime := time.Now()
		msg.Humidity = float64(rand.Intn(100))
		msg.Temperature = float64(rand.Intn(100))
		err := b.Publish(ctx, topicName, msg)
		assert.Nil(t, err)
		elapsedTime := time.Since(startTime) / time.Millisecond
		fmt.Printf("Publish %d, elapsed time: %dms, Humidity: %.2f Temperature: %.2f\n",
			i, elapsedTime, msg.Humidity, msg.Temperature)
	}

	fmt.Printf("total send %d messages\n", count)

	<-interrupt
}

func Test_Aliyun_Subscribe(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	endpoint := ""
	accessKey := ""
	secretKey := ""
	instanceId := ""
	topicName := ""
	groupName := "GID_DEFAULT"

	b := NewBroker(
		broker.WithCodec("json"),
		rocketmqOption.WithEnableTrace(),
		rocketmqOption.WithNameServerDomain(endpoint),
		rocketmqOption.WithAccessKey(accessKey),
		rocketmqOption.WithSecretKey(secretKey),
		rocketmqOption.WithInstanceName(instanceId),
	)

	_ = b.Init()

	if err := b.Connect(); err != nil {
		t.Logf("cant connect to broker, skip: %v", err)
		t.Skip()
	}
	defer b.Disconnect()

	_, err := b.Subscribe(topicName,
		api.RegisterHygrothermographJsonHandler(handleHygrothermograph),
		api.HygrothermographCreator,
		broker.WithQueueName(groupName),
	)
	assert.Nil(t, err)

	<-interrupt
}
