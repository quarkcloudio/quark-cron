package model

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	appmodel "github.com/quarkcms/quark-go/pkg/app/model"
	"github.com/quarkcms/quark-go/pkg/component/admin/form/fields/selectfield"
	"github.com/quarkcms/quark-go/pkg/dal/db"
	"gorm.io/gorm"
)

// 调度器模型
type Scheduler struct {
	Id        int               `json:"id" gorm:"autoIncrement"`
	Name      string            `json:"name" gorm:"size:200;not null"`
	Status    int               `json:"status" gorm:"size:4;not null;default:1"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt gorm.DeletedAt    `json:"deleted_at"`
	Handler   *gocron.Scheduler `json:"-" gorm:"-"`
}

// Seeder
func (m *Scheduler) Seeder() {

	// 如果菜单已存在，不执行Seeder操作
	if (&appmodel.Menu{}).IsExist(18) {
		return
	}

	// 创建菜单
	menuSeeders := []*appmodel.Menu{
		{Id: 18, Name: "任务管理", GuardName: "admin", Icon: "icon-comment", Type: "default", Pid: 0, Sort: 0, Path: "/task", Show: 1, Status: 1},
		{Id: 19, Name: "调度列表", GuardName: "admin", Icon: "", Type: "engine", Pid: 18, Sort: 0, Path: "/api/admin/scheduler/index", Show: 1, Status: 1},
	}
	db.Client.Create(&menuSeeders)
}

// 插入数据
func (m *Scheduler) Insert(schedule *Scheduler) {
	db.Client.Create(&schedule)
}

// 属性值
func (m *Scheduler) Options() (options []*selectfield.Option) {
	var schedulers []*Scheduler
	db.Client.Find(&schedulers)

	for _, v := range schedulers {
		options = append(options, &selectfield.Option{
			Label: v.Name,
			Value: v.Id,
		})
	}

	return options
}

// 获取调度器信息
func (m *Scheduler) GetInfoById(id interface{}) (scheduler *Scheduler) {
	db.Client.Where("id = ?", id).Find(&scheduler)

	return
}

var schedulers []*Scheduler

// 加载所有调度器
func (m *Scheduler) LoadServices() (err error) {
	db.Client.Where("status = ?", 1).Find(&schedulers)
	for k, v := range schedulers {

		// 创建实例
		scheduler := gocron.NewScheduler(time.UTC)

		// 加载当前调度器的所有作业
		err = (&Job{}).LoadServices(v.Id, scheduler)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 启动调度器
		scheduler.StartAsync()

		// 将调度器挂载到全局变量
		v.Handler = scheduler

		schedulers[k] = v
	}

	return
}

// 清理所有调度器
func (m *Scheduler) ClearServices() (err error) {
	for _, v := range schedulers {
		v.Handler.Clear()
	}

	return
}

// 重新加载所有调度器
func (m *Scheduler) ReloadServices() (err error) {

	// 先清理所有调度器
	m.ClearServices()

	// 查询当前调度器
	db.Client.Where("status = ?", 1).Find(&schedulers)

	// 遍历处理
	for k, v := range schedulers {

		// 创建实例
		scheduler := gocron.NewScheduler(time.UTC)

		// 加载当前调度器的所有作业
		err = (&Job{}).LoadServices(v.Id, scheduler)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 启动调度器
		scheduler.StartAsync()

		// 将调度器挂载到全局变量
		v.Handler = scheduler

		schedulers[k] = v
	}

	return
}
