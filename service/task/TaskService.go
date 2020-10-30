package s_task

import (
	"apiproject/log"
	"go.uber.org/zap"
)

/**
@author 王世彪
	个人博客: https://sofineday.com?from=apiproject
	微信: 645102170
	QQ: 645102170
*/

var TaskService = &taskService{}

type taskService struct {
}

/**
执行任务1
*/
func (this *taskService) Task1(para interface{}) {
	log.Logger.Info("执行任务1, 完成", zap.Any("para", para))
}

/**
执行任务2
*/
func (this *taskService) Task2() {
	log.Logger.Info("执行任务2, 完成")
}
