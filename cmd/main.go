package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"reporter"
	"reporter/producer"
	"sync"
)

func main() {
	router := gin.New()

	kafkaProducer := producer.CreateKafkaProducer([]string{"192.168.3.101:32769", "192.168.3.101:32768"}, "user_action_collection")
	log.Println("connect kafka ok")
	s := reporter.NewDataCollectionServer(kafkaProducer)
	ctx, cancel := context.WithCancel(context.Background())
	w := new(sync.WaitGroup)
	for i := 0; i < 4; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			s.PublishMessage(ctx)
		}()
	}
	router.POST("/reporter/user_action", s.DoRequest)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	done := make(chan struct{})

	go func() {
		if err := router.Run(":8082"); err != nil {
			log.Println(err.Error())
			done<- struct {}{}
		}
	}()

	select {
	case <-sig:
		cancel()
		w.Wait()
	case <-done:
		cancel()
		w.Wait()
	}
}
