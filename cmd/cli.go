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
	"os"

	"github.com/spf13/cobra"
)

var (
	uiSize    int
	macNotify bool
)

type Cli struct {
	rootCmd *cobra.Command
}

func NewCli() *Cli {
	cli := &Cli{
		rootCmd: &cobra.Command{
			Use:   "kconsole",
			Short: "container terminal manager.",
			Long:  PrintLogo(),
		},
	}
	cli.rootCmd.SetOut(os.Stdout)
	cli.rootCmd.SetErr(os.Stderr)
	cli.setFlags()
	cli.rootCmd.DisableAutoGenTag = true
	return cli
}

func (cli *Cli) setFlags() {
	flags := cli.rootCmd.PersistentFlags()
	flags.IntVar(&uiSize, "ui-size", 4, "number of list items to show in menu at once")
	flags.BoolVarP(&macNotify, "mac-notify", "m", false, "enable to display Mac notification banner")
}

// Run command
func (cli *Cli) Run() error {
	return cli.rootCmd.Execute()
}
