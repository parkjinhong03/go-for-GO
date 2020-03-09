// load balancing(부하 분산)이란 하나의 인터넷 서비스가 발생하는 트래픽이 많을 때 여러 대의 서버가 분산하여 처리하는 것 이다.
// RR, WRR, LC, WLC와 같은 스케줄링 알고리즘을 이용하여 서버의 로드율 증가, 부하량, 속도저하 등을 막을 수 있다.

package main

import (
	"bytes"
	"fmt"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"sync/atomic"
)

type Strategy interface {
	NextService() (*Service, error)
	SetServices([]*Service)
}

type Service struct {
	Host	 	string
	Port 		int32
	Prefix 		string
	MaxConnect	int32
	CurConnect 	int32
}

func (e *Service) HTTPRequest(r *http.Request) (*http.Response, error) {
	body, _ := ioutil.ReadAll(r.Body)

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, fmt.Sprintf("http://%s:%d%s%s", e.Host, e.Port, e.Prefix, r.URL.Path), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "LC Load Balancer")

	atomic.AddInt32(&e.CurConnect, 1)
	resp, err := client.Do(req)
	atomic.AddInt32(&e.CurConnect, -1)

	if err != nil {
		return nil, err
	}
	return resp, err
}

func NewService(host string, port int32, prefix string, connections int32) *Service {
	return &Service{
		Host:        host,
		Port:        port,
		Prefix:      prefix,
		MaxConnect:  connections,
		CurConnect:  0,
	}
}

type LCStrategy struct {
	Services []*Service
}

func (l *LCStrategy) NextService() (*Service, error) {
	var min int32 = -1
	var service *Service

	for _, s := range l.Services {
		if s.CurConnect == s.MaxConnect {
			continue
		} else if min == -1 {
			min = s.CurConnect
			service = s
		} else if s.CurConnect < min {
			min = s.CurConnect
			service = s
		}
	}

	if min == -1 {
		return &Service{}, fmt.Errorf("no more APIs are available for connection")
	}
	return service, nil
}

func (l *LCStrategy) SetServices(services []*Service) {
	l.Services = services
}

type Loadbalancer struct {
	strategy Strategy
}

func NewLoadBalancer(strategy Strategy, services []*Service) Loadbalancer {
	strategy.SetServices(services)
	return Loadbalancer{
		strategy: strategy,
	}
}

func (l *Loadbalancer) GetService() (*Service, error) {
	return l.strategy.NextService()
}

func (l *Loadbalancer) UpdateServices(services []*Service) {
	l.strategy.SetServices(services)
}

type LoadBalancingHandler struct {
	lb Loadbalancer
}

func (lbh *LoadBalancingHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	service, err := lbh.lb.GetService()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusTooManyRequests)
		return
	}
	fmt.Println(service.Port, service.CurConnect)

	resp, err := service.HTTPRequest(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	_, _ = rw.Write(body)
}

func main() {
	endpoints := []*Service{
		NewService("localhost", 801, "/search", 5),
		NewService("localhost", 802, "/search", 5),
	}

	handler := LoadBalancingHandler{
		lb:   NewLoadBalancer(&LCStrategy{}, endpoints),
	}

	http.Handle("/search/", http.StripPrefix("/search", &handler))
	log.Fatal(http.ListenAndServe(":80", nil))
}