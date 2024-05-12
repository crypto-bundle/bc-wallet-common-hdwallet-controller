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
	"errors"
	"net"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrUnableToFindActiveFileSocket = errors.New("unable to find active file socket")
	ErrMissingDirEntry              = errors.New("missing dir entry, probably socket files not found")
)

type socketDialler struct {
	dirName     string
	filePattern string

	dirEntries           []os.DirEntry
	count                uint
	currentEntryPosition uint
}

func (d *socketDialler) next() (os.DirEntry, bool) {
	position := d.currentEntryPosition
	d.currentEntryPosition++

	if d.currentEntryPosition <= d.count {
		err := d.prepare()
		if err != nil {
			return nil, false
		}

		return d.dirEntries[position], false
	}

	return d.dirEntries[position], true
}

func (d *socketDialler) Prepare() (func(context.Context, string) (net.Conn, error), error) {
	err := d.prepare()
	if err != nil {
		return nil, err
	}

	return d.DialCallback, nil
}

func (d *socketDialler) prepare() error {
	for {
		count, err := d.reset()
		if err == nil && count > 0 {
			return nil
		}

		time.Sleep(time.Second)
	}
}

func (d *socketDialler) reset() (uint, error) {
	d.dirEntries = make([]os.DirEntry, 0)
	d.count = 0
	d.currentEntryPosition = 0

	files, err := os.ReadDir(d.dirName)
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		match, loopErr := filepath.Match(d.filePattern, file.Name())
		if loopErr != nil {
			return 0, loopErr
		}

		if match {
			d.dirEntries = append(d.dirEntries, file)
			d.count++
		}
	}

	return d.count, nil
}

func (d *socketDialler) DialCallback(ctx context.Context, _ string) (net.Conn, error) {
	file, hasNext := d.next()
	if file == nil && !hasNext {
		return nil, ErrMissingDirEntry
	}

	filePath := filepath.Join(d.dirName, file.Name())

	resolved, err := net.ResolveUnixAddr("unix", filePath)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUnix("unix", nil, resolved)
	if err != nil {
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
