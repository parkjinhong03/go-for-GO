package broker

import (
	"fmt"
	"github.com/InVisionApp/go-health/v2"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2/server"
	"log"
)

func TTLCheckHandler(s server.Server, cs *api.Client) func(s *health.State) {
	var isRejected = false
	var cid = fmt.Sprintf("service:%s-%s", s.Options().Name, s.Options().Id)

	return func(s *health.State) {
		if s.Status == "ok" && isRejected {
			log.Printf("[%s] %s recovered", s.Name, cid)
			err := cs.Agent().PassTTL(cid, "broker server recovered.")
			if err != nil { log.Printf("[%s] consul agent error (err: %v)", s.Name, err); return }
			isRejected = false
		}
		if s.Status == "failed" && !isRejected {
			log.Printf("[%s] %s rejected (reason: %s)", s.Name, cid, s.Err)
			err := cs.Agent().FailTTL(cid, "broker server downed.")
			if err != nil { log.Printf("[%s] consul agent error (err: %v)", s.Name, err); return }
			isRejected = true
		}
	}
}