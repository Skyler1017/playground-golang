package main

import (
	"encoding/json"
	"fmt"
	"git.code.oa.com/omc-dev/mq"
	"github.com/Shopify/sarama"
	"sync"
	"time"
)

type KafkaCfg struct {
	Brokers         []string `yaml:"brokers"`
	Topic           string   `yaml:"topic"`
	ConsumerGroupId string   `yaml:"consumer_group_id"`
	MaxConcurrency  int      `yaml:"max_concurrency"`
}

type KafkaClient struct {
	*mq.KafkaClient
}

// 获取kafka配置
func getKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_1_0 // 内网自己搭的kafka版本
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Flush.Frequency = time.Millisecond * 10
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.MaxMessages = 10240
	config.Producer.RequiredAcks = sarama.RequiredAcks(-1)
	return config
}

type NpmPublishEvent struct {
	PackageName string `json:"package_name"`
	Publisher   string `json:"publisher"`
	Version     string `json:"version"`
	Timestamp   int64  `json:"timestamp"`
}

func printEvent(event *NpmPublishEvent) {
	timeFormatLayout := "2006-01-02 15:04:05"
	formatTimeStr := time.Unix(int64(event.Timestamp), 0).Format(timeFormatLayout)
	fmt.Printf("%s在%s发布了%s包%s\n", event.Publisher, formatTimeStr, event.PackageName, event.Version)
}

func (c *KafkaClient) handleMsg(msg mq.Message) (err error) {
	needACK := true
	payload := msg.Payload()
	defer func() {
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		if needACK {
			if err = c.ACK(msg); err != nil {
				fmt.Printf("%+v\n", err)
			}
		}
	}()
	pe := &NpmPublishEvent{}
	err = json.Unmarshal(payload, pe)
	if err != nil {
		fmt.Printf("无法解析json结构\n")
		return err
	}
	printEvent(pe)
	return
}

func (c *KafkaClient) sendMsg() {
	payload, err := json.Marshal(&NpmPublishEvent{
		PackageName: "test-pkg",
		Publisher:   "hui",
		Timestamp:   time.Now().Unix(),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	r := &mq.KafkaMessage{
		Message: &mq.GenericMessage{
			Buf: payload,
		},
		Key: fmt.Sprintf("%s:%s", "test-pkg", time.Now()),
	}
	fmt.Println("发送消息...")
	err = c.Send("mirrors_npm_publish", r)
	if err != nil {
		fmt.Println("发送失败", err)
	}
}

func main() {
	cfg := &KafkaCfg{
		Brokers:         []string{"127.0.0.1:9092"},
		Topic:           "mirrors_npm_publish",
		ConsumerGroupId: "mirrors_npm_publish_consumer",
		MaxConcurrency:  10,
	}
	kc, err := mq.NewKafkaClient(&mq.KafkaClientConfig{
		Addrs:        cfg.Brokers,
		GroupID:      cfg.ConsumerGroupId,
		SaramaConfig: getKafkaConfig(),
	})
	if err != nil {
		panic(err)
	}
	client := &KafkaClient{kc}
	err = client.Subscribe(client.handleMsg, cfg.Topic)
	if err != nil {
		panic(err)
	}
	go func() {
		for i := 0; i < 1; i++ {
			client.sendMsg()
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
