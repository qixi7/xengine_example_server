package basemodule

import "fmt"

/*
	netmodule.go: 机器人网络模块
*/

// 基础信息模块
type BaseInfoModule struct {
	roleID uint32 // roleID
	name   string // name
}

// --------------------------- 内置函数 ---------------------------

// --------------------------- 外置函数 ---------------------------

// init
func (c *BaseInfoModule) Init(roleID uint32, name string) {
	c.roleID = roleID
	c.name = name
}

// 基础信息
func (c *BaseInfoModule) BaseString() string {
	return fmt.Sprintf("roleID=%d, name=%s", c.roleID, c.name)
}

// set
func (c *BaseInfoModule) SetBaseInfo(roleID uint32, name string) {
	c.roleID = roleID
	c.name = name
}

// get
func (c *BaseInfoModule) GetRoleID() uint32 {
	return c.roleID
}

// get
func (c *BaseInfoModule) GetName() string {
	return c.name
}
