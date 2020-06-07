package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
	"os"
)

func main() {
	name := os.Getenv("name")
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	t, c, err := cfg.New("hello-world", config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	defer c.Close()

	t.StartSpan("say-hello")
	span := t.StartSpan("say-hello")
	span.SetTag("hello-to", name)

	helloStr := fmt.Sprintf("Hello, %s!", name)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	fmt.Println(helloStr)
	span.LogKV("event", "println")

	span.Finish()
}