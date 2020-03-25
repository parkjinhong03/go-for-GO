package main

import (
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
)

func main() {
	// config.Load와 file.NewSource 함수를 이용하여 파일 소스로부터 config를 불러올 수 있다.
	if err := config.Load(file.NewSource(
		// 가져올 파일의 위치를 file.WithPath 함수를 이용하여 등록할 수 있다.
		file.WithPath("./config.json"),
	)); err != nil {
		fmt.Println(err)
		return
	}

	// json 데이터에 접근하기 위해 데이터를 담을 객체를 정의한다.
	type Host struct {
		Address string `json:"address"`
		Port    int `json:"port"`
	}

	var host Host

	// config.Get 함수를 이용해 읽어올 json 데이터의 세부적인 위치를 정하고, Scan 메서드를 이용하여 내용을 구조체에 담을 수 있다.
	if err := config.Get("hosts", "database").Scan(&host); err != nil {
		fmt.Println(err)
		return
	}

	// database | host: 10.0.0.1, port: 3306
	fmt.Printf("database | host: %s, port: %d\n", host.Address, host.Port)
}