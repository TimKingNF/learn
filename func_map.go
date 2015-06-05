//	注册模板函数
package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	Edu "learn/EDU"
	"learn/models"
	"strconv"
	"strings"
	"time"
)

func RegisterFuncMap() {
	beego.AddFuncMap("getSex", func(n int) string {
		switch n {
		case 1:
			return "男"
		case 2:
			return "女"
		default:
			return "未知"
		}
	})

	beego.AddFuncMap("len", func(in interface{}) int {
		switch in.(type) {
		case []*models.StudentCourse:
			return len(in.([]*models.StudentCourse))
		case []string:
			return len(in.([]string))
		case []int:
			return len(in.([]int))
		case []*Edu.EduLog:
			return len(in.([]*Edu.EduLog))
		case []*models.Student:
			return len(in.([]*models.Student))
		case []*models.StudentHomework:
			return len(in.([]*models.StudentHomework))
		case []*models.StudentCheck:
			return len(in.([]*models.StudentCheck))
		case []*models.TeacherCourseHomework:
			return len(in.([]*models.TeacherCourseHomework))
		}
		return -1
	})

	beego.AddFuncMap("getKey", func(in int) int {
		return in + 1
	})

	beego.AddFuncMap("FilterRepeat", func(n interface{}) interface{} {
		switch n.(type) {
		case []string:
			return models.FilterRepeat(n).([]string)
		}
		return nil
	})

	beego.AddFuncMap("GetTermsYear", func(n []*models.Term) []string {
		if len(n) <= 0 {
			return nil
		}
		var result []string
		for _, v := range n {
			result = append(result, fmt.Sprintf("%d", v.StartYear)+"-"+fmt.Sprintf("%d", v.EndYear))
		}
		return result
	})

	beego.AddFuncMap("GetTermYear", func(n *models.Term) string {
		if n == nil {
			return ""
		}
		return fmt.Sprintf("%d", n.StartYear) + "-" + fmt.Sprintf("%d", n.EndYear)
	})

	beego.AddFuncMap("GetTermNumber", func(n *models.Term) string {
		if n == nil {
			return ""
		}
		return fmt.Sprintf("%d", n.Number)
	})

	beego.AddFuncMap("GetTermsNumber", func(n []*models.Term) []string {
		if len(n) <= 0 {
			return nil
		}
		var result []string
		for _, v := range n {
			result = append(result, fmt.Sprintf("%d", v.Number))
		}
		return result
	})

	beego.AddFuncMap("GetWeekTime", func(in string) []string {
		var time map[string]string
		if err := json.Unmarshal([]byte(in), &time); err != nil {
			return nil
		}
		return strings.Split(time["week_time"], ",")
	})

	beego.AddFuncMap("Map2Json", func(in string) map[string]string {
		var maps map[string]string
		if err := json.Unmarshal([]byte(in), &maps); err != nil {
			return nil
		}
		return maps
	})

	beego.AddFuncMap("GetWeekTimeSmall", func(in []string) int {
		if len(in) <= 0 {
			return 0
		}
		num1, _ := strconv.Atoi(in[0])
		small := num1
		for _, v := range in {
			num2, _ := strconv.Atoi(v)
			if small > num2 {
				small = num2
			}
		}
		return small
	})

	beego.AddFuncMap("GetCourseColor", func(key1, key2 int) int {
		if key1 > 13 {
			return key1 - key2
		} else if key1 > key2 {
			return key1 - key2
		} else if key1 < key2 {
			return key1 + key2
		} else if key1 > 3 {
			return key1 + key2
		}
		return 0
	})

	beego.AddFuncMap("WeekIsPart", func(now_week_string string, t_course *models.TeacherCourse) bool {
		if len(now_week_string) == 0 {
			if 1 >= t_course.StartWeek && 1 <= t_course.EndWeek {
				return true
			}
		}
		now_week, err := strconv.Atoi(now_week_string)
		if err != nil {
			return false
		}
		if now_week >= t_course.StartWeek && now_week <= t_course.EndWeek {
			return true
		}
		return false
	})

	beego.AddFuncMap("GetEduLogLen", func(userinfo map[string]string, edu_log ...[]*Edu.EduLog) int {
		var result = 0
		if userinfo != nil {
			result += 1
		}
		for _, v := range edu_log {
			result += len(v)
		}
		return result
	})

	beego.AddFuncMap("JudgeCourseSpacing", func(now_week_string string, in interface{}) string {
		switch in.(type) {
		case map[string]*models.StudentCourse:
			for _, v := range in.(map[string]*models.StudentCourse) {
				switch v.TeacherCourse.Spacing {
				case 1:
					now_week, _ := strconv.Atoi(now_week_string)
					if now_week%2 == 1 {
						return "单周"
					}
				case 2:
					now_week, _ := strconv.Atoi(now_week_string)
					if now_week%2 == 0 {
						return "双周"
					}
				case 0:
					return "单,双周"
				}
			}
			return ""
		case map[string]*models.TeacherCourse:
			for _, v := range in.(map[string]*models.TeacherCourse) {
				switch v.Spacing {
				case 1:
					now_week, _ := strconv.Atoi(now_week_string)
					if now_week%2 == 1 {
						return "单周"
					}
				case 2:
					now_week, _ := strconv.Atoi(now_week_string)
					if now_week%2 == 0 {
						return "双周"
					}
				case 0:
					return "单,双周"
				}
			}
			return ""
		}
		return ""
	})

	beego.AddFuncMap("GetWeekDayInChinese", func(in string) string {
		switch in {
		case "1":
			return "一"
		case "2":
			return "二"
		case "3":
			return "三"
		case "4":
			return "四"
		case "5":
			return "五"
		case "6":
			return "六"
		case "7":
			return "日"
		}
		return ""
	})

	beego.AddFuncMap("getDepartmentName", func(in *models.Department) string {
		if in == nil {
			return ""
		}
		if len(in.Name) > 0 {
			return in.Name
		} else {
			if tmp, err := models.GetDepartmentById(in.Id); err == nil {
				return tmp.Name
			} else {
				return ""
			}
		}
		return ""
	})

	beego.AddFuncMap("getCourseCredit", func(in []*models.StudentCourse) (result float64) {
		if in == nil {
			return 0
		}
		//	剔除 成绩与绩点均为 0 的 课程
		in = models.FilterStudentCourse(in)
		if in == nil {
			return 0
		}
		for k, _ := range in {
			if in[k].TeacherCourse != nil && in[k].TeacherCourse.Course != nil {
				result += in[k].TeacherCourse.Course.Credit
			}
		}
		return
	})

	beego.AddFuncMap("getCourseGPA", func(in []*models.StudentCourse) string {
		if in == nil {
			return "0"
		}
		//	剔除 成绩与绩点均为 0 的 课程
		in = models.FilterStudentCourse(in)
		if in == nil {
			return "0"
		}
		var result float64 = 0
		for k, _ := range in {
			if in[k].Score != nil {
				result += in[k].Score.GradePointAverage
			}
		}
		return fmt.Sprintf("%.2f", result)
	})

	beego.AddFuncMap("getCourseAverageGrade", func(in []*models.StudentCourse) (result int) {
		if in == nil {
			return 0
		}
		//	剔除 成绩与绩点均为 0 的 课程
		in = models.FilterStudentCourse(in)
		if in == nil {
			return 0
		}
		for k, _ := range in {
			if in[k].Score != nil {
				result += in[k].Score.Grade
			}
		}
		return result / len(in)
	})

	beego.AddFuncMap("getCourseAverageGPA", func(in []*models.StudentCourse) string {
		if in == nil {
			return "0"
		}
		//	剔除 成绩与绩点均为 0 的 课程
		in = models.FilterStudentCourse(in)
		if in == nil {
			return "0"
		}
		var result float64 = 0
		for k, _ := range in {
			if in[k].Score != nil {
				result += in[k].Score.GradePointAverage
			}
		}
		return fmt.Sprintf("%.2f", result/float64(len(in)))
	})

	beego.AddFuncMap("lowerYear", func(in int) int {
		return in - 2000
	})

	beego.AddFuncMap("getAllCourseCredit", func(in []*models.Term) (result float64) {
		if in == nil {
			return 0
		}
		var s_course []*models.StudentCourse
		for k, _ := range in {
			s_course = append(s_course, in[k].Orgs["s_course"].([]*models.StudentCourse)...)
		}
		//	剔除 成绩与绩点均为 0 的 课程
		s_course = models.FilterStudentCourse(s_course)
		if s_course == nil {
			return 0
		}
		for k, _ := range s_course {
			if s_course[k].TeacherCourse != nil && s_course[k].TeacherCourse.Course != nil {
				result += s_course[k].TeacherCourse.Course.Credit
			}
		}
		return
	})

	beego.AddFuncMap("getAllCourseGPA", func(in []*models.Term) string {
		if in == nil {
			return "0"
		}
		var s_course []*models.StudentCourse
		for k, _ := range in {
			s_course = append(s_course, in[k].Orgs["s_course"].([]*models.StudentCourse)...)
		}
		//	剔除 成绩与绩点均为 0 的 课程
		s_course = models.FilterStudentCourse(s_course)
		if s_course == nil {
			return "0"
		}
		var result float64 = 0
		for k, _ := range s_course {
			if s_course[k].Score != nil {
				result += s_course[k].Score.GradePointAverage
			}
		}
		return fmt.Sprintf("%.2f", result)
	})

	beego.AddFuncMap("getAllCourseAverageGrade", func(in []*models.Term) (result int) {
		if in == nil {
			return 0
		}
		var s_course []*models.StudentCourse
		for k, _ := range in {
			s_course = append(s_course, in[k].Orgs["s_course"].([]*models.StudentCourse)...)
		}
		//	剔除 成绩与绩点均为 0 的 课程
		s_course = models.FilterStudentCourse(s_course)
		if s_course == nil {
			return 0
		}
		for k, _ := range s_course {
			if s_course[k].Score != nil {
				result += s_course[k].Score.Grade
			}
		}
		return result / len(s_course)
	})

	beego.AddFuncMap("getAllCourseAverageGPA", func(in []*models.Term) string {
		if in == nil {
			return "0"
		}
		var s_course []*models.StudentCourse
		for k, _ := range in {
			s_course = append(s_course, in[k].Orgs["s_course"].([]*models.StudentCourse)...)
		}
		//	剔除 成绩与绩点均为 0 的 课程
		s_course = models.FilterStudentCourse(s_course)
		if s_course == nil {
			return "0"
		}
		var result float64 = 0
		for k, _ := range s_course {
			if s_course[k].Score != nil {
				result += s_course[k].Score.GradePointAverage
			}
		}
		return fmt.Sprintf("%.2f", result/float64(len(s_course)))
	})

	beego.AddFuncMap("isOutOfDate", func(in string) bool {
		in_time := models.GetTime(in)
		if in_time.After(time.Now()) {
			return true
		}
		return false
	})

	beego.AddFuncMap("getHomeworkGrade", func(in []*models.StudentHomework) []string {
		var a, b, c, d, e = 0, 0, 0, 0, 0
		for k, _ := range in {
			switch in[k].Grade {
			case "A":
				a++
			case "B":
				b++
			case "C":
				c++
			case "D":
				d++
			case "E":
				e++
			}
		}
		return strings.Split(fmt.Sprintf("%d", a)+","+fmt.Sprintf("%d", b)+","+fmt.Sprintf("%d", c)+","+fmt.Sprintf("%d", d)+","+fmt.Sprintf("%d", e), ",")
	})

	beego.AddFuncMap("getHomeworkByStudentCourses", func(in []*models.StudentCourse) (s_homeworks []*models.StudentHomework) {
		if len(in) <= 0 {
			return nil
		}
		for k, _ := range in {
			if in[k].Orgs["s_homeworks"] != nil {
				s_homeworks = append(s_homeworks, in[k].Orgs["s_homeworks"].([]*models.StudentHomework)...)
			}
		}
		return
	})

	beego.AddFuncMap("getHomeworkByTerm", func(in []*models.Term) (s_courses []*models.StudentCourse) {
		if len(in) <= 0 {
			return nil
		}
		for k, _ := range in {
			if in[k].Orgs["s_course"] != nil {
				s_courses = append(s_courses, in[k].Orgs["s_course"].([]*models.StudentCourse)...)
			}
		}
		return
	})

	beego.AddFuncMap("getCheckByStudentCourses", func(in []*models.StudentCourse) (s_checkes []*models.StudentCheck) {
		if len(in) <= 0 {
			return nil
		}
		for k, _ := range in {
			if in[k].Orgs["s_checks"] != nil {
				s_checkes = append(s_checkes, in[k].Orgs["s_checks"].([]*models.StudentCheck)...)
			}
		}
		return
	})

	beego.AddFuncMap("getCheckByTerm", func(in []*models.Term) (s_courses []*models.StudentCourse) {
		if len(in) <= 0 {
			return nil
		}
		for k, _ := range in {
			if in[k].Orgs["s_course"] != nil {
				s_courses = append(s_courses, in[k].Orgs["s_course"].([]*models.StudentCourse)...)
			}
		}
		return
	})

	beego.AddFuncMap("getWeekByStudentCheck", func(in []*models.StudentCheck) (weeks []int) {
		if len(in) <= 0 {
			return nil
		}
		for k, _ := range in {
			if in[k] != nil {
				weeks = append(weeks, in[k].Week)
			}
		}
		return models.FilterRepeat(weeks).([]int)
	})

	beego.AddFuncMap("getStudentCheckInChineseByWeek", func(in []*models.StudentCheck, week int) []string {
		if len(in) <= 0 || week <= 0 || week > 20 {
			return nil
		}
		var a, b, c, d = 0, 0, 0, 0
		for k, _ := range in {
			if in[k] != nil && in[k].Week == week {
				switch in[k].Result {
				case "未到":
					a++
				case "迟到":
					b++
				case "已到":
					c++
				case "请假":
					d++
				}
			}
		}
		return strings.Split(fmt.Sprintf("%d", a)+","+fmt.Sprintf("%d", b)+","+fmt.Sprintf("%d", c)+","+fmt.Sprintf("%d", d), ",")
	})

	beego.AddFuncMap("getStudentCheckInChinese", func(in []*models.StudentCheck) []string {
		var a, b, c, d = 0, 0, 0, 0
		for k, _ := range in {
			if in[k] != nil {
				switch in[k].Result {
				case "未到":
					a++
				case "迟到":
					b++
				case "已到":
					c++
				case "请假":
					d++
				}
			}
		}
		return strings.Split(fmt.Sprintf("%d", a)+","+fmt.Sprintf("%d", b)+","+fmt.Sprintf("%d", c)+","+fmt.Sprintf("%d", d), ",")
	})

	beego.AddFuncMap("getStudentHomeworksByTeacherCourseHomeworks", func(in []*models.TeacherCourseHomework) (result []*models.StudentHomework) {
		if len(in) <= 0 {
			return nil
		}
		for k, _ := range in {
			result = append(result, in[k].StudentHomeworks...)
		}
		return
	})
}
