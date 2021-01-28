package reporter

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"reporter/producer"
)

type DataSource struct {
	Topic string `json:"topic"`
	Content []byte
}

type DataCollectionServer struct {
	//dependency
	Producer     producer.Producer
	dataChan chan DataSource
}

func NewDataCollectionServer(producer producer.Producer) *DataCollectionServer {
	return &DataCollectionServer{
		Producer: producer,
		dataChan: make(chan DataSource),
	}
}

//CollectAndPublish 获取数据，并提交到kafka producer
//CollectAndPublish 将并发执行, 每个goroutines同时从dataChan读取消息，并调用同一个kafka.Writer发布消息
func (d *DataCollectionServer)PublishMessage(ctx context.Context) {
	defer d.Producer.Close()
	loop:
	for {
		select {
		case data := <-d.dataChan:
			log.Printf("get new data: %v", data)
			d.Producer.Publish(ctx, nil, data.Content)
		case <-ctx.Done():
			log.Println(ctx.Err().Error())
			break loop
		}
	}
}

func (d *DataCollectionServer)DoRequest(c *gin.Context) {
	var b = DataSource{
		Topic:   "test001",
		Content: []byte("message a"),
	}
	//c.ShouldBind(&b)
	d.dataChan <- b
	c.JSON(200, gin.H{"message": "ok"})
}










