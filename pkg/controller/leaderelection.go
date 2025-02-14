/*
Copyright 2019 The HAProxy Ingress Controller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"os"
	"time"

	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

// LeaderElector ...
type LeaderElector interface {
	IsLeader() bool
	Run()
}

// LeaderSubscriber ...
type LeaderSubscriber interface {
	OnStartedLeading(stop <-chan struct{})
	OnStoppedLeading()
	OnNewLeader(identity string)
}

type leaderelector struct {
	logger *logger
	le     *leaderelection.LeaderElector
}

// NewLeaderElector ...
func NewLeaderElector(id string, logger *logger, cache *cache, subscriber LeaderSubscriber) LeaderElector {
	hostname, _ := os.Hostname()
	namespace, podname, err := cache.GetIngressPodName()
	if err != nil {
		logger.Fatal("error reading ingress controller pod: %v", err)
	}

	lock := &resourcelock.ConfigMapLock{
		Client:        cache.client.CoreV1(),
		ConfigMapMeta: metav1.ObjectMeta{Namespace: namespace, Name: id},
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: podname,
			EventRecorder: record.NewBroadcaster().NewRecorder(scheme.Scheme, api.EventSource{
				Component: "haproxy-ingress-leader-elector",
				Host:      hostname,
			}),
		},
	}
	callbacks := leaderelection.LeaderCallbacks{
		OnStartedLeading: func(stop <-chan struct{}) {
			if subscriber != nil {
				subscriber.OnStartedLeading(stop)
			}
		},
		OnStoppedLeading: func() {
			if subscriber != nil {
				subscriber.OnStoppedLeading()
			}
		},
		OnNewLeader: func(identity string) {
			if subscriber != nil {
				subscriber.OnNewLeader(identity)
			}
		},
	}

	le, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock:          lock,
		LeaseDuration: 30 * time.Second,
		RenewDeadline: 15 * time.Second,
		RetryPeriod:   10 * time.Second,
		Callbacks:     callbacks,
	})
	if err != nil {
		logger.Fatal("error starting leader election: %v", err)
	}
	return &leaderelector{
		logger: logger,
		le:     le,
	}
}

func (l *leaderelector) IsLeader() bool {
	return l.le.IsLeader()
}

func (l *leaderelector) Run() {
	go wait.Forever(func() {
		l.le.Run()
	}, 0)
}
