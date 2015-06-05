package EDU

import (
	"errors"
	"fmt"
	"learn/models"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func getUserScoreHref(data string) (href string, err error) {
	if len(data) <= 0 {
		return "", errors.New("getUserScoreHref: data is null")
	}
	re, _ := regexp.Compile(`学生考试查询</a></li><li><a href="(.*)" target='zhuti' onclick=".*">成绩查询</a></li>`)
	maps := re.FindStringSubmatch(data)
	if len(maps) > 0 && len(maps[0]) > 0 {
		return maps[1], nil
	}
	return "", errors.New("getUserScoreHref: not find string")
}

func getUserScore(data string) (info map[string]map[string]string, err error) {
	if len(data) <= 0 {
		return nil, errors.New("getUserScore: data is null")
	}
	//	取出table内容
	re, _ := regexp.Compile(`<table class="datelist" cellspacing="0" cellpadding="3" border="0" id="Datagrid1" style="width:100%;border-collapse:collapse;">([\s\S]*)</table>[\s\S]*<table width="100%">`)
	maps := re.FindStringSubmatch(data)
	if len(maps) <= 0 {
		return nil, errors.New("getUserScore: table is null")
	}
	//	去除所有html标签
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	data = re.ReplaceAllString(maps[1], "\n")
	data = strings.Replace(data, "&nbsp;", "-", -1)
	//	去除多余换行符 和空格
	re, _ = regexp.Compile("\\s{2,}")
	data = re.ReplaceAllString(data, "\n")
	//	根据换行符分割data
	arr := strings.Split(data, "\n")
	if len(arr) < 0 {
		return nil, errors.New("getUserScore: data is null")
	}
	arr = arr[1:]
	info = make(map[string]map[string]string, len(arr)-1)
	for k, _ := range arr {
		if k >= 15 && k%15 == 0 {
			k1 := k + 15
			if k1 >= len(arr) {
				break
			}
			var s = make(map[string]string, k1-k)
			for i := 0; k < k1; k++ {
				s[arr[i]] = arr[k]
				i++
			}
			info[s["课程代码"]] = s
		}
	}
	return info, nil
}

func GetStudentScore(user *models.Student, data []byte, cookies []*http.Cookie) (map[string]map[string]string, error) {
	userscorehref, _ := getUserScoreHref(string(data))
	userscore_param, _, err := client.SendRequest("GET", Host_url+"/"+userscorehref, "", cookies, header)
	if err != nil {
		return nil, models.ErrorInfo("GetStudentScore", errors.New("连接失败,错误原因："+err.Error()))
	}
	vierstate, eventvalidation, _ := getLoginKey(string(userscore_param))
	userscore_info, _, err := client.SendRequest("POST", Host_url+"/"+userscorehref, url.Values{
		"__VIEWSTATE":       {vierstate},
		"__EVENTVALIDATION": {eventvalidation},
		"ddIXN":             {""},
		"ddIXQ":             {""},
		"Button1":           {models.Utf82gbk("按学期查询")},
	}.Encode(), cookies, header)
	if err != nil {
		return nil, models.ErrorInfo("GetStudentScore", errors.New("连接失败,错误原因："+err.Error()))
	}
	userscore, _ := getUserScore(string(userscore_info))
	for _, v := range userscore {
		if len(v) > 0 {
			if exist := models.CourseExist(v["课程名称"]); !exist {
				if exist := models.DepartmentExist(v["学院名称"]); !exist {
					if _, err := models.AddDepartment(&models.Department{Name: v["学院名称"]}); err == nil {
						EduLogCreate(&EduLog{Student: user, Content: v["学院名称"], Type: "department", Result: 1})
					} else {
						EduLogCreate(&EduLog{Student: user, Content: v["学院名称"], Type: "department", Result: 0})
					}
				}
				if department, err := models.GetDepartmentByName(v["学院名称"]); err == nil {
					f, _ := strconv.ParseFloat(v["学分"], 64)
					if err := models.AddCourse(&models.Course{Id: v["课程代码"], Name: v["课程名称"], Department: department, Remark: v["课程归属"], Type: v["课程性质"], Credit: f}); err == nil {
						EduLogCreate(&EduLog{Student: user, Content: v["课程名称"], Type: "course", Result: 1})
					} else {
						EduLogCreate(&EduLog{Student: user, Content: v["课程名称"], Type: "course", Result: 0})
					}
				}
			}
		}
	}
	return userscore, nil
}

func ControlSQLByScheduleAndScore(user *models.Student, schedules map[*models.Term][]*models.TeacherCourse, score map[string]map[string]string) error {
	if user == nil || len(schedules) <= 0 || len(score) <= 0 {
		return models.ErrorInfo("ControlSQLByScheduleAndScore", errors.New("user is nil"))
	}
	for k, v := range schedules {
		//	取出学年度 学期
		if k == nil || v == nil {
			continue
		}
		if exist := models.TermExist(k.Number, k.StartYear, k.EndYear); !exist {
			if err := models.AddTerm(k); err == nil {
				EduLogCreate(&EduLog{Student: user, Content: `{"学年度": "` + fmt.Sprintf("%d", k.StartYear) + `-` + fmt.Sprintf("%d", k.EndYear) + `", "学期": "` + fmt.Sprintf("%d", k.Number) + `"}`, Type: "term", Result: 1})
			} else {
				EduLogCreate(&EduLog{Student: user, Content: `{"学年度": "` + fmt.Sprintf("%d", k.StartYear) + `-` + fmt.Sprintf("%d", k.EndYear) + `", "学期": "` + fmt.Sprintf("%d", k.Number) + `"}`, Type: "term", Result: 0})
			}
		}
		_, err := models.SearchTerm(k.Number, k.StartYear, k.EndYear)
		if err != nil {
			models.Info("ControlSQLByScheduleAndScore", err)
			continue
		}
		for _, t_course := range v {
			if exist := models.TeacherCourseExist(t_course.Course.Id, t_course.Time, t_course.Term.Id, t_course.Teacher.Id); !exist {
				//	add teacher course
				if err := models.AddTeacherCourse(t_course); err == nil {
					EduLogCreate(&EduLog{Student: user, Content: `教师 ` + t_course.Teacher.Name + ` 教授的课程：` + t_course.Course.Name + `, 上课时间：` + t_course.Time + `, 上课地点：` + t_course.Place, Type: "t_course", Result: 1})
				} else {
					models.Info("EduGetStudentSchedule", err)
					EduLogCreate(&EduLog{Student: user, Content: `教师 ` + t_course.Teacher.Name + ` 教授的课程：` + t_course.Course.Name + `, 上课时间：` + t_course.Time + `, 上课地点：` + t_course.Place, Type: "t_course", Result: 0})
				}
			}
			if temp_t_course, err := models.SearchTeacherCourse(t_course.Course.Id, t_course.Time, t_course.Term.Id, t_course.Teacher.Id); err == nil {
				if exist := models.StudentCourseExist(user.Id, temp_t_course.Id); !exist {
					//	add student course
					if err := models.AddStudentCourse(&models.StudentCourse{Student: user, TeacherCourse: temp_t_course}); err == nil {
						EduLogCreate(&EduLog{Student: user, Content: `教师 ` + temp_t_course.Teacher.Name + ` 教授的课程：` + temp_t_course.Course.Name + `, 上课时间：` + temp_t_course.Time + `, 上课地点：` + temp_t_course.Place, Type: "s_course", Result: 1})
					} else {
						EduLogCreate(&EduLog{Student: user, Content: `教师 ` + temp_t_course.Teacher.Name + ` 教授的课程：` + temp_t_course.Course.Name + `, 上课时间：` + temp_t_course.Time + `, 上课地点：` + temp_t_course.Place, Type: "s_course", Result: 0})
					}
				}
				if temp_s_course, err := models.SearchStudentCourse(user.Id, temp_t_course.Id); err == nil {
					if exist := models.StudentGradeExist(temp_s_course.Id); !exist {
						//	取出成绩
						var (
							score = score[t_course.Course.Id]
							grade = 0
						)
						if score["重修成绩"] != "-" {
							num, _ := strconv.Atoi(score["重修成绩"])
							grade = num
						} else if score["补考成绩"] != "-" {
							num, _ := strconv.Atoi(score["补考成绩"])
							grade = num
						} else if len(score["成绩"]) > 0 {
							num, _ := strconv.Atoi(score["成绩"])
							grade = num
						}
						gpa, _ := strconv.ParseFloat(score["绩点"], 64)
						if err := models.AddStudentGrade(&models.StudentGrade{StudentCourse: temp_s_course, GradePointAverage: gpa, Grade: grade}); err == nil {
							EduLogCreate(&EduLog{Student: user, Content: `科目：` + temp_s_course.TeacherCourse.Course.Name + `的成绩为：` + fmt.Sprintf("%d", grade) + `，绩点为：` + fmt.Sprintf("%.2f", gpa), Type: "s_course_grade", Result: 1})
						} else {
							EduLogCreate(&EduLog{Student: user, Content: `科目：` + temp_s_course.TeacherCourse.Course.Name + `的成绩为：` + fmt.Sprintf("%d", grade) + `，绩点为：` + fmt.Sprintf("%.2f", gpa), Type: "s_course_grade", Result: 0})
						}
					} else {
						//	取出成绩
						var (
							score = score[t_course.Course.Id]
							grade = 0
						)
						if score["重修成绩"] != "-" {
							num, _ := strconv.Atoi(score["重修成绩"])
							grade = num
						} else if score["补考成绩"] != "-" {
							num, _ := strconv.Atoi(score["补考成绩"])
							grade = num
						} else if len(score["成绩"]) > 0 {
							num, _ := strconv.Atoi(score["成绩"])
							grade = num
						}
						gpa, _ := strconv.ParseFloat(score["绩点"], 64)
						//	查找成绩
						if student_grades, err := models.SearchStudentGradeByStudentCourse(temp_s_course); err == nil {
							for _, v := range student_grades {
								//	更新
								models.UpdateStudentGrade(&models.StudentGrade{Id: v.Id, StudentCourse: temp_s_course, GradePointAverage: gpa, Grade: grade})
							}
						}
					}
				}
			}
		}
	}
	return nil
}
