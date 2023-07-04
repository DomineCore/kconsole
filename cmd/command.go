// MIT License
//
// # Copyright (c) 2023 Core
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

import "github.com/spf13/cobra"

type Command interface {
	Init()
	CobraCmd() *cobra.Command
}

type BaseCommand struct {
	command *cobra.Command
}

func (bc BaseCommand) Init() {
}

func (bc BaseCommand) CobraCmd() *cobra.Command {
	return bc.command
}

func (bc *BaseCommand) AddCommands(children ...Command) {
	for _, child := range children {
		child.Init()
		childCmd := child.CobraCmd()
		bc.CobraCmd().AddCommand(childCmd)
	}
}

func NewBaseCommand() *BaseCommand {
	cli := NewCli()
	baseCmd := &BaseCommand{
		command: cli.rootCmd,
	}
	baseCmd.AddCommands(&ConsoleCmd{})
	baseCmd.AddCommands(&DownloadCmd{})
	baseCmd.AddCommands(&UploadCmd{})
	baseCmd.AddCommands(&LogCmd{})
	return baseCmd
}
