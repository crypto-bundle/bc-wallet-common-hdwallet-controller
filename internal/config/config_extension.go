/*
 *
 *
 * MIT-License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

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
