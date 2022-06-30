//go:build dev
// +build dev

package webapp

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"

	"github.com/lafriks/go-spaproxy"
	"github.com/sirupsen/logrus"
)

var Dir string

func WebHandler(l *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	proxy, err := newAngularDevProxy()
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

func newAngularDevProxy() (spaproxy.SpaDevProxy, error) {
	port, err := spaproxy.GetFreePort(0)
	if err != nil {
		return nil, err
	}

	args := make([]string, 0)
	args = append(args, "--port", strconv.Itoa(port))
	args = append(args, "--host", "localhost")
	args = append(args, "--open", "false")

	return spaproxy.NewSpaDevProxy(&spaproxy.SpaDevProxyOptions{
		RunnerType:  spaproxy.RunnerTypeNpm,
		ScriptName:  "start",
		Dir:         Dir,
		Env:         []string{},
		Args:        args,
		Port:        port,
		StartRegexp: regexp.MustCompile("is listening on"),
	})
}
