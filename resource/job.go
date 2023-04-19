package resource

import (
	"strings"

	"github.com/quarkcms/quark-cron/action"
	"github.com/quarkcms/quark-cron/model"
	"github.com/quarkcms/quark-cron/search"
	"github.com/quarkcms/quark-go/pkg/app/handler/admin/actions"
	"github.com/quarkcms/quark-go/pkg/app/handler/admin/searches"
	"github.com/quarkcms/quark-go/pkg/builder"
	"github.com/quarkcms/quark-go/pkg/builder/template/adminresource"
	"github.com/quarkcms/quark-go/pkg/component/admin/form/fields/selectfield"
	"github.com/quarkcms/quark-go/pkg/component/admin/form/rule"
	"gorm.io/gorm"
)

type Job struct {
	adminresource.Template
}

// 初始化
func (p *Job) Init() interface{} {

	// 初始化模板
	p.TemplateInit()

	// 标题
	p.Title = "作业"

	// 模型
	p.Model = &model.Job{}

	// 分页
	p.PerPage = 10

	return p
}

func (p *Job) Fields(ctx *builder.Context) []interface{} {
	field := &adminresource.Field{}
	options := (&model.Scheduler{}).Options()

	return []interface{}{
		field.ID("id", "ID"),

		field.Text("name", "名称").
			SetRules([]*rule.Rule{
				rule.Required(true, "名称必须填写"),
			}),

		field.Select("schedule_id", "调度器").
			SetOptions(options).
			SetRules([]*rule.Rule{
				rule.Required(true, "请选择调度器"),
			}),

		field.Select("type", "类型").
			SetOptions([]*selectfield.Option{
				{
					Label: "Shell脚本",
					Value: 1,
				},
				{
					Label: "访问URL",
					Value: 2,
				},
				{
					Label: "执行应用",
					Value: 3,
				},
			}).
			SetRules([]*rule.Rule{
				rule.Required(true, "请选择类型"),
			}).
			SetDefault(1),

		field.Space("执行周期", []interface{}{
			field.Select("cycle_type").
				SetOptions([]*selectfield.Option{
					{
						Label: "每天",
						Value: 1,
					},
					{
						Label: "N天",
						Value: 2,
					},
					{
						Label: "每小时",
						Value: 3,
					},
					{
						Label: "N小时",
						Value: 4,
					},
					{
						Label: "N分钟",
						Value: 5,
					},
					{
						Label: "N秒",
						Value: 6,
					},
					{
						Label: "每星期",
						Value: 7,
					},
					{
						Label: "每月",
						Value: 8,
					},
				}).
				SetWhen(1, func() interface{} {
					return []interface{}{

						field.Number("hour").
							SetAddonAfter("点").
							SetWidth(70).
							SetMin(0).
							SetMax(23).
							SetDefault(1).
							OnlyOnForms(),

						field.Number("minute").
							SetAddonAfter("分").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(30).
							OnlyOnForms(),

						field.Number("second").
							SetAddonAfter("秒").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(30).
							OnlyOnForms(),
					}
				}).
				SetWhen(2, func() interface{} {
					return []interface{}{

						field.Number("day").
							SetAddonAfter("天").
							SetWidth(70).
							SetMin(1).
							SetMax(31).
							SetDefault(3).
							OnlyOnForms(),

						field.Number("hour").
							SetAddonAfter("小时").
							SetWidth(70).
							SetMin(0).
							SetMax(23).
							SetDefault(1).
							OnlyOnForms(),

						field.Number("minute").
							SetAddonAfter("分").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(30).
							OnlyOnForms(),

						field.Number("second").
							SetAddonAfter("秒").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(30).
							OnlyOnForms(),
					}
				}).
				SetWhen(3, func() interface{} {
					return []interface{}{
						field.Number("minute").
							SetAddonAfter("分").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(30).
							OnlyOnForms(),

						field.Number("second").
							SetAddonAfter("秒").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(30).
							OnlyOnForms(),
					}
				}).
				SetWhen(4, func() interface{} {
					return []interface{}{

						field.Number("hour").
							SetAddonAfter("小时").
							SetWidth(70).
							SetMin(0).
							SetMax(23).
							SetDefault(1).
							OnlyOnForms(),

						field.Number("minute").
							SetAddonAfter("分").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(59).
							OnlyOnForms(),

						field.Number("second").
							SetAddonAfter("秒").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(30).
							OnlyOnForms(),
					}
				}).
				SetWhen(5, func() interface{} {
					return []interface{}{
						field.Number("minute").
							SetAddonAfter("分").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(59).
							OnlyOnForms(),

						field.Number("second").
							SetAddonAfter("秒").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(30).
							OnlyOnForms(),
					}
				}).
				SetWhen(6, func() interface{} {
					return []interface{}{
						field.Number("second").
							SetAddonAfter("秒").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(30).
							OnlyOnForms(),
					}
				}).
				SetWhen(7, func() interface{} {
					return []interface{}{
						field.Select("week").
							SetOptions([]*selectfield.Option{
								{
									Label: "周一",
									Value: 1,
								},
								{
									Label: "周二",
									Value: 2,
								},
								{
									Label: "周三",
									Value: 3,
								},
								{
									Label: "周四",
									Value: 4,
								},
								{
									Label: "周五",
									Value: 5,
								},
								{
									Label: "周六",
									Value: 6,
								},
								{
									Label: "周日",
									Value: 7,
								},
							}).
							SetWidth(70).
							SetDefault(1).
							OnlyOnForms(),

						field.Number("hour").
							SetAddonAfter("点").
							SetWidth(70).
							SetMin(0).
							SetMax(23).
							SetDefault(1).
							OnlyOnForms(),

						field.Number("minute").
							SetAddonAfter("分").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(59).
							OnlyOnForms(),

						field.Number("second").
							SetAddonAfter("秒").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(30).
							OnlyOnForms(),
					}
				}).
				SetWhen(8, func() interface{} {
					return []interface{}{

						field.Number("day").
							SetAddonAfter("日").
							SetWidth(70).
							SetMin(1).
							SetMax(31).
							SetDefault(3).
							OnlyOnForms(),

						field.Number("hour").
							SetAddonAfter("点").
							SetWidth(70).
							SetMin(0).
							SetMax(23).
							SetDefault(1).
							OnlyOnForms(),

						field.Number("minute").
							SetAddonAfter("分").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(59).
							OnlyOnForms(),

						field.Number("second").
							SetAddonAfter("秒").
							SetWidth(70).
							SetMin(0).
							SetMax(59).
							SetDefault(30).
							OnlyOnForms(),
					}
				}).
				SetWidth(100).
				SetDefault(1).
				OnlyOnForms(),
		}),

		field.Dependency().
			SetWhen("type", 1, func() interface{} {
				return []interface{}{
					field.TextArea("shell", "脚本内容").
						SetRows(10).
						OnlyOnForms(),
				}
			}).
			SetWhen("type", 2, func() interface{} {
				return []interface{}{
					field.Text("url", "URL地址").
						SetWidth(350).
						OnlyOnForms(),
				}
			}).
			SetWhen("type", 3, func() interface{} {
				return []interface{}{
					field.Text("path", "应用路径").
						SetWidth(350).
						OnlyOnForms(),
				}
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
func (p *Job) Searches(ctx *builder.Context) []interface{} {

	return []interface{}{
		(&searches.Input{}).Init("name", "名称"),
		(&search.Status{}).Init(),
	}
}

// 行为
func (p *Job) Actions(ctx *builder.Context) []interface{} {

	return []interface{}{
		(&actions.Delete{}).Init("批量删除"),
		(&actions.Disable{}).Init("批量停止"),
		(&actions.Enable{}).Init("批量启动"),
		(&action.CreateModal{}).Init(p.Title),
		(&action.EditModal{}).Init("编辑"),
		(&actions.Delete{}).Init("删除"),
	}
}

// 数据保存后回调
func (p *Job) AfterSaved(ctx *builder.Context, id int, data map[string]interface{}, result *gorm.DB) interface{} {
	if result.Error != nil {
		return ctx.SimpleError(result.Error.Error())
	}

	// 重载服务
	(&model.Scheduler{}).ReloadServices()

	return ctx.SimpleSuccess("操作成功！", strings.Replace("/index?api="+adminresource.IndexRoute, ":resource", ctx.Param("resource"), -1))
}
