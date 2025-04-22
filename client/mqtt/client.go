package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	klog "github.com/jiushengTech/common/log/klog/logger"
	"sync"
)

type Client struct {
	mqtt.Client
	subscriptions *sync.Map // 使用sync.Map来替代map
}

func NewClient(opts *mqtt.ClientOptions) *Client {
	srv := &Client{
		Client:        mqtt.NewClient(opts),
		subscriptions: &sync.Map{},
	}
	return srv
}

func (c *Client) GetSubscriptions() *sync.Map {
	return c.subscriptions
}

// Subscribe 订阅主题，并保存订阅信息以便断线重连时恢复
func (c *Client) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	// 订阅主题
	token := c.Client.Subscribe(topic, qos, callback)
	if token.Wait() && token.Error() != nil {
		return token.Error() // 返回订阅错误
	}

	// 使用sync.Map存储订阅信息，避免频繁创建新的map
	subscriptions, _ := c.subscriptions.LoadOrStore(c.Client.OptionsReader().ClientID, &sync.Map{})
	subscriptions.(*sync.Map).Store(topic, qos)

	klog.Log.Infof("成功订阅主题: %s (QoS: %d)", topic, qos)
	return nil
}

// Unsubscribe 取消订阅主题，并从保存的订阅信息中移除
func (c *Client) Unsubscribe(topic string) error {
	// 取消订阅主题
	token := c.Client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		return token.Error() // 返回取消订阅错误
	}

	// 使用sync.Map移除订阅信息
	subscriptions, ok := c.subscriptions.Load(c.Client.OptionsReader().ClientID)
	if ok {
		subscriptions.(*sync.Map).Delete(topic)
	}

	klog.Log.Infof("成功取消订阅主题: %s", topic)
	return nil
}

// Publish 发布消息
func (c *Client) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	// 发布消息
	token := c.Client.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error() // 返回发布消息错误
	}

	klog.Log.Infof("成功发布消息至主题: %s", topic)
	return nil
}
