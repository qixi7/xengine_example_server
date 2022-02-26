package global

import (
	"robot/root"
	"xcore/xcontainer/job"
	"xcore/xcontainer/timer"
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
