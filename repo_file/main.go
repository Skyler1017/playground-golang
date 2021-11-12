package main

import (
	"encoding/json"
	"fmt"
	"git.code.oa.com/omc-dev/mq"
	"github.com/Shopify/sarama"
	"sync"
	"time"
)

type KafkaClient struct {
	*mq.KafkaClient
}

type Event struct {
	Platform  string `json:"platform"`
	RepoName  string `json:"repo_name"`
	Path      string `json:"path"`
	User      string `json:"user"`
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

func printEvent(event *Event) {
	timeFormatLayout := "01-02 15:04:05"
	formatTimeStr := time.Unix(event.Timestamp, 0).Format(timeFormatLayout)
	fmt.Printf("%s【%s】(%s)%s \n", formatTimeStr, event.Platform, event.User, event.Path)
}

func (c *KafkaClient) sendMsg() {
	payload, err := json.Marshal(&Event{
		RepoName:  "test-pkg",
		User:      "hui",
		Path:      "/your/path",
		Platform:  "maven",
		Timestamp: time.Now().Unix(),
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
	err = c.Send("repo_file", r)
	if err != nil {
		fmt.Println("发送失败", err)
		return
	}
	fmt.Println("发送成功...")
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
	pe := &Event{}
	err = json.Unmarshal(payload, pe)
	if err != nil {
		fmt.Printf("无法解析json结构\n")
		return err
	}
	printEvent(pe)
	return
}

func main() {
	kc, err := mq.NewKafkaClient(&mq.KafkaClientConfig{
		Addrs:        []string{"91.136.37.191:9092"},
		GroupID:      "mirrors_npm_publish_consumer",
		SaramaConfig: getKafkaConfig(),
	})
	if err != nil {
		panic(err)
	}
	client := &KafkaClient{kc}
	err = client.Subscribe(client.handleMsg, "repo_file")
	if err != nil {
		panic(err)
	}
	go func() {
		cg := sarama.NewConfig()
		cg.Producer.Return.Successes = true
		cg.Version = sarama.V1_1_1_0
		producer, err := mq.NewKafkaClient(&mq.KafkaClientConfig{
			Addrs:        []string{"91.136.37.191:9092"},
			GroupID:      "mirrors_npm_publish_consumer",
			SaramaConfig: cg,
		})

		if err != nil {
			fmt.Println("生产者初始化失败:", err)
			return
		}
		for i := 0; i < 5; i++ {
			// 从本地发送一条消息，测试是否能正常从kafka中解析
			c := KafkaClient{producer}
			c.sendMsg()
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

// 获取kafka配置
func getKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_1_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Flush.Frequency = time.Millisecond * 10
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.MaxMessages = 1024
	config.Producer.RequiredAcks = sarama.RequiredAcks(-1)
	//config.Net.SASL.Enable = true
	//config.Net.SASL.User = "ckafka-ydrv28xg#mirrors_npm"
	//config.Net.SASL.Password = "mirrors_npm_prod"
	return config
}
