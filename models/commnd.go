package models

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)

const (
	Host = "127.0.0.1"
	Port = "8080"
)

func Info(method string, err ...interface{}) {
	fmt.Println(method, "：", err)
}

func ErrorInfo(method string, err interface{}) error {
	return errors.New(fmt.Sprintln(method, "：", err))
}

func GetSignature(token, timestamp, sessid string) string {
	tmpArr := []string{token, timestamp, sessid}
	sort.Strings(tmpArr)
	tmpStr := tmpArr[0] + tmpArr[1] + tmpArr[2]
	return Str2Sha1(tmpStr)
}

func Str2Sha1(data string) string {
	t := sha1.New()         //	创建sha1字符串
	io.WriteString(t, data) //	调用io包写入数据
	//	返回加密字符串
	return fmt.Sprintf("%x", t.Sum(nil))
}

func GetMathRand(scope int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(scope)
}
