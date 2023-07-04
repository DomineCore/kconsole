// MIT License
//
// Copyright (c) 2023 Core
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCli_Run(t *testing.T) {
	// 创建Cli对象
	cli := NewCli()

	// 设置命令行参数
	args := []string{"--lines", "10"}

	// 将命令行参数设置为args
	cli.rootCmd.SetArgs(args)

	// 创建一个bytes.Buffer对象，用于捕获输出
	buf := new(bytes.Buffer)

	// 将输出设置为buf
	cli.rootCmd.SetOut(buf)

	// 执行命令
	err := cli.Run()

	// 断言命令执行没有错误
	assert.NoError(t, err)

}

func TestCli_setFlags(t *testing.T) {
	// 创建Cli对象
	cli := NewCli()

	// 设置命令行参数
	args := []string{"--lines", "10"}

	// 将命令行参数设置为args
	cli.rootCmd.SetArgs(args)

	// 解析命令行参数
	err := cli.rootCmd.ParseFlags(args)

	// 断言命令解析没有错误
	assert.NoError(t, err)

}

func TestNewCli(t *testing.T) {
	// 创建Cli对象
	cli := NewCli()

	// 断言Cli对象不为空
	assert.NotNil(t, cli)

	// 断言rootCmd对象不为空
	assert.NotNil(t, cli.rootCmd)

	// 断言rootCmd对象的Use属性被正确设置
	assert.Equal(t, "kconsole", cli.rootCmd.Use)

	// 断言rootCmd对象的Short属性被正确设置
	assert.Equal(t, "container terminal manager.", cli.rootCmd.Short)
}

func TestMain(m *testing.M) {
	// 在测试之前设置
	os.Exit(m.Run())
}
