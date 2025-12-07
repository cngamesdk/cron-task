package sql_cleaning

import (
	"cngamesdk.com/cron-task/global"
	"cngamesdk.com/cron-task/internal/service"
	"cngamesdk.com/cron-task/model/sql/cron_task"
	"context"
	"github.com/cngamesdk/go-core/model/sql"
	"github.com/duke-git/lancet/v2/datetime"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"strings"
	"time"
)

const (
	startDateTimeKey = "StartDateTime"
	endDateTimeKey   = "EndDateTime"
	startDateKey     = "StartDate"
	endDateKey       = "EndDate"
)

// NewSqlCleaningService 实例化服务
func NewSqlCleaningService(req *cron_task.DimCronTaskConfigModel) *SqlCleaningService {
	myService := &SqlCleaningService{}
	myService.Config = req
	return myService
}

type SqlCleaningService struct {
	service.BaseService
}

func (receiver *SqlCleaningService) presetVariableStartDateTime(req string) string {
	findStartDateTime := "{{" + startDateTimeKey + "}}"
	startDateTime, ok := receiver.Config.Config[startDateTimeKey]
	startDateTimeStr := datetime.FormatTimeToStr(time.Now().Add(time.Hour*(-12)), "yyyy-mm-dd hh:mm:ss")
	if ok {
		startDateTimeStr = startDateTime.(string)
	}
	req = strings.ReplaceAll(req, findStartDateTime, startDateTimeStr)
	return req
}

func (receiver *SqlCleaningService) presetVariableEndDateTime(req string) string {
	findEndDateTime := "{{" + endDateTimeKey + "}}"
	if strings.Index(req, findEndDateTime) >= 0 {
		endDateTime, ok := receiver.Config.Config[endDateTimeKey]
		endDateTimeStr := datetime.FormatTimeToStr(time.Now(), "yyyy-mm-dd hh:mm:ss")
		if ok {
			endDateTimeStr = endDateTime.(string)
		}
		req = strings.ReplaceAll(req, findEndDateTime, endDateTimeStr)
	}
	return req
}

func (receiver *SqlCleaningService) presetVariableStartDate(req string) string {
	findStartDateTime := "{{" + startDateKey + "}}"
	startDateTime, ok := receiver.Config.Config[startDateKey]
	startDateTimeStr := datetime.FormatTimeToStr(time.Now().Add(time.Hour*(-12)), "yyyy-mm-dd")
	if ok {
		startDateTimeStr = startDateTime.(string)
	}
	req = strings.ReplaceAll(req, findStartDateTime, startDateTimeStr)
	return req
}

func (receiver *SqlCleaningService) presetVariableEndDate(req string) string {
	findEndDateTime := "{{" + endDateKey + "}}"
	if strings.Index(req, findEndDateTime) >= 0 {
		endDateTime, ok := receiver.Config.Config[endDateKey]
		endDateTimeDeal, _ := datetime.FormatStrToTime(datetime.FormatTimeToStr(time.Now().Add(time.Hour*24), "yyyy-mm-dd"), "yyyy-mm-dd")
		endDateTimeStr := datetime.FormatTimeToStr(endDateTimeDeal.Add(time.Second*-1), "yyyy-mm-dd hh:mm:ss")
		if ok {
			endDateTimeStr = endDateTime.(string)
		}
		req = strings.ReplaceAll(req, findEndDateTime, endDateTimeStr)
	}
	return req
}

func (receiver *SqlCleaningService) presetVariables() string {
	var presetVariables []func(string2 string) string
	presetVariables = append(presetVariables, receiver.presetVariableStartDateTime)
	presetVariables = append(presetVariables, receiver.presetVariableEndDateTime)
	presetVariables = append(presetVariables, receiver.presetVariableStartDate)
	presetVariables = append(presetVariables, receiver.presetVariableEndDate)
	content := receiver.Config.Content
	for _, item := range presetVariables {
		content = item(content)
	}
	return content
}

func (receiver *SqlCleaningService) PreEvent(ctx context.Context) (resp string, err error) {
	resp = receiver.presetVariables()
	return
}

func (receiver *SqlCleaningService) Run(ctx context.Context) (err error) {
	startTime := time.Now()
	var execSql string
	defer func() {
		endTime := time.Now()
		if err != nil {
			if failErr := receiver.FailEvent(ctx); failErr != nil {
				err = failErr
				global.Logger.ErrorCtx(ctx, "失败事件异常", zap.Error(err))
			}
		} else {
			if successErr := receiver.SuccessEvent(ctx); successErr != nil {
				err = successErr
				global.Logger.ErrorCtx(ctx, "成功事件异常", zap.Error(err))
			}
		}
		logModel := cron_task.NewOdsCronTaskLogModel()
		logModel.ConfigId = receiver.Config.Id
		logModel.StartTime = startTime
		logModel.EndTime = endTime
		logModel.Latency = cast.ToInt(endTime.Sub(startTime).Milliseconds() / 1000)
		logModel.Status = sql.StatusSuccess
		if err != nil {
			logModel.Status = sql.StatusFail
		}
		logModel.Result = execSql
		receiver.BaseService.TaskLog = logModel
		if completeErr := receiver.CompleteEvent(ctx); completeErr != nil {
			err = completeErr
			global.Logger.ErrorCtx(ctx, "完成事件异常", zap.Error(err))
		}
	}()

	preEventResult, preEventErr := receiver.PreEvent(ctx)
	if preEventErr != nil {
		global.Logger.ErrorCtx(ctx, "前置事件异常", zap.Error(err))
		return
	}
	execSql = preEventResult
	err = global.MyDb.WithContext(ctx).Exec(execSql).Error
	return
}

func (receiver *SqlCleaningService) SuccessEvent(ctx context.Context) (err error) {
	successEndTimeStr := datetime.FormatTimeToStr(time.Now(), "yyyy-mm-dd hh:mm:ss")
	receiver.Config.Config[startDateTimeKey] = successEndTimeStr

	//记录表,保留最后执行时间
	model := cron_task.NewDimCronTaskConfigModel()
	if model.Config == nil {
		model.Config = make(sql.CustomMapType)
	}
	model.Config[startDateTimeKey] = successEndTimeStr
	if saveErr := model.Updates(ctx, "id = ?", receiver.Config.Id); saveErr != nil {
		err = saveErr
		global.Logger.ErrorCtx(ctx, "保存异常", zap.Error(saveErr))
		return
	}
	return
}

func (receiver *SqlCleaningService) FailEvent(ctx context.Context) (err error) {
	return
}

func (receiver *SqlCleaningService) CompleteEvent(ctx context.Context) (err error) {
	if err = receiver.BaseService.CompleteEvent(ctx); err != nil {
		global.Logger.ErrorCtx(ctx, "保存异常", zap.Error(err), zap.Any("data", receiver.TaskLog))
	}
	return
}
