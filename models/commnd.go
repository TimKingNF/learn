package models

import (
	"code.google.com/p/mahonia"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sort"
	"time"
)

// 13322947525
// 441522199611031017

const (
	Host string = "127.0.0.1"
	Port string = "8080"
)

func Info(method string, err ...interface{}) {
	fmt.Println(method, "：", err)
}

func ErrorInfo(method string, err interface{}) error {
	return errors.New(method + "：  " + fmt.Sprintf("%s", err))
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

func CreateDir(name string) error {
	src := name + "/"
	if pathIsExist(name) {
		return nil
	}
	err := os.MkdirAll(src, 0777)
	if err != nil {
		if os.IsPermission(err) {
			Info("不够权限创建文件")
		}
	}
	return err
}

func pathIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func GetRandString(s int) string {
	r := Str2Md5(fmt.Sprintf("%d", GetMathRand(s))) //	生成随机数
	g := Str2Md5(fmt.Sprintf("%d", time.Now().Unix()))
	z := Str2Sha1(fmt.Sprintf("%d", time.Now().Unix()))
	x := Str2Sha1(fmt.Sprintf("%d", GetMathRand(s)))
	str := r[0:8] + "-" + g[8:20] + "-" + z[20:28] + "-" + x[28:40]
	return str
}

func Str2Md5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func getCookie(in []*http.Cookie, name string) (string, error) {
	if len(in) <= 0 {
		return "", errors.New("getCookie: slice is null")
	}
	re, _ := regexp.Compile(name + "=(.*); Path=")
	for _, v := range in {
		maps := re.FindStringSubmatch(v.String())
		if len(maps[0]) > 0 {
			return maps[1], nil
		}
	}
	return "", errors.New("getCookie: cookie is null")
}

func utf82gbk(utfStr string) string {
	enc := mahonia.NewEncoder("gbk")
	if ret, ok := enc.ConvertStringOK(utfStr); ok {
		return ret
	}
	return ""
}

func SaveFile(path, filename, filesuffix string, data []byte) (string, error) {
	if pathIsExist(path) != true {
		//	如果路径不存在则创建新路径
		if err := os.MkdirAll(path, 0777); err != nil {
			if os.IsPermission(err) {
				return "", errors.New("SaveFile: not power")
			}
			return "", err
		}
	}
	filePath := path + filename + "." + filesuffix
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	err = ioutil.WriteFile(filePath, data, 0666) //写入文件(字节数组)
	if err != nil {
		return "", err
	}
	return "/" + filePath, nil
}

func SwitchDate(in string) int {
	switch in {
	case "周一":
		return 1
	case "周二":
		return 2
	case "周三":
		return 3
	case "周四":
		return 4
	case "周五":
		return 5
	case "周六":
		return 6
	case "周日":
		return 7
	}
	return 0
}
