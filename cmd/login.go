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
	"kconsole/utils/bcs"
	"kconsole/utils/errorx"

	"github.com/pingcap/errors"
	"github.com/spf13/cobra"
)

const (
	mode  = "mode"
	host  = "host"
	token = "token"
)

var (
	ModeLocal = "local"
	ModeBcs   = "bcs"
)

type LoginCmd struct {
	BaseCommand
}

func (cl *LoginCmd) Init() {
	cl.command = &cobra.Command{
		Use:   "login",
		Short: "Set cluster authentication information",
		Long:  "Set cluster authentication information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cl.runLogin(cmd, args)
		},
	}
	cl.command.DisableFlagsInUseLine = true
	cl.command.Flags().StringP(mode, "m", "local", "-mode=local will use the same cluster as the kubectl command will be used; -mode=bcs will use the cluster of bcs")
	cl.command.Flags().StringP(host, "H", "", "-host=xxx set host of the bcs")
	cl.command.Flags().StringP(token, "t", "", "-token=xxx set token of the bcs")
}

func (cl *LoginCmd) validate(args []string) (modeval, hostval, tokenval string) {
	modeval, err := cl.command.Flags().GetString(mode)
	errorx.CheckError(err)
	if modeval != ModeLocal && modeval != ModeBcs {
		errorx.CheckError(errors.New("invalid mode"))
	}
	if modeval == ModeBcs {
		hostval, err = cl.command.Flags().GetString(host)
		errorx.CheckError(err)
		if hostval == "" {
			errorx.CheckError(fmt.Errorf("host must not be empty when mode is %s", modeval))
		}
		tokenval, err = cl.command.Flags().GetString(token)
		errorx.CheckError(err)
		if tokenval == "" {
			errorx.CheckError(fmt.Errorf("token must not be empty when mode is %s", tokenval))
		}
	}
	return modeval, hostval, tokenval
}

func (cl *LoginCmd) runLogin(cmd *cobra.Command, args []string) error {
	mode, host, token := cl.validate(args)
	switch mode {
	case ModeLocal:
		// set config auth to local
		config.UpdateConfilefile(map[string]string{
			"auth": config.LocalConfigAuth,
		})
	case ModeBcs:
		// set config auth to bcs, add bcshost and bcstoken
		config.UpdateConfilefile(map[string]string{
			"auth":     config.BcsAuth,
			"bcshost":  host,
			"bcstoken": token,
		})
		// ping bcs host
		bcs.PingBCSOrdie(cmd.Context())
	}
	fmt.Println("login successfully.")
	return nil
}
