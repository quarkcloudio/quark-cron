package model

import (
	"fmt"
	"strconv"
	"time"

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
				At(strconv.Itoa(v.Hour) + ":" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				Do(func() {
					fmt.Println("执行周期：" + strconv.Itoa(v.Hour) + "点" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
				})
		case 2:
			// 每隔几天几小时几分几秒执行一次
			job, err = scheduler.
				Every((v.Day * 24 * 60 * 60) + (v.Hour * 60 * 60) + (v.Minute * 60) + v.Second).
				Seconds().
				Do(func() {
					fmt.Println("执行周期：" + strconv.Itoa(v.Hour) + "点" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
				})
		case 3:
			// 限定每小时的几分几秒执行一次
			job, err = scheduler.
				Every(1).
				Day().
				At("01:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("02:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("03:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("04:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("05:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("06:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("07:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("08:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("09:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("10:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("11:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("12:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("13:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("14:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("15:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("16:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("17:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("18:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("19:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("20:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("21:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("22:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				At("23:" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				Do(func() {
					fmt.Println("执行周期：" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
				})
		case 4:
			// 每隔几小时几分几秒执行一次
			job, err = scheduler.
				Every((v.Hour * 60 * 60) + (v.Minute * 60) + v.Second).
				Seconds().
				Do(func() {
					fmt.Println("执行周期：" + strconv.Itoa(v.Hour) + "点" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
				})
		case 5:
			// 每隔几分几秒执行一次
			job, err = scheduler.
				Every((v.Minute * 60) + v.Second).
				Seconds().
				Do(func() {
					fmt.Println("执行周期：" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
				})
		case 6:
			// 每隔几秒钟执行一次
			job, err = scheduler.
				Every(v.Second).
				Seconds().
				Do(func() {
					fmt.Println("执行周期：" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
				})
		case 7:
			// 每周的几点几分几秒执行一次
			switch v.Week {
			case 1:
				job, err = scheduler.
					Every(1).
					Monday().
					At(strconv.Itoa(v.Hour) + ":" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
					Do(func() {
						fmt.Println("执行周期：" + "周一" + strconv.Itoa(v.Hour) + "点" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
					})
			case 2:
				job, err = scheduler.
					Every(1).
					Thursday().
					At(strconv.Itoa(v.Hour) + ":" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
					Do(func() {
						fmt.Println("执行周期：" + "周二" + strconv.Itoa(v.Hour) + "点" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
					})
			case 3:
				job, err = scheduler.
					Every(1).
					Wednesday().
					At(strconv.Itoa(v.Hour) + ":" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
					Do(func() {
						fmt.Println("执行周期：" + "周三" + strconv.Itoa(v.Hour) + "点" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
					})
			case 4:
				job, err = scheduler.
					Every(1).
					Thursday().
					At(strconv.Itoa(v.Hour) + ":" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
					Do(func() {
						fmt.Println("执行周期：" + "周四" + strconv.Itoa(v.Hour) + "点" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
					})
			case 5:
				job, err = scheduler.
					Every(1).
					Friday().
					At(strconv.Itoa(v.Hour) + ":" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
					Do(func() {
						fmt.Println("执行周期：" + "周五" + strconv.Itoa(v.Hour) + "点" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
					})
			case 6:
				job, err = scheduler.
					Every(1).
					Saturday().
					At(strconv.Itoa(v.Hour) + ":" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
					Do(func() {
						fmt.Println("执行周期：" + "周六" + strconv.Itoa(v.Hour) + "点" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
					})
			case 7:
				job, err = scheduler.
					Every(1).
					Sunday().
					At(strconv.Itoa(v.Hour) + ":" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
					Do(func() {
						fmt.Println("执行周期：" + "周日" + strconv.Itoa(v.Hour) + "点" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
					})
			}
		case 8:
			// 每月几号几点几分执行一次
			job, err = scheduler.
				Every(1).
				Month(v.Day).
				At(strconv.Itoa(v.Hour) + ":" + strconv.Itoa(v.Minute) + ":" + strconv.Itoa(v.Second)).
				Do(func() {
					fmt.Println("执行周期：" + strconv.Itoa(v.Day) + "号" + strconv.Itoa(v.Hour) + "点" + strconv.Itoa(v.Minute) + "分" + strconv.Itoa(v.Second))
				})
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
