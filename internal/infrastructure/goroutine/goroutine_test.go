package goroutine

import (
	"context"
	"testing"
	"time"
)

func Test_SafeGo(t *testing.T) {
	ch := make(chan struct{}, 1)

	f := func(c context.Context, _ struct{}) {
		defer func() {
			ch <- struct{}{}
		}()
		panic("test")
	}

	SafeGo(context.TODO(), f, struct{}{})
	<-ch
	time.Sleep(time.Millisecond * 10)
}
