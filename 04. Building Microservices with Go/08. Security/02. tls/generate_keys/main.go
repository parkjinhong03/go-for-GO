package main

import (
	"flag"
	"tls.generate-key.com/generate"
	"tls.generate-key.com/template"
)

func main() {
	password := flag.String("password", "password", "Flag to use to encrypt private keys")
	flag.Parse()

	rootKey, rootCert := generate.KeyAndCertificate(
		&template.RootTemplate,
		nil, nil,
		"../key/root_key.pem",
		"../key/root_cert.pem",
		*password,
		generate.DurationDecade,
	)

	applicationKey, applicationCert := generate.KeyAndCertificate(
		&template.ApplicationTemplate,
		rootCert, rootKey,
		"../key/application_key.pem",
		"../key/application_cert.pem",
		*password,
		generate.DurationYear,
	)

	_, _ = generate.KeyAndCertificate(
		&template.InstanceTemplate,
		applicationCert, applicationKey,
		"../key/instance_key.pem",
		"../key/instance_cert.pem",
		"",
		generate.DurationMonth,
	)
}