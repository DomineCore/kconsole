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
	"fmt"
	"kconsole/utils/errorx"

	"github.com/spf13/cobra"
)

type LogDownCmd struct {
	BaseCommand
}

func (cl *LogDownCmd) Init() {
	cl.command = &cobra.Command{
		Use:   "logdown",
		Short: "download pod's log for a container incluster.",
		Long:  "download pod's log for a container incluster. Only the latest 150 lines.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.runLogDown(cmd, args)
		},
	}
	cl.command.DisableFlagsInUseLine = true
}

func (cl LogDownCmd) validateArgs(args []string) (downFilename string) {
	if len(args) < 1 {
		errorx.CheckErrorWithCode(fmt.Errorf("provide at least one file name for storing logs."), errorx.ErrorArgsErr)
	}
	downFilename = args[0]
	return
}

func (cl LogDownCmd) runLogDown(cmd *cobra.Command, args []string) error {
	// validate args logfilename
	// call utils get pods
	downFilename := cl.validateArgs(args)
	podname, namespace, selectcontainer := SelectContainer()
	// build exec real command
	err := SaveLogs(namespace, podname, selectcontainer, downFilename)
	return err
}
