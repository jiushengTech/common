package rocketmqClientGo

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/jiushengTech/common/broker"
)

type publication struct {
	topic  string
	err    error
	m      *broker.Message
	ctx    context.Context
	reader rocketmq.PushConsumer
	rm     *primitive.Message
}

func (p *publication) Topic() string {
	return p.topic
}

func (p *publication) Message() *broker.Message {
	return p.m
}

func (p *publication) RawMessage() interface{} {
	return p.rm
}

func (p *publication) Ack() error {
	return nil
}

func (p *publication) Error() error {
	return p.err
}
