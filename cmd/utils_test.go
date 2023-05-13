// MIT License
//
// Copyright (c) 2023 Core
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

func TestPrintLogo(t *testing.T) {
	logo := PrintLogo()
	assert.Contains(t, logo, "Exec your container more easily.")
}

func TestDefaultKubeConfig(t *testing.T) {
	config := defaultKubeConfig()
	assert.NotNil(t, config)
}

func TestDefaulClientSet(t *testing.T) {
	clientset := defaulClientSet()
	assert.NotNil(t, clientset)
	discovery := clientset.Discovery()
	assert.NotNil(t, discovery)
	coreV1 := clientset.CoreV1()
	assert.NotNil(t, coreV1)
}

func TestAllPodList(t *testing.T) {
	pods := allPodList()
	assert.NotNil(t, pods)
	assert.NotEmpty(t, pods.Items)
}

func TestGetPod(t *testing.T) {
	pod, err := getPod("test-pod", "default")
	assert.Nil(t, pod)
	assert.NotNil(t, err)
}

func TestListAllPods(t *testing.T) {
	podNames := ListAllPods()
	assert.NotEmpty(t, podNames)
	assert.Contains(t, podNames[0], "/")
}
