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
	config.Version = sarama.V1_1_1_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Flush.Frequency = time.Millisecond * 10
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.MaxMessages = 10240
	config.Producer.RequiredAcks = sarama.RequiredAcks(-1)
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "admin"
	config.Net.SASL.Password = "admin-secret"
	config.Net.SASL.Handshake = true
	config.Net.SASL.Mechanism = "PLAIN"
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
	err = c.Send("mirrors_npm_publish", r)
	if err != nil {
		fmt.Println("发送失败", err)
		return
	}
	fmt.Println("发送成功...")
}

func main() {
	cfg := &KafkaCfg{
		Brokers:         []string{"9.136.103.27:9092"},
		Topic:           "mirrors_npm_publish",
		ConsumerGroupId: "mirrors_npm_publish_consumer",
		MaxConcurrency:  20,
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
		cg := sarama.NewConfig()
		cg.Producer.Return.Successes = true
		cg.Version = sarama.V1_1_1_0
		//producer, err := mq.NewKafkaClient(&mq.KafkaClientConfig{
		//	Addrs:        []string{"9.136.103.27:9092"},
		//	GroupID:      "mirrors_npm_publish_consumer",
		//	SaramaConfig: cg,
		//})
		//
		//if err != nil {
		//	fmt.Println("生产者初始化失败:", err)
		//	return
		//}
		//for i := 0; i < 5; i++ {
		//	// 从本地发送一条消息，测试是否能正常从kafka中解析
		//	c := KafkaClient{producer}
		//	c.sendMsg()
		//}
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
