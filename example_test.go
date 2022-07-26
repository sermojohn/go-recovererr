package recovererr

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	mhwr := retrierMessageHandler{
		messageHandler: &dummyMessageHandler{},
		retrier:        NewRetrier(WithRetryPolicy(RetryRecoverablePolicy), WithIntervalGenerator(StaticIntervalGenerator(10*time.Millisecond))),
		timeout:        100 * time.Millisecond,
	}
	err := mhwr.Handle(&Message{})

	fmt.Printf("handler completed with output: %v\n", err)
}

type Message struct {
	Topic     string    `json:"topic"`
	Key       []byte    `json:"key"`
	Value     []byte    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type MessageHandler interface {
	Handle(*Message) error
}

type dummyMessageHandler struct{}

func (dmh *dummyMessageHandler) Handle(_ *Message) error {
	fmt.Println("called handler")
	time.Sleep(10 * time.Millisecond)
	return Recoverable(fmt.Errorf("handling timed out"))
}

type retrierMessageHandler struct {
	retrier        *Retrier
	timeout        time.Duration
	messageHandler MessageHandler
}

func (rmh *retrierMessageHandler) Handle(m *Message) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), rmh.timeout)
	defer cancelFunc()

	retrier := NewRetrier(
		WithRetryPolicy(RetryRecoverablePolicy),
		WithIntervalGenerator(StaticIntervalGenerator(rmh.timeout)),
	)

	return retrier.Do(ctx, func() error { return rmh.messageHandler.Handle(m) })
}
