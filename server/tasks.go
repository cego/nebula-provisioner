package server

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula"
	"github.com/slyngdk/nebula-provisioner/server/store"
)

type tasks struct {
	config *nebula.Config
	store  *store.Store

	l    *logrus.Logger
	quit chan interface{}
}

func NewTasks(l *logrus.Logger, config *nebula.Config, store *store.Store) *tasks {
	return &tasks{l: l, config: config, store: store, quit: make(chan interface{})}
}

func (t *tasks) Start() {
	t.l.Infoln("Starting task scheduler")

	renewCertTicker := time.NewTicker(t.config.GetDuration("tasks.certRenew.interval", 1*time.Hour))
	renewCATicker := time.NewTicker(t.config.GetDuration("tasks.caRenew.interval", 24*time.Hour))
	dbGCTicker := time.NewTicker(t.config.GetDuration("tasks.dbGC.interval", 5*time.Minute))

	go func() {
		for {
			select {
			case <-renewCertTicker.C:
				t.renewCerts()
			case <-renewCATicker.C:
				t.renewCAs()
			case <-dbGCTicker.C:
				t.dbGC()
			case <-t.quit:
				renewCertTicker.Stop()
				return
			}
		}
	}()
}

func (t *tasks) Stop() {
	t.l.Infoln("Stopping task scheduler")
	t.quit <- struct{}{}
}

func (t *tasks) renewCerts() {
	t.l.Infoln("Task: renew agent certificates")
	err := t.store.RenewCertForAgents()
	if err != nil {
		t.l.WithError(err).Errorln("error when renewing certificates for agents")
	}
}

func (t *tasks) renewCAs() {
	t.l.Infoln("Task: renew ca certificates")
	err := t.store.RenewCAs()
	if err != nil {
		t.l.WithError(err).Errorln("error when renewing ca certificates")
	}
}

func (t *tasks) dbGC() {
	t.l.Debugln("Task: db garbage collection")
	t.store.GC()
}
