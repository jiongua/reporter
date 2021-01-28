package reporter
//
//import (
//	"context"
//	"github.com/stretchr/testify/mock"
//	"testing"
//)
//
//type MockKafkaProducer struct {
//	mock.Mock
//}
//
//func (m *MockKafkaProducer) Publish(ctx context.Context, key, value []byte) error {
//	m.Called(ctx, key, value)
//	return nil
//}
//
//func (m *MockKafkaProducer) Close() error {
//	m.Called()
//	return nil
//}
//
//func TestDataCollectionServer_CollectAndPublish(t *testing.T) {
//	m := new(MockKafkaProducer)
//	m.On("Publish", context.Background(), nil, []byte("hello")).Return(nil)
//	s := DataCollectionServer{
//		Producer:    m,
//		dataChan:    make(chan DataSource),
//	}
//	s.PublishMessage(context.Background())
//	m.AssertExpectations(t)
//}