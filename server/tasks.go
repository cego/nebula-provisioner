package server

import (
	"time"

	"github.com/cego/nebula-provisioner/server/store"
	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula/config"
)

type tasks struct {
	config *config.C
	store  *store.Store

	l    *logrus.Logger
	quit chan interface{}
}

func NewTasks(l *logrus.Logger, config *config.C, store *store.Store) *tasks {
	return &tasks{l: l, config: config, store: store, quit: make(chan interface{})}
}

func (t *tasks) Start() {
	t.l.Infoln("Starting task scheduler")

	renewCertDuration := func() time.Duration { return t.config.GetDuration("tasks.certRenew.interval", 1*time.Hour) }
	renewCADuration := func() time.Duration { return t.config.GetDuration("tasks.caRenew.interval", 24*time.Hour) }
	dbGCDuration := func() time.Duration { return t.config.GetDuration("tasks.dbGC.interval", 5*time.Minute) }

	renewCertTicker := time.NewTicker(renewCertDuration())
	renewCATicker := time.NewTicker(renewCADuration())
	dbGCTicker := time.NewTicker(dbGCDuration())

	t.config.RegisterReloadCallback(func(_ *config.C) {
		t.l.Info("Reloading task scheduler")
		renewCertTicker.Reset(renewCertDuration())
		renewCATicker.Reset(renewCADuration())
		dbGCTicker.Reset(dbGCDuration())
	})

	go func() {
		for {
			select {
			case <-renewCertTicker.C:
				t.renewCerts()
			case <-renewCATicker.C:
				t.renewCAs()
				t.expireCAs()
			case <-dbGCTicker.C:
				t.dbGC()
			case <-t.quit:
				renewCertTicker.Stop()
				renewCATicker.Stop()
				dbGCTicker.Stop()
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

func (t *tasks) expireCAs() {
	t.l.Infoln("Task: update status for expired ca certificates")
	err := t.store.UpdateExpiredCAs()
	if err != nil {
		t.l.WithError(err).Errorln("error when updating expire for ca certificates")
	}
}

func (t *tasks) dbGC() {
	t.l.Debugln("Task: db garbage collection")
	t.store.GC()
}
