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
	"strings"

	"github.com/spf13/cobra"
)

var CMDS = []string{"sh", "bash"}

type ConsoleCmd struct {
	BaseCommand
}

func (cl *ConsoleCmd) Init() {
	cl.command = &cobra.Command{
		Use:   "console",
		Short: "Exec a command for a container incluster.",
		Long:  "Exec a command for a container incluster.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.runConsole(cmd, args)
		},
	}
	cl.command.DisableFlagsInUseLine = true
}

func (cl ConsoleCmd) runConsole(cmd *cobra.Command, args []string) error {
	// call utils get pods
	pods := ListAllPods()
	selectpod := SelectUI(pods, "select a pod")
	// pod: namespace/podname
	namespace_pod := strings.Split(selectpod, "/")
	namespace := namespace_pod[0]
	podname := namespace_pod[1]
	// call utils get container
	containers := ListContainersByPod(namespace, podname)
	selectcontainer := SelectUI(containers, "select a container")
	// select command
	selectcmd := SelectUI(CMDS, "select a cmd")
	// build exec real command
	err := ExecPodContainer(namespace, podname, selectcontainer, selectcmd)
	return err
}
