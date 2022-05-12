package jkubemq

import (
	"context"
	"os"
	"strconv"

	"github.com/jigadhirasu/common/j"
	"github.com/kubemq-io/kubemq-go"
)

type KubeMQFunc func(tags j.Tags, bytes j.Bytes) j.Bytes

var kubeMQ *kubemq.Client

// use context or gen new context
func Client(ctxx ...context.Context) *kubemq.Client {
	if kubeMQ != nil {
		return kubeMQ
	}

	address := os.Getenv("KUBEMQ_ADDR")
	port, _ := strconv.Atoi(os.Getenv("KUBEMQ_PORT"))
	hostname := os.Getenv("HOSTNAME")

	if len(ctxx) < 1 {
		ctxx = append(ctxx, context.Background())
	}
	ctx := ctxx[0]

	c, err := kubemq.NewClient(ctx,
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAddress(address, port),
		kubemq.WithClientId(hostname),
		kubemq.WithAutoReconnect(true),
		kubemq.WithMaxReconnects(16),
	)

	if err != nil {
		panic(err)
	}

	kubeMQ = c
	return kubeMQ
}
