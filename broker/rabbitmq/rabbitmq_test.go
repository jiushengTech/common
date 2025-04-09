package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/stretchr/testify/assert"

	"github.com/jiushengTech/common/broker"
	api "github.com/jiushengTech/common/testing/api/manual"
	"github.com/jiushengTech/common/tracing"
)

const (
	testBroker = "amqp://user:bitnami@127.0.0.1:5672"

	testExchange = "test_exchange"
	testQueue    = "test_queue"
	testRouting  = "test_routing_key"
)

func handleHygrothermograph(_ context.Context, topic string, headers broker.Headers, msg *api.Hygrothermograph) error {
	log.Infof("Topic %s, Headers: %+v, Payload: %+v\n", topic, headers, msg)
	return nil
}

func Test_Publish_WithRawData(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx := context.Background()

	b := NewBroker(
		broker.WithOptionContext(ctx),
		broker.WithAddress(testBroker),
		WithExchangeName(testExchange),
		WithDurableExchange(),
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
		msg.Humidity = float64(i)
		msg.Temperature = float64(rand.Intn(100))
		buf, _ := json.Marshal(&msg)
		err := b.Publish(ctx, testRouting, buf)
		assert.Nil(t, err)
		elapsedTime := time.Since(startTime) / time.Millisecond
		fmt.Printf("Publish %d, elapsed time: %dms, Humidity: %.2f Temperature: %.2f\n",
			i, elapsedTime, msg.Humidity, msg.Temperature)
	}

	fmt.Printf("total send %d messages\n", count)

	<-interrupt
}

func Test_Subscribe_WithRawData(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	b := NewBroker(
		broker.WithAddress(testBroker),
		WithExchangeName(testExchange),
		WithDurableExchange(),
	)

	_ = b.Init()

	if err := b.Connect(); err != nil {
		t.Logf("cant connect to broker, skip: %v", err)
		t.Skip()
	}
	defer b.Disconnect()

	_, err := b.Subscribe(testRouting,
		api.RegisterHygrothermographRawHandler(handleHygrothermograph),
		nil,
		broker.WithQueueName(testQueue),
		// broker.WithDisableAutoAck(),
		WithDurableQueue(),
	)
	assert.Nil(t, err)

	<-interrupt
}

func Test_Publish_WithJsonCodec(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx := context.Background()

	b := NewBroker(
		broker.WithAddress(testBroker),
		broker.WithCodec("json"),
		WithExchangeName(testExchange),
		WithDurableExchange(),
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
		msg.Humidity = float64(i)
		msg.Temperature = float64(rand.Intn(100))
		err := b.Publish(ctx, testRouting, msg)
		assert.Nil(t, err)
		elapsedTime := time.Since(startTime) / time.Millisecond
		fmt.Printf("Publish %d, elapsed time: %dms, Humidity: %.2f Temperature: %.2f\n",
			i, elapsedTime, msg.Humidity, msg.Temperature)
	}

	fmt.Printf("total send %d messages\n", count)

	<-interrupt
}

func Test_Subscribe_WithJsonCodec(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	b := NewBroker(
		broker.WithAddress(testBroker),
		broker.WithCodec("json"),
		WithExchangeName(testExchange),
		WithDurableExchange(),
	)

	_ = b.Init()

	if err := b.Connect(); err != nil {
		t.Logf("cant connect to broker, skip: %v", err)
		t.Skip()
	}
	defer b.Disconnect()

	_, err := b.Subscribe(testRouting,
		api.RegisterHygrothermographJsonHandler(handleHygrothermograph),
		api.HygrothermographCreator,
		broker.WithQueueName(testQueue),
		// broker.WithDisableAutoAck(),
		WithDurableQueue(),
	)
	assert.Nil(t, err)

	<-interrupt
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
			"rabbitmq-tracer",
		)
	case "zipkin":
		return broker.WithTracerProvider(tracing.NewTracerProvider(exporterName,
			"http://localhost:9411/api/v2/spans",
			serviceName,
			"test",
			"1.0.0",
			1.0,
		),
			"rabbitmq-tracer",
		)
	}

	return nil
}

func Test_Publish_WithTracer(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx := context.Background()

	b := NewBroker(
		broker.WithOptionContext(ctx),
		broker.WithAddress(testBroker),
		broker.WithCodec("json"),
		WithExchangeName(testExchange),
		WithDurableExchange(),
		createTracerProvider("otlp-grpc", "publish_tracer_tester"),
		//broker.WithPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})),
	)

	_ = b.Init()

	if err := b.Connect(); err != nil {
		t.Logf("cant connect to broker, skip: %v", err)
		t.Skip()
	}
	defer b.Disconnect()

	var msg api.Hygrothermograph
	const count = 1
	for i := 0; i < count; i++ {
		startTime := time.Now()
		msg.Humidity = float64(rand.Intn(100))
		msg.Temperature = float64(rand.Intn(100))
		err := b.Publish(ctx, testRouting, msg, WithMessageId(strconv.Itoa(i)))
		assert.Nil(t, err)
		elapsedTime := time.Since(startTime) / time.Millisecond
		fmt.Printf("Publish %d, elapsed time: %dms, Humidity: %.2f Temperature: %.2f\n",
			i, elapsedTime, msg.Humidity, msg.Temperature)
	}

	fmt.Printf("total send %d messages\n", count)

	<-interrupt
}

func Test_Subscribe_WithTracer(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	b := NewBroker(
		broker.WithAddress(testBroker),
		broker.WithCodec("json"),
		createTracerProvider("otlp-grpc", "subscribe_tracer_tester"),
		//broker.WithPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})),
		WithExchangeName(testExchange),
		WithDurableExchange(),
	)

	_ = b.Init()

	if err := b.Connect(); err != nil {
		t.Logf("cant connect to broker, skip: %v", err)
		t.Skip()
	}
	defer b.Disconnect()

	_, err := b.Subscribe(testRouting,
		api.RegisterHygrothermographJsonHandler(handleHygrothermograph),
		api.HygrothermographCreator,
		broker.WithQueueName(testQueue),
		// broker.WithDisableAutoAck(),
		WithDurableQueue(),
	)
	assert.Nil(t, err)

	<-interrupt
}
