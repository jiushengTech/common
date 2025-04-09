package rocketmq

import (
	"github.com/jiushengTech/common/broker"

	rocketmqOption "github.com/jiushengTech/common/broker/rocketmq/option"

	aliyunMQ "github.com/jiushengTech/common/broker/rocketmq/aliyun"
	rocketmqV2 "github.com/jiushengTech/common/broker/rocketmq/rocketmq-client-go"
	rocketmqV5 "github.com/jiushengTech/common/broker/rocketmq/rocketmq-clients"
)

func NewBroker(driverType rocketmqOption.DriverType, opts ...broker.Option) broker.Broker {

	switch driverType {
	case rocketmqOption.DriverTypeAliyun:
		return aliyunMQ.NewBroker(opts...)
	case rocketmqOption.DriverTypeV2:
		return rocketmqV2.NewBroker(opts...)
	case rocketmqOption.DriverTypeV5:
		return rocketmqV5.NewBroker(opts...)
	default:
		return nil
	}
}
