package resource

import (
	"github.com/quarkcms/quark-cron/model"
	"github.com/quarkcms/quark-cron/search"
	"github.com/quarkcms/quark-go/pkg/app/handler/admin/actions"
	"github.com/quarkcms/quark-go/pkg/app/handler/admin/searches"
	"github.com/quarkcms/quark-go/pkg/builder"
	"github.com/quarkcms/quark-go/pkg/builder/template/adminresource"
	"github.com/quarkcms/quark-go/pkg/component/admin/form/rule"
	"gorm.io/gorm"
)

type Scheduler struct {
	adminresource.Template
}

// 初始化
func (p *Scheduler) Init() interface{} {

	// 初始化模板
	p.TemplateInit()

	// 标题
	p.Title = "调度器"

	// 模型
	p.Model = &model.Scheduler{}

	// 分页
	p.PerPage = 10

	return p
}

func (p *Scheduler) Fields(ctx *builder.Context) []interface{} {
	field := &adminresource.Field{}

	return []interface{}{
		field.ID("id", "ID"),

		field.Text("name", "名称").
			SetRules([]*rule.Rule{
				rule.Required(true, "名称必须填写"),
			}),

		field.Datetime("created_at", "创建时间").OnlyOnIndex(),

		field.Switch("status", "状态").
			SetFalseValue("已停止").
			SetTrueValue("运行中").
			SetEditable(true).
			OnlyOnIndex(),
	}
}

// 搜索
func (p *Scheduler) Searches(ctx *builder.Context) []interface{} {

	return []interface{}{
		(&searches.Input{}).Init("name", "名称"),
		(&search.Status{}).Init(),
	}
}

// 行为
func (p *Scheduler) Actions(ctx *builder.Context) []interface{} {

	return []interface{}{
		(&actions.Delete{}).Init("批量删除"),
		(&actions.Disable{}).Init("批量停止"),
		(&actions.Enable{}).Init("批量启动"),
		(&actions.CreateModal{}).Init(p.Title),
		(&actions.EditModal{}).Init("编辑"),
		(&actions.Delete{}).Init("删除"),
	}
}

// 执行行为后回调
func (p *Scheduler) AfterAction(ctx *builder.Context, uriKey string, query *gorm.DB) interface{} {
	if uriKey == "delete" {
		// 重载服务
		(&model.Scheduler{}).ReloadServices()
	}

	return nil
}

// 执行列表编辑完成后回调
func (p *Scheduler) AfterEditable(ctx *builder.Context, id interface{}, field string, value interface{}) interface{} {
	if field == "status" {
		// 重载服务
		(&model.Scheduler{}).ReloadServices()
	}

	return ctx.SimpleSuccess("操作成功")
}
