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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsoleCmd_Init(t *testing.T) {
	// Create a mock BaseCommand
	baseCmd := &BaseCommand{}

	// Create a ConsoleCmd with the mock BaseCommand
	ConsoleCmd := &ConsoleCmd{
		BaseCommand: *baseCmd,
	}

	// Call Init on the ConsoleCmd
	ConsoleCmd.Init()

	// Check that the ConsoleCmd's command has the expected Use, Short, and Long fields
	assert.Equal(t, "console", ConsoleCmd.command.Use)
	assert.Equal(t, "Exec a command for a container incluster.", ConsoleCmd.command.Short)
	assert.Equal(t, "Exec a command for a container incluster.", ConsoleCmd.command.Long)
	assert.True(t, ConsoleCmd.command.DisableFlagsInUseLine)
	assert.NotNil(t, ConsoleCmd.command.RunE)
}
