package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

//# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kfksender kfksender.go
func main() {
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes = true

	// 使用给定代理地址和配置创建一个同步生产者
	producer, err := sarama.NewSyncProducer([]string{"kafka-node1.8btc-vpc.com:9092", "kafka-node2.8btc-vpc.com:9092", "kafka-node3.8btc-vpc.com:9092"}, config)
	if err != nil {
		panic(err)
	}

	defer producer.Close()

	topic := "pref-pay-notify"
	value := `{"cost":100,"order_no":"-1000","trade_no":"TC20200824000000054"}`

	//构建发送的消息，
	msg := &sarama.ProducerMessage{
		//Topic: "test",//包含了消息的主题
		Partition: int32(10),                   //
		Key:       sarama.StringEncoder("key"), //
	}

	msg.Topic = topic
	//将字符串转换为字节数组
	msg.Value = sarama.ByteEncoder(value)
	//fmt.Println(value)
	//SendMessage：该方法是生产者生产给定的消息
	//生产成功的时候返回该消息的分区和所在的偏移量
	//生产失败的时候返回error
	partition, offset, err := producer.SendMessage(msg)

	if err != nil {
		fmt.Println("Send message Fail")
	}
	fmt.Printf("Partition = %d, offset=%d\n", partition, offset)
}
