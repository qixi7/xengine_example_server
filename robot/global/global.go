package global

import (
	"github.com/qixi7/xengine_core/xcontainer/job"
	"github.com/qixi7/xengine_core/xcontainer/timer"
	"xserver/robot/root"
)

/*
	global.go: some global function
*/

func GetJobMgr() *job.Controller {
	return root.StaticRoot.Getter(root.ESJob).Get().(*job.Controller)
}

func GetTimerMgr() *timer.Controller {
	return root.GameRoot.Getter(root.EGTimer).Get().(*timer.Controller)
}
