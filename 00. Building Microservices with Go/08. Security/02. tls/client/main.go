// 우리가 만든 서버의 루트 인증서는 자체 서명되었기 때문에 클라이언트가 이 인증서를 신뢰하고 사용할 수 없다.
// 따라서 클라이언트 측에서 서버로 부터 받은 인증서의 루트 기관이 신뢰할 수 있는 기관이라고 요청을 보내기 전에 명시를 해야한다.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// x509.NewCertPool 함수를 사용하여 새로운 인증서 풀을 생성할 수 있다.
	rootCA := x509.NewCertPool()

	// ioutil.ReadFile 함수를 이용하여 루트 인증서를 바이트 배열 형식으로 변환한다.
	rootCert, err := ioutil.ReadFile("../key/root_cert.pem")
	if err != nil {
		log.Fatal(err)
	}

	// 그리고 위에서 생성한 인증서 풀의 AppendCertsFromPEM 메서드를 이용하여 새 인증서를 풀에 추가할 수 있다.
	ok := rootCA.AppendCertsFromPEM(rootCert)
	if !ok {
		log.Fatal("failed to parse root certificate")
	}

	// 인증서가 유효한 것으로 확인되려면 중간 인증서와 투르 인증서가 모두 필요하기 때문에 어플리케이션 인증서로 인증서 풀에 추가한다.
	applicationCert, err := ioutil.ReadFile("../key/application_cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	ok = rootCA.AppendCertsFromPEM(applicationCert)
	if !ok {
		log.Fatal("failed to parse application certificate")
	}

	// 아래 세 구문을 통해 tls 설정(인증서 설정)을 클라이언트의 전송 계층에 추가시킬 수 있다.
	tlsConfig := &tls.Config{RootCAs: rootCA}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// 마지막으로 client.Get 함수를 통해서 핸드쉐이크와 세션을 진행하고 https 통신을 마무리 한다.
	resp, err := client.Get("https://localhost:8080")

	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}