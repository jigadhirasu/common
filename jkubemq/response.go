package jkubemq

import (
	"time"

	"github.com/kubemq-io/kubemq-go"
)

func Response(receive interface{}) *kubemq.Response {
	switch r := receive.(type) {
	case *kubemq.CommandReceive:
		return kubeMQ.R().SetRequestId(r.Id).SetResponseTo(r.ResponseTo).SetExecutedAt(time.Now())
	case *kubemq.QueryReceive:
		return kubeMQ.R().SetRequestId(r.Id).SetResponseTo(r.ResponseTo).SetExecutedAt(time.Now())
	}
	return nil
}
