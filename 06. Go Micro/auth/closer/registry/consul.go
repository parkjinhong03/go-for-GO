package registry

import (
	"auth/tool/addr"
	"fmt"
	"github.com/hashicorp/consul/api"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	"strconv"
	"strings"
)

func ConsulServiceRegister(s server.Server, cs *api.Client) func() error {
	return func() (err error) {
		ps := strings.Split(s.Options().Address, ":")[3]
		port, err := strconv.Atoi(ps)
		ip := addr.GetLocalAddr().IP
		if err != nil { log.Fatal(err) }
		sid := s.Options().Name + "-" + s.Options().Id
		cid := "service:" + sid

		asr := &api.AgentServiceRegistration{
			ID:      sid,
			Name:    s.Options().Name,
			Port:    port,
			Address: ip.String(),
		}
		err = cs.Agent().ServiceRegister(asr)
		if err != nil { log.Fatal(err) }

		asc := api.AgentServiceCheck{
			Name:   s.Options().Name,
			Status: "passing",
			TTL:    "8640s",
		}
		acr := &api.AgentCheckRegistration{
			ID:                cid,
			Name:              fmt.Sprintf("service '%s' check", s.Options().Name),
			ServiceID:         sid,
			AgentServiceCheck: asc,
		}
		err = cs.Agent().CheckRegister(acr)
		if err != nil { log.Fatal(err) }

		log.Infof("succeed to registry service and check to consul!! (service id: %s | check id: %s)\n", sid, cid)
		return
	}
}