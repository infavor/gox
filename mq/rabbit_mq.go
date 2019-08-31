package mq

import (
	"errors"
	"fmt"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/uuid"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type Binding struct {
	QueueName    string
	ExchangeName string
	BindingKey   string
}

type RabbitClient struct {
	ConnectionString   string
	connection         *amqp.Connection
	channel            *amqp.Channel
	stopWorld          *sync.Mutex
	AutoACK            bool
	QueueNames         []string
	ConsumeQueues      []string
	ExchangeNames      []string
	DefaultContentType string
	Bindings           []Binding
	PrefetchCount      int
	ConsumeHandler     func(queueName string, d amqp.Delivery) error
}

func (c *RabbitClient) Init() {
	c.stopWorld = new(sync.Mutex)
	if c.DefaultContentType == "" {
		c.DefaultContentType = "text/plain"
	}
	if c.PrefetchCount == 0 {
		c.PrefetchCount = 1
	}
	c.reconnect()
}

func (c *RabbitClient) Publish(exchange, bindingKey string, data []byte) error {
	if c.channel == nil {
		return errors.New("channel is not ready")
	}
	if err := c.channel.Publish(
		exchange,   // exchange
		bindingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  c.DefaultContentType,
			Body:         data,
		},
	); err != nil {
		go c.reconnect()
		return err
	}
	return nil
}

func (c *RabbitClient) Consume() error {
	consumes := make(map[string]<-chan amqp.Delivery)
	for _, name := range c.ConsumeQueues {
		d, err := c.channel.Consume(
			name,      // queue
			"",        // consumer
			c.AutoACK, // auto-ack
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)
		if err != nil {
			logger.Error("error consume queue: ", err)
			go c.reconnect()
			return err
		} else {
			logger.Info("declare consumer success")
			consumes[name] = d
		}
	}

	for name, msgs := range consumes {
		_name, _msgs := name, msgs
		go func() {
			for d := range _msgs {
				dd := d
				if err := c.ConsumeHandler(_name, dd); err != nil {
					logger.Error("error consume msg: ", err)
				}
			}
			logger.Error("msgs break")
		}()
	}
	return nil
}

func (c *RabbitClient) reconnect() {
	c.stopWorld.Lock()
	defer c.stopWorld.Unlock()
	c.Close()
	c.channel = nil
	c.connection = nil
	c.connection = c.getConnection()
}

func (c *RabbitClient) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.connection != nil {
		c.connection.Close()
	}
}

func (c *RabbitClient) getConnection() *amqp.Connection {
	for {
		conn, err := amqp.Dial(c.ConnectionString)
		if err == nil {
			logger.Warn("now connection is ", &conn)
			closeConn := make(chan *amqp.Error)
			go func() {
				way := uuid.UUID()
				fmt.Println("way:", way)
				for {
					err := <-closeConn
					if err != nil {
						logger.Error("rabbit connection is closed: ", err)
						go c.reconnect()
						break
					} else {
						logger.Warn("??????::", way)
						time.Sleep(time.Second)
					}
				}
			}()
			conn.NotifyClose(closeConn)
			c.connection = conn

			if err = c.initChannel(); err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
			if err = c.declareQueues(); err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
			if err = c.declareExchanges(); err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
			if err = c.bindQueue(); err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
			if err = c.Consume(); err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
			return conn
		}
		logger.Error("failed connect to rabbitmq server: ", err)
		time.Sleep(5 * time.Second)
		logger.Info("trying to reconnect to RabbitMQ at ", c.ConnectionString)
	}
}

func (c *RabbitClient) initChannel() error {
	// 打开管道
	ch, err := c.connection.Channel()
	if err != nil {
		logger.Error("failed to open channel: ", err)
		return err
	}
	if err = ch.Qos(c.PrefetchCount, 0, false); err != nil {
		logger.Error(err)
		return err
	}
	c.channel = ch
	return nil
}

func (c *RabbitClient) declareQueues() error {
	for _, name := range c.QueueNames {
		_, err := c.channel.QueueDeclare(
			name,  // name
			true,  // durable 是否持久化消息
			false, // delete when unused 未使用时删除
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			logger.Error("failed to declare queue: ", name)
			return err
		}
		logger.Info("bind queue success: ", name)
	}
	return nil
}

func (c *RabbitClient) declareExchanges() error {
	for _, name := range c.ExchangeNames {
		// 声明交换器
		err := c.channel.ExchangeDeclare(
			name,
			"topic",
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			logger.Error("failed to declare exchange: ", name)
			return err
		}
		logger.Info("bind exchange success: ", name)
	}
	return nil
}

func (c *RabbitClient) bindQueue() error {
	for _, bind := range c.Bindings {
		err := c.channel.QueueBind(bind.QueueName, bind.BindingKey, bind.ExchangeName, false, nil)
		if err != nil {
			logger.Error("failed to bind queue '", bind.QueueName, "' with exchange '", bind.ExchangeName, "': ", err)
			return err
		}
		logger.Info("bind queue success: ", bind.BindingKey)
	}
	return nil
}
