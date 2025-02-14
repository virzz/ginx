package ginx_test

import (
	"testing"
	"time"

	"github.com/virzz/ginx"
	"github.com/virzz/vlog"
)

func TestNew(t *testing.T) {
	httpSrv, err := ginx.New(&ginx.Config{
		Addr:    "127.0.0.1",
		Port:    3003,
		Pprof:   true,
		Metrics: true,
	}, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil {
			vlog.Error("Failed to run http server", "err", err.Error())
		}
	}()

	time.Sleep(10 * time.Second)

}
