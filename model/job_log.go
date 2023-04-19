package model

import (
	"time"

	appmodel "github.com/quarkcms/quark-go/pkg/app/model"
	"github.com/quarkcms/quark-go/pkg/dal/db"
	"gorm.io/gorm"
)

// 作业日志模型
type JobLog struct {
	Id         int            `json:"id" gorm:"autoIncrement"`
	ScheduleId int            `json:"schedule_id" gorm:"size:11;not null;default:0"`
	JobId      int            `json:"job_id" gorm:"size:11;not null;default:0"`
	Result     string         `json:"result" gorm:"size:5000;null"`
	Status     int            `json:"status" gorm:"size:4;not null;default:1"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

// Seeder
func (m *JobLog) Seeder() {

	// 如果菜单已存在，不执行Seeder操作
	if (&appmodel.Menu{}).IsExist(21) {
		return
	}

	// 创建菜单
	menuSeeders := []*appmodel.Menu{
		{Id: 21, Name: "作业日志", GuardName: "admin", Icon: "", Type: "engine", Pid: 18, Sort: 0, Path: "/api/admin/jobLog/index", Show: 1, Status: 1},
	}
	db.Client.Create(&menuSeeders)
}

// 插入数据
func (m *JobLog) Insert(jobLog *JobLog) {
	db.Client.Create(&jobLog)
}
