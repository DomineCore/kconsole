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
	"kconsole/config"
	"kconsole/utils/errorx"

	"github.com/spf13/cobra"
)

type SwitchCmd struct {
	BaseCommand
}

func (cl *SwitchCmd) Init() {
	cl.command = &cobra.Command{
		Use:   "switch",
		Short: "Select a cluster.",
		Long:  "Select a cluster.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.runSwitch(cmd, args)
		},
	}
	cl.command.DisableFlagsInUseLine = true
}

func (cl *SwitchCmd) validate(args []string) {
	if config.GetKconsoleConfig().Auth != config.BcsAuth {
		errorx.CheckErrorWithCode(fmt.Errorf("the switch command must be used when auth=bcsauth. Run the login command first, for example, 'kconsole login --mode bcs --host xxx --token xxx'."), errorx.ErrorArgsErr)
	}
}

func (cl *SwitchCmd) runSwitch(cmd *cobra.Command, args []string) error {
	cl.validate(args)
	switch config.GetKconsoleConfig().Auth {
	case config.BcsAuth:
		clusterid := selectBCSCluster()
		config.UpdateConfilefile(map[string]string{"bcscluster": clusterid})
		fmt.Println("checkout cluster: ", clusterid, "~")
	}
	return nil
}
