package reporter

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reporter/producer"
)

type DataSource struct {
	Topic string `form:"topic" binding:"required"`
	Content []byte 	`form:"content" binding:"required"`
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
			log.Println("get new data...")
			d.Producer.Publish(ctx, nil, data.Content)
		case <-ctx.Done():
			log.Println(ctx.Err().Error())
			break loop
		}
	}
}

func (d *DataCollectionServer)DoRequest(c *gin.Context) {
	var b DataSource
	err := c.ShouldBind(&b)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if b.Topic != "user_action_collection" {
		log.Printf("invalid topic: %s\n", b.Topic)
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("bad topic <%s>, want <user_action_collection>", b.Topic)})
		return
	}
	//may send to kafka directly
	d.dataChan <- b
	c.JSON(200, gin.H{"message": "ok"})
}










