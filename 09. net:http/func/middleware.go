package net

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type Middleware func(next HandlerFunc) HandlerFunc

func LogHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		// next(c)가 실행하기 전에 현재 시간을 기록
		t := time.Now()

		// 다음 핸들러 수행
		next(c)

		// 웹 요청 정보와 전체 소요 시간을 로그로 남김
		log.Printf("[%s] %q %v\n",
			c.Request.Method,
			c.Request.URL.String(),
			time.Now().Sub(t))
	}
}

func RecoverHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(c.ResponseWriter,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()
		next(c)
	}
}

func ParseFormHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		c.Request.ParseForm()
		fmt.Println(c.Request.PostForm)
		for k, v := range c.Request.PostForm {
			if len(k) > 0 {
				c.Params[k] = v[0]
			}
		}
		next(c)
	}
}

func ParseJsonBodyHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		var m map[string]interface{}
		if json.NewDecoder(c.Request.Body).Decode(&m); len(m)>0 {
			for k, v := range m {
				c.Params[k] = v
			}
		}
		next(c)
	}
}

func staticHandler(next HandlerFunc) HandlerFunc {
	var (
		dir = http.Dir(".")
		indexFile = "index.html"
	)

	return func(c *Context) {
		// http 메서드가 GET이나 HEAD가 아니면 바로 다음 핸들러 수행
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			next(c)
			return
		}

		file := c.Request.URL.Path
		// URL 경로에 해당하는 파일 열기 시도

		f, err := dir.Open(file)
		if err != nil {
			// URL 경로에 해당하는 파일 열기에 실패하면 바로 다음 핸들러 수행
			next(c)
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			// 파일의 상태가 정상이 아니라면 바로 다음 핸들러 수행
			next(c)
			return
		}

		// URL 경로가 디렉터리라면 indexFile을 사용
		if fi.IsDir() {
			if !strings.HasSuffix(c.Request.URL.Path, "/") {
				http.Redirect(c.ResponseWriter, c.Request, c.Request.URL.Path+"/", http.StatusFound)
				return
			}

			// 디렉터리를 가르키는 URL 경로에 indexFile 이름을 붙여서 전체 파일 경로 생성
			file = path.Join(file, indexFile)

			// indexFile 열기 시도
			f, err = dir.Open(file)
			if err != nil {
				next(c)
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				next(c)
				return
			}
		}

		// file의 내용 전달(next 핸들러로 제어권을 넘기지 않고 요청 처리를 종료함)
		http.ServeContent(c.ResponseWriter, c.Request, file, fi.ModTime(), f)
	}
}