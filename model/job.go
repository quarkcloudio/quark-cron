package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/axgle/mahonia"
	"github.com/go-co-op/gocron"
	appmodel "github.com/quarkcms/quark-go/pkg/app/model"
	"github.com/quarkcms/quark-go/pkg/dal/db"
	"gorm.io/gorm"
)

// 作业模型
type Job struct {
	Id         int            `json:"id" gorm:"autoIncrement"`
	Name       string         `json:"name" gorm:"size:200;not null"`
	ScheduleId int            `json:"schedule_id" gorm:"size:11;not null;default:0"`
	Type       int            `json:"type" gorm:"size:4;not null;default:1"`
	CycleType  int            `json:"cycle_type" gorm:"size:4;not null;default:1"`
	Week       int            `json:"week" gorm:"size:11;not null;default:1"`
	Day        int            `json:"day" gorm:"size:11;not null;default:3"`
	Hour       int            `json:"hour" gorm:"size:11;not null;default:1"`
	Minute     int            `json:"minute" gorm:"size:11;not null;default:30"`
	Second     int            `json:"second" gorm:"size:11;not null;default:30"`
	Shell      string         `json:"shell" gorm:"size:5000;null"`
	Url        string         `json:"url" gorm:"size:1000;null"`
	Path       string         `json:"path" gorm:"size:1000;null"`
	Status     int            `json:"status" gorm:"size:4;not null;default:1"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

// Seeder
func (m *Job) Seeder() {

	// 如果菜单已存在，不执行Seeder操作
	if (&appmodel.Menu{}).IsExist(20) {
		return
	}

	// 创建菜单
	menuSeeders := []*appmodel.Menu{
		{Id: 20, Name: "作业列表", GuardName: "admin", Icon: "", Type: "engine", Pid: 18, Sort: 0, Path: "/api/admin/job/index", Show: 1, Status: 1},
	}
	db.Client.Create(&menuSeeders)
}

// 插入数据
func (m *Job) Insert(job *Job) {
	db.Client.Create(&job)
}

// 获取作业信息
func (m *Job) GetInfoById(id interface{}) (job *Job) {
	db.Client.Where("id = ?", id).Find(&job)

	return
}

// 加载调度器的所有Job
func (m *Job) LoadServices(schedulerId int, scheduler *gocron.Scheduler) (err error) {
	var jobs []*Job
	var job *gocron.Job
	db.Client.Where("schedule_id", schedulerId).Where("status = ?", 1).Find(&jobs)

	for _, v := range jobs {
		switch v.CycleType {
		case 1:
			// 限定每天的几点几分几秒执行一次
			job, err = scheduler.
				Every(1).
				Day().
				At(strconv.Itoa(v.Hour)+":"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				Do(execute, v)
		case 2:
			// 每隔几天几小时几分几秒执行一次
			job, err = scheduler.
				Every((v.Day*24*60*60)+(v.Hour*60*60)+(v.Minute*60)+v.Second).
				Seconds().
				Do(execute, v)
		case 3:
			// 限定每小时的几分几秒执行一次
			job, err = scheduler.
				Every(1).
				Day().
				At("01:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("02:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("03:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("04:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("05:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("06:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("07:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("08:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("09:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("10:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("11:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("12:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("13:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("14:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("15:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("16:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("17:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("18:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("19:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("20:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("21:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("22:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				At("23:"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				Do(execute, v)
		case 4:
			// 每隔几小时几分几秒执行一次
			job, err = scheduler.
				Every((v.Hour*60*60)+(v.Minute*60)+v.Second).
				Seconds().
				Do(execute, v)
		case 5:
			// 每隔几分几秒执行一次
			job, err = scheduler.
				Every((v.Minute*60)+v.Second).
				Seconds().
				Do(execute, v)
		case 6:
			// 每隔几秒钟执行一次
			job, err = scheduler.
				Every(v.Second).
				Seconds().
				Do(execute, v)
		case 7:
			// 每周的几点几分几秒执行一次
			switch v.Week {
			case 1:
				// 周一
				job, err = scheduler.
					Every(1).
					Monday().
					At(strconv.Itoa(v.Hour)+":"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
					Do(execute, v)
			case 2:
				// 周二
				job, err = scheduler.
					Every(1).
					Thursday().
					At(strconv.Itoa(v.Hour)+":"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
					Do(execute, v)
			case 3:
				// 周三
				job, err = scheduler.
					Every(1).
					Wednesday().
					At(strconv.Itoa(v.Hour)+":"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
					Do(execute, v)
			case 4:
				// 周四
				job, err = scheduler.
					Every(1).
					Thursday().
					At(strconv.Itoa(v.Hour)+":"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
					Do(execute, v)
			case 5:
				// 周五
				job, err = scheduler.
					Every(1).
					Friday().
					At(strconv.Itoa(v.Hour)+":"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
					Do(execute, v)
			case 6:
				// 周六
				job, err = scheduler.
					Every(1).
					Saturday().
					At(strconv.Itoa(v.Hour)+":"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
					Do(execute, v)
			case 7:
				// 周日
				job, err = scheduler.
					Every(1).
					Sunday().
					At(strconv.Itoa(v.Hour)+":"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
					Do(execute, v)
			}
		case 8:
			// 每月几号几点几分执行一次
			job, err = scheduler.
				Every(1).
				Month(v.Day).
				At(strconv.Itoa(v.Hour)+":"+strconv.Itoa(v.Minute)+":"+strconv.Itoa(v.Second)).
				Do(execute, v)
		}

		// 返回错误信息
		if err != nil {
			return err
		}

		// 标记作业
		job.Tag(strconv.Itoa(v.Id))
	}

	return
}

// 执行
func execute(job *Job) {
	var (
		shellPath = "./shell/"
		result    = ""
		status    = 1
	)

	// 获取调度器信息
	schedulerInfo := (&Scheduler{}).GetInfoById(job.ScheduleId)

	// 日志
	// log := "调度器:" + schedulerInfo.Name + " 作业:" + job.Name + " 执行时间:" + time.Now().Format(time.DateTime)
	// fmt.Println(log)

	switch job.Type {
	case 1:
		// Shell脚本
		hashContent := strconv.Itoa(schedulerInfo.Id) +
			strconv.Itoa(job.Id) +
			job.CreatedAt.Format(time.DateTime) +
			job.UpdatedAt.Format(time.DateTime) +
			job.Shell

		// 计算脚本hash
		hashValue, err := hash(hashContent)
		if err != nil {
			fmt.Println(err)
			return
		}

		// shell文件夹是否存在，不存在则创建
		if !isExist(shellPath) {
			os.Mkdir(shellPath, 0777)
		}

		// 判断操作系统类型
		if runtime.GOOS == "windows" {
			shellPath = shellPath + hashValue + ".bat"
		} else {
			shellPath = shellPath + hashValue + ".sh"
		}

		// shell文件是否存在
		if !isExist(shellPath) {
			f, err := os.OpenFile(shellPath, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0777)
			if err != nil {
				fmt.Println(err)
				return
			}

			// 关闭文件
			defer f.Close()

			// 写入内容
			_, err = f.WriteString(job.Shell)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		// 获取shell的绝对路径
		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}

		// 判断操作系统类型
		if runtime.GOOS == "windows" {
			output, err := exec.Command("cmd.exe", "/C", path+"\\shell\\"+hashValue+".bat").Output()
			if err != nil {
				status = 0
				result = err.Error()
			} else {
				result = string(output)
			}
		} else {
			output, err := exec.Command("sh", "/C", path+"/shell/"+hashValue+".sh").Output()
			if err != nil {
				status = 0
				result = err.Error()
			} else {
				result = string(output)
			}
		}
	case 2:
		// 访问URL
		resp, err := http.Get(job.Url)
		if err != nil {
			status = 0
			result = err.Error()
			return
		} else {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				status = 0
				result = err.Error()
			} else {
				result = string(body)
			}
		}
	case 3:
		// 执行应用
		output, err := exec.Command(job.Path).Output()
		if err != nil {
			status = 0
			result = err.Error()
		} else {
			result = string(output)
		}
	}

	decoder := mahonia.NewDecoder("gbk")
	utf8Result := decoder.ConvertString(result)

	// 插入日志
	(&JobLog{}).Insert(&JobLog{
		ScheduleId: schedulerInfo.Id,
		JobId:      job.Id,
		Result:     string(utf8Result),
		Status:     status,
	})
}

// 计算文件哈希值
func hash(content string) (value string, err error) {
	sha256New := sha256.New()
	byteReader := bytes.NewReader([]byte(content))
	_, err = io.Copy(sha256New, byteReader)
	if err != nil {
		return
	}
	value = hex.EncodeToString(sha256New.Sum(nil))
	return
}

// 检查路径是否存在
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
