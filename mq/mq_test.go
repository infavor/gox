package mq_test

import (
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/mq"
	"github.com/hetianyi/gox/timer"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

// -t 1 --redis-server=123123@192.168.245.142:6379 --mq-server=
func TestInitRabbitMQClient(t *testing.T) {
	client := &mq.RabbitClient{
		ConnectionString: "amqp://admin:123456@192.168.245.142:5672/",
		QueueNames:       []string{"url_list_queue", "html_content_queue"},
		ExchangeNames:    []string{"exchange_spider_man"},
		Bindings: []mq.Binding{
			{
				QueueName:    "url_list_queue",
				ExchangeName: "exchange_spider_man",
				BindingKey:   "binding_key_url",
			},
			{
				QueueName:    "html_content_queue",
				ExchangeName: "exchange_spider_man",
				BindingKey:   "binding_key_html",
			},
		},
		AutoACK:       false,
		ConsumeQueues: []string{"url_list_queue", "html_content_queue"},
		ConsumeHandler: func(queueName string, d amqp.Delivery) error {
			logger.Info("从队列 ", queueName, " 收到消息：", string(d.Body))
			d.Ack(false)
			return nil
		},
	}
	client.Init()
	logger.Info("启动成功...")

	timer.Start(0, 0, time.Second*3, func(t *timer.Timer) {
		if err := client.Publish("exchange_spider_man", "binding_key_url", []byte("Hello from client")); err != nil {
			logger.Error("error send msg:", err)
		} else {
			logger.Info("send msg success")
		}
	})

	c := make(chan int)
	<-c
}
