package template

import (
	"crypto/x509"
	"crypto/x509/pkix"
)

var rootTemplate = x509.Certificate{
	// Subject 필드를 이용하여 인증서 요청에 통합될 정보(DN)를 입력할 수 있다.
	Subject: pkix.Name{
		Country:            []string{"KR"},
		Organization:       []string{"PJH Co"},
		OrganizationalUnit: []string{"Tech"},
		CommonName:         "Root",
	},

	KeyUsage: x509.KeyUsageKeyEncipherment |
		x509.KeyUsageDigitalSignature |
		x509.KeyUsageCertSign |
		x509.KeyUsageCRLSign,

	// 아래의 두 필드는 해당 인증서가 인증된 기관(CA)로 부터 생성된 것인지를 임의로 설정할 수 있는 필드이다.
	BasicConstraintsValid: true,
	IsCA:                  true,
}

// x509.Certificate 객체를 이용하여 application 단계의 키 및 인증서를 발급받기 위해 초기 설정이 되어있는 템플릿 선언
var applicationTemplate = x509.Certificate{
	Subject: pkix.Name{
		Country:            []string{"KR"},
		Organization:       []string{"PJH Co"},
		OrganizationalUnit: []string{"Tech"},
		CommonName:         "Application",
	},

	KeyUsage: x509.KeyUsageKeyEncipherment |
		x509.KeyUsageDigitalSignature |
		x509.KeyUsageCertSign |
		x509.KeyUsageCRLSign,
	BasicConstraintsValid: true,
	IsCA:                  true,
}

// x509.Certificate 객체를 이용하여 instance 단계의 키 인증서를 발급받기 위해 초기 설정이 되어있는 템플릿 선언
var instanceTemplate = x509.Certificate{
	Subject: pkix.Name{
		Country:            []string{"KR"},
		Organization:       []string{"PJH Co"},
		OrganizationalUnit: []string{"Tech"},
		CommonName:         "Instance",
	},

	KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	BasicConstraintsValid: true,
	// DNSNames 필드를 이용하여 해당 인증서를 사용할 수 있는 DNS를 표시할 수 있다.
	DNSNames:              []string{"localhost"},
}