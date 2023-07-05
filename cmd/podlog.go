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

import (
	"kconsole/utils/errorx"

	"github.com/spf13/cobra"
)

type LogCmd struct {
	BaseCommand
}

func (cl *LogCmd) Init() {
	cl.command = &cobra.Command{
		Use:   "log",
		Short: "show pod's log for a container incluster.",
		Long:  "show pod's log for a container incluster. Only the latest 150 lines.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.runConsole(cmd, args)
		},
	}
	cl.command.DisableFlagsInUseLine = true
}

func (cl LogCmd) runConsole(cmd *cobra.Command, args []string) error {
	// call utils get pods
	podname, namespace, selectcontainer := SelectContainer()
	// build exec real command
	lines, err := cmd.Flags().GetInt64(flagLines)
	errorx.CheckError(err)
	err = PrintLogs(namespace, podname, selectcontainer, lines)
	return err
}
