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

package errorx

import (
	"os"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

const (
	// ErrorSelectExit The selector is closed
	ErrorSelectExit = 0
	// ErrorSelectExit
	ErrorConfigErr = 1
	// ErrorAuthConfigErr
	ErrorAuthConfigErr = 2
	// ErrorBCSAuthConfigErr bcs auth config error
	ErrorBCSAuthConfigErr = 3
	// ErrorGetBCSUserProjErr get bcs project unknown error
	ErrorGetBCSUserProjErr = 4
	// ErrorGetBCSUserProj get bcs cluster proj unknown error
	ErrorGetBCSUserClusterErr = 5
	// ErrorUnknow Unexpected error, need to contact the developer
	ErrorUnknow = 20
)

// CheckError if error is not nil, call fatal.
func CheckError(err error) {
	if err != nil {
		debug.PrintStack()
		Fatal(ErrorUnknow, err)
	}
}

func CheckErrorWithCode(err error, exitcode int) {
	if err != nil {
		Fatal(exitcode, err)
	}
}

func Fatal(exitcode int, args ...interface{}) {
	exitfunc := func() {
		os.Exit(exitcode)
	}
	log.RegisterExitHandler(exitfunc)
	log.Fatal(args...)
}
