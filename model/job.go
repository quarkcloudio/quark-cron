package model

import (
	"time"

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
