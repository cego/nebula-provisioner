//go:build dev
// +build dev

package webapp

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	spaproxy "github.com/lafriks/go-spaproxy"
	"github.com/sirupsen/logrus"
)

var Dir string

func WebHandler(l *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	proxy, err := spaproxy.NewAngularDevProxy(&spaproxy.AngularDevProxyOptions{
		Dir: Dir,
	})
	if err != nil {
		panic(err)
	}

	err = proxy.Start(context.Background())
	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = <-sigs
		err := proxy.Stop()
		if err != nil {
			l.WithError(err).Error("error when closing dev server")
		}
	}()

	return proxy.HandleFunc
}
