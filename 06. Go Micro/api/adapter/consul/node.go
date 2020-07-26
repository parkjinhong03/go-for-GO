package consul

import (
	"fmt"
	topic "gateway/topic/golang"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2/registry"
)

const (
	StatusMustBePassing = "Status==passing"
)

func GetServiceNodes(cs *api.Client) (nds []*registry.Node, err error) {
	hcs, _, err := cs.Health().Checks(topic.AuthService, &api.QueryOptions{Filter: StatusMustBePassing})
	if err != nil { return }

	var as *api.AgentService
	for _, hc := range hcs {
		as, _, err = cs.Agent().Service(hc.ServiceID, nil)
		if err != nil { return }
		var md = map[string]string{"CheckID": hc.CheckID}
		nd := &registry.Node{Id: as.ID, Address: fmt.Sprintf("%s:%d", as.Address, as.Port), Metadata: md}
		nds = append(nds, nd)
	}
	return
}