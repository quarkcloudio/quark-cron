package resource

import (
	"html"

	"github.com/quarkcms/quark-cron/model"
	"github.com/quarkcms/quark-cron/search"
	"github.com/quarkcms/quark-go/pkg/app/handler/admin/actions"
	"github.com/quarkcms/quark-go/pkg/app/handler/admin/searches"
	"github.com/quarkcms/quark-go/pkg/builder"
	"github.com/quarkcms/quark-go/pkg/builder/template/adminresource"
	"gorm.io/gorm"
)

type JobLog struct {
	adminresource.Template
}

// 初始化
func (p *JobLog) Init() interface{} {

	// 初始化模板
	p.TemplateInit()

	// 标题
	p.Title = "作业日志"

	// 模型
	p.Model = &model.JobLog{}

	// 分页
	p.PerPage = 10

	// 重写排序，兼容sqlite
	p.IndexOrder = "job_logs.id desc"

	// 自动刷新
	p.IndexPolling = 10

	return p
}

// 查询条件
func (p *JobLog) IndexQuery(ctx *builder.Context, query *gorm.DB) *gorm.DB {
	return query.
		Select("job_logs.id,job_logs.schedule_id,job_logs.job_id,job_logs.result,job_logs.status,job_logs.created_at").
		Joins("LEFT JOIN jobs ON job_logs.job_id = jobs.id")
}

func (p *JobLog) Fields(ctx *builder.Context) []interface{} {
	field := &adminresource.Field{}

	return []interface{}{
		field.ID("id", "ID"),

		field.Select("schedule_id", "调度器", func() interface{} {
			schedulerInfo := (&model.Scheduler{}).GetInfoById(p.Field["schedule_id"])

			return schedulerInfo.Name
		}),

		field.Select("job_id", "作业名称", func() interface{} {
			jobInfo := (&model.Job{}).GetInfoById(p.Field["job_id"])

			return jobInfo.Name
		}),

		field.Datetime("created_at", "执行时间"),

		field.Switch("status", "状态").
			SetSpan(2).
			SetFalseValue("失败").
			SetTrueValue("成功"),

		field.Text("result", "执行结果", func() interface{} {

			return "<pre>" + html.EscapeString(p.Field["result"].(string)) + "</pre>"
		}).
			SetSpan(2).
			OnlyOnDetail(),
	}
}

// 搜索
func (p *JobLog) Searches(ctx *builder.Context) []interface{} {

	return []interface{}{
		(&searches.Input{}).Init("jobs.name", "作业名称"),
		(&search.JobLogStatus{}).Init(),
	}
}

// 行为
func (p *JobLog) Actions(ctx *builder.Context) []interface{} {

	return []interface{}{
		(&actions.Delete{}).Init("批量删除"),
		(&actions.DetailLink{}).Init("详情"),
		(&actions.Delete{}).Init("删除"),
		(&actions.FormBack{}).Init(),
		(&actions.FormExtraBack{}).Init(),
	}
}
