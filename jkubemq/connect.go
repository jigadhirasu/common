package jkubemq

import (
	"context"
	"os"
	"strconv"

	"github.com/jigadhirasu/jgorm/jcommon"
	"github.com/kubemq-io/kubemq-go"
)

type KubeMQFunc func(tags jcommon.Tags, bytes jcommon.Bytes) jcommon.Bytes

var kubeMQ *kubemq.Client

func KubeMQ(ctxx ...context.Context) *kubemq.Client {
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
