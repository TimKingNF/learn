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
	"strconv"
	"time"
)

const (
	Host    string = "127.0.0.1"
	Port    string = "8080"
	Package string = "g:/mygo/src/learn"
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
	if PathIsExist("/" + name) {
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

func PathIsExist(path string) bool {
	_, err := os.Stat(Package + path)
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

func Utf82gbk(utfStr string) string {
	enc := mahonia.NewEncoder("gbk")
	if ret, ok := enc.ConvertStringOK(utfStr); ok {
		return ret
	}
	return ""
}

func SaveFile(path, filename, filesuffix string, data []byte) (string, error) {
	if PathIsExist(path) != true {
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

func FilterRepeat(m interface{}) interface{} {
	switch m.(type) {
	case []*Term:
		for key1, val1 := range m.([]*Term) {
			for key2 := key1 + 1; key2 < len(m.([]*Term)); key2++ {
				if val1.Id == m.([]*Term)[key2].Id {
					m = Remove(m.([]*Term), key2, key2+1)
					key2--
				}
			}
		}
		return m
	case []string:
		for key1, val1 := range m.([]string) {
			for key2 := key1 + 1; key2 < len(m.([]string)); key2++ {
				if val1 == m.([]string)[key2] {
					m = Remove(m.([]string), key2, key2+1)
					key2--
				}
			}
		}
		return m
	case []int:
		for key1, val1 := range m.([]int) {
			for key2 := key1 + 1; key2 < len(m.([]int)); key2++ {
				if val1 == m.([]int)[key2] {
					m = Remove(m.([]int), key2, key2+1)
					key2--
				}
			}
		}
		return m
	case []*StudentCourse:
		for key1, val1 := range m.([]*StudentCourse) {
			for key2 := key1 + 1; key2 < len(m.([]*StudentCourse)); key2++ {
				if val1.TeacherCourse.Course.Name == m.([]*StudentCourse)[key2].TeacherCourse.Course.Name {
					m = Remove(m.([]*StudentCourse), key2, key2+1)
					key2--
				}
			}
		}
		return m
	case []map[string]interface{}:
		for key1, val1 := range m.([]map[string]interface{}) {
			for key2 := key1 + 1; key2 < len(m.([]map[string]interface{})); key2++ {
				if val1["type"] == m.([]map[string]interface{})[key2]["type"] {
					switch val1["type"] {
					case "teacher":
						if val1["sender"].(*Teacher).Id == m.([]map[string]interface{})[key2]["sender"].(*Teacher).Id {
							m = Remove(m.([]map[string]interface{}), key2, key2+1)
							key2--
						}
					case "admin":
						if val1["sender"].(*Admin).Id == m.([]map[string]interface{})[key2]["sender"].(*Admin).Id {
							m = Remove(m.([]map[string]interface{}), key2, key2+1)
							key2--
						}
					}
				}
			}
		}
		return m
	}
	return nil
}

//	删除数据
func Remove(slice interface{}, start, end int) interface{} {
	switch slice.(type) {
	case []*Term:
		return append(slice.([]*Term)[:start], slice.([]*Term)[end:]...)
	case []string:
		return append(slice.([]string)[:start], slice.([]string)[end:]...)
	case []int:
		return append(slice.([]int)[:start], slice.([]int)[end:]...)
	case []*StudentCourse:
		return append(slice.([]*StudentCourse)[:start], slice.([]*StudentCourse)[end:]...)
	case []map[string]interface{}:
		return append(slice.([]map[string]interface{})[:start], slice.([]map[string]interface{})[end:]...)
	}
	return nil
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

func Sort(m interface{}) SortMap {
	if m == nil {
		return nil
	}
	in := NewSortMap(m)
	sort.Sort(in)
	return in
}

type SortMap []SortMapItem

type SortMapItem struct {
	Key int
	Val interface{}
}

func NewSortMap(m interface{}) SortMap {
	switch m.(type) {
	case map[int]map[int]map[string]*StudentCourse:
		ms := make(SortMap, 0, len(m.(map[int]map[int]map[string]*StudentCourse)))
		for k, v := range m.(map[int]map[int]map[string]*StudentCourse) {
			ms = append(ms, SortMapItem{k, v})
		}
		return ms
	case map[int]map[string]*StudentCourse:
		ms := make(SortMap, 0, len(m.(map[int]map[string]*StudentCourse)))
		for k, v := range m.(map[int]map[string]*StudentCourse) {
			ms = append(ms, SortMapItem{k, v})
		}
		return ms
	case map[int]map[int]map[string]*TeacherCourse:
		ms := make(SortMap, 0, len(m.(map[int]map[int]map[string]*TeacherCourse)))
		for k, v := range m.(map[int]map[int]map[string]*TeacherCourse) {
			ms = append(ms, SortMapItem{k, v})
		}
		return ms
	case map[int]map[string]*TeacherCourse:
		ms := make(SortMap, 0, len(m.(map[int]map[string]*TeacherCourse)))
		for k, v := range m.(map[int]map[string]*TeacherCourse) {
			ms = append(ms, SortMapItem{k, v})
		}
		return ms
	}
	return nil
}

func (ms SortMap) Len() int {
	return len(ms)
}

func (ms SortMap) Less(i, j int) bool {
	return ms[i].Key < ms[j].Key // 按键值排序
}

func (ms SortMap) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func GetTime(date string) time.Time {
	if len(date) < 0 {
		return time.Now()
	}
	re, _ := regexp.Compile("(\\d){1,5}")
	maps := re.FindAllString(date, 5)

	if len(maps) >= 3 {
		year, _ := strconv.Atoi(maps[0])
		num, _ := strconv.Atoi(maps[1])
		mouth := time.Month(num)
		day, _ := strconv.Atoi(maps[2])
		hour := 0
		minute := 0
		if len(maps) > 3 {
			if len(maps[3]) > 0 {
				hour, _ = strconv.Atoi(maps[3])
			}
			if len(maps[4]) > 0 {
				minute, _ = strconv.Atoi(maps[4])
			}
		}
		time := time.Date(year, mouth, day, hour, minute, 0, 0, time.Local)
		return time
	}
	return time.Now()
}
