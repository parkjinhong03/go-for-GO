package parser

import (
	"errors"
	"strconv"
	"strings"
)

var InvalidError = errors.New("that error is not invalid gorm error")

func DBErrorParse(errStr string) (code int, err error) {
	strArr := strings.Split(errStr, " ")
	if strArr[0] != "Error" { err = InvalidError; return }
	if code, err = strconv.Atoi(strArr[1][:4]); err != nil { err = InvalidError }
	return
}