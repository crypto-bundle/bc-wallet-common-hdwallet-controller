package config

import "time"

func (c *MangerConfig) GetHostName() string {
	return c.baseAppCfgSrv.GetHostName()
}

func (c *MangerConfig) GetEnvironmentName() string {
	return c.baseAppCfgSrv.GetEnvironmentName()
}

func (c *MangerConfig) IsProd() bool {
	return c.baseAppCfgSrv.IsProd()
}

func (c *MangerConfig) IsStage() bool {
	return c.baseAppCfgSrv.IsStage()
}

func (c *MangerConfig) IsTest() bool {
	return c.baseAppCfgSrv.IsTest()
}

func (c *MangerConfig) IsDev() bool {
	return c.baseAppCfgSrv.IsDev()
}

func (c *MangerConfig) IsDebug() bool {
	return c.baseAppCfgSrv.IsDebug()
}

func (c *MangerConfig) IsLocal() bool {
	return c.baseAppCfgSrv.IsLocal()
}

func (c *MangerConfig) GetStageName() string {
	return c.baseAppCfgSrv.GetStageName()
}

func (c *MangerConfig) GetApplicationPID() int {
	return c.baseAppCfgSrv.GetApplicationPID()
}

func (c *MangerConfig) GetApplicationName() string {
	return c.baseAppCfgSrv.GetApplicationName()
}

func (c *MangerConfig) SetApplicationName(appName string) {
	c.baseAppCfgSrv.SetApplicationName(appName)
}

func (c *MangerConfig) GetReleaseTag() string {
	return c.baseAppCfgSrv.GetReleaseTag()
}

func (c *MangerConfig) GetCommitID() string {
	return c.baseAppCfgSrv.GetCommitID()
}

func (c *MangerConfig) GetShortCommitID() string {
	return c.baseAppCfgSrv.GetShortCommitID()
}

func (c *MangerConfig) GetBuildNumber() uint64 {
	return c.baseAppCfgSrv.GetBuildNumber()
}

func (c *MangerConfig) GetBuildDateTS() int64 {
	return c.baseAppCfgSrv.GetBuildDateTS()
}

func (c *MangerConfig) GetBuildDate() time.Time {
	return c.baseAppCfgSrv.GetBuildDate()
}
