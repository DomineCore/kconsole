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

type UploadCmd struct {
	BaseCommand
}

func (cl *UploadCmd) Init() {
	cl.command = &cobra.Command{
		Use:   "upload",
		Short: "Copy files locally to remote",
		Long:  "Copy files locally to remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.runCluster(cmd, args)
		},
	}
	cl.command.DisableFlagsInUseLine = true
}

func (cl UploadCmd) runCluster(cmd *cobra.Command, args []string) error {
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
	// input src file
	inputsourcecmd := InputUI("input local source file path", "local", "")
	// input dest file
	inputdestcmd := InputUI("input container dest file path", "/", "")
	// build exec real command
	err := copyToPod(namespace, podname, selectcontainer, inputsourcecmd, inputdestcmd)
	return err
}
