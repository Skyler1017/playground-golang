package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"os"
	"sync"
	"time"
)

type Producer struct {
	broker  string
	topic   string
	p       sarama.AsyncProducer
	signals chan os.Signal
}

func NewProducer(broker, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true //指定返回
	p, err := sarama.NewAsyncProducer([]string{broker}, config)
	if err != nil {
		return nil, err
	}
	signals := make(chan os.Signal, 1)
	return &Producer{
		broker:  broker,
		topic:   topic,
		p:       p,
		signals: signals,
	}, nil
}

func (p *Producer) AsyncSend(msg string) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(p sarama.AsyncProducer) {
		defer wg.Done()
		time.Sleep(time.Second * 5)
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					fmt.Println("错误:", err)
					return
				} else {
					return
				}
			case suc := <-success:
				fmt.Println("成功:", suc)
				return
			}
		}
	}(p.p)

	content := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(msg),
	}
	time.Sleep(time.Second * 3)
	p.p.Input() <- content
	wg.Wait()
}

func main() {
	p, err := NewProducer("127.0.0.1:9092", "kafka-test")
	if err != nil {
		fmt.Println("生产者初始化失败")
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		time.Sleep(time.Millisecond * 50)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p.AsyncSend(fmt.Sprintf("发送了第%d条消息", i+1))
		}(i)
	}
	wg.Wait()
}
