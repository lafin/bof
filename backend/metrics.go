package main

import (
	"github.com/go-pkgz/lgr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var pusher *push.Pusher
var l = lgr.New(lgr.Msec, lgr.Debug, lgr.CallerFile, lgr.CallerFunc)

var (
	taskErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "lafin_bof_errors",
	}, []string{"error"})
)

func metrics(url, job string) {
	if url != "" && job != "" {
		registry := prometheus.NewRegistry()
		registry.MustRegister(taskErrors)
		pusher = push.New(url, job).Gatherer(registry)
	}
}

func pushMetrics() {
	if pusher != nil {
		if err := pusher.Push(); err != nil {
			l.Logf("ERROR could not push to Pushgateway, %v", err)
		}
		taskErrors.Reset()
	}
}
