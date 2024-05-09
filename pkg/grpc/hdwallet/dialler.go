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

package hdwallet

import (
	"context"
	"net"
	"os"
	"path/filepath"
)

type socketDialler struct {
	dirName     string
	filePattern string

	dirEntries           []os.DirEntry
	count                uint
	currentEntryPosition uint
}

func (d *socketDialler) next() (os.DirEntry, bool) {
	d.currentEntryPosition++

	if d.currentEntryPosition > d.count {
		return d.dirEntries[d.currentEntryPosition], false
	}

	return d.dirEntries[d.currentEntryPosition], true
}

func (d *socketDialler) Prepare() error {
	return d.prepare()
}

func (d *socketDialler) prepare() error {
	d.dirEntries = make([]os.DirEntry, 0)
	d.count = 0
	d.currentEntryPosition = 0

	files, err := os.ReadDir(d.dirName)
	if err != nil {
		return err
	}

	for _, file := range files {
		match, loopErr := filepath.Match(d.filePattern, file.Name())
		if loopErr != nil {
			return loopErr
		}

		if match {
			d.dirEntries = append(d.dirEntries, file)
			d.count++
		}
	}

	return nil
}

func (d *socketDialler) DialCallback(ctx context.Context, _ string) (net.Conn, error) {
	file, hasNext := d.next()
	filePath := filepath.Join(d.dirName, file.Name())

	resolved, err := net.ResolveUnixAddr("unix", filePath)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("unix", resolved.String())
	if err != nil {
		if !hasNext {
			return nil, d.prepare()
		}

		return nil, err
	}

	return conn, nil
}

func newSocketDialer(dirName, filePattern string) *socketDialler {
	return &socketDialler{
		dirName:     dirName,
		filePattern: filePattern,

		dirEntries:           nil,
		count:                0,
		currentEntryPosition: 0,
	}
}
