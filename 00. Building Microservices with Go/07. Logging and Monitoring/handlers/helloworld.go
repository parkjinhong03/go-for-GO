package handlers

import (
	"building-microservices-with-go.com/logging/entities"
	"building-microservices-with-go.com/logging/httputil"
	"encoding/json"
	"github.com/alexcesaro/statsd"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

type helloWorldHandler struct {
	statsD *statsd.Client
	logger *logrus.Logger
}

func (h *helloWorldHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var status int

	name := r.Context().Value("name").(string)
	response := entities.HelloWorldResponse{Message: "Hello " + name}
	encoder := json.NewEncoder(rw)

	err := encoder.Encode(response)
	if err != nil {
		status = http.StatusInternalServerError
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		status = http.StatusOK
	}

	request := httputil.RequestSerializer{Request: r}
	entry := h.logger.WithFields(logrus.Fields{
		"group": "handler",
		"segment": "helloWorld",
		"outcome": status,
	})

	if err != nil {
		entry.Fatal(request.ToJSON())
	} else {
		entry.Info(request.ToJSON())
	}

	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
}

func NewHelloWorldHandler(logger *logrus.Logger) *helloWorldHandler {
	return &helloWorldHandler{
		logger: logger,
	}
}