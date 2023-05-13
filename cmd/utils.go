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
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pingcap/errors"
	"github.com/pterm/pterm"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

// ----
// base utils
// ----
func PrintLogo() string {
	panel := pterm.DefaultHeader.WithMargin(8).
		WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).
		WithTextStyle(pterm.NewStyle(pterm.FgLightWhite)).Sprint("Exec your container more easily.")
	logo := pterm.FgLightGreen.Sprint(`
| | __  ___   ___   _ __   ___   ___  | |  ___ 
| |/ / / __| / _ \ | '_ \ / __| / _ \ | | / _ \
|   < | (__ | (_) || | | |\__ \| (_) || ||  __/
|_|\_\ \___| \___/ |_| |_||___/ \___/ |_| \___|
`)
	pterm.Info.Prefix = pterm.Prefix{
		Text:  "Tips",
		Style: pterm.NewStyle(pterm.BgBlue, pterm.FgLightWhite),
	}
	return fmt.Sprintf(`
%s%s
`, panel, logo)
}

// ----
// kube utils
// ----

func defaultKubeConfig() *rest.Config {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}

	// 构建kubeconfig文件路径
	kubeconfig := filepath.Join(home, ".kube", "config")

	// 加载kubeconfig文件
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return config
}

func defaulClientSet() *kubernetes.Clientset {
	config := defaultKubeConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func allPodList() *v1.PodList {
	pods, err := defaulClientSet().CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return pods
}

func getPod(podname string, namespace string) (*v1.Pod, error) {
	pod, err := defaulClientSet().CoreV1().Pods(namespace).Get(context.Background(), podname, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func ListAllPods() []string {
	pods := allPodList()
	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))
	}
	return podNames
}

func ListContainersByPod(namespace string, podname string) (containers []string) {
	pod, err := getPod(podname, namespace)
	if err != nil {
		panic(err.Error())
	}
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return
}

// ----
// ui utils
// ----

func SelectUI(data []string, title string) string {
	searcher := func(input string, index int) bool {
		item := data[index]
		loweritem := strings.Replace(strings.ToLower(item), " ", "", -1)
		return strings.Contains(loweritem, input)
	}

	prompt := promptui.Select{
		Label:    title,
		Items:    data,
		Searcher: searcher,
	}

	_, result, err := prompt.Run()
	if err != nil {
		panic(err.Error())
	}
	return result
}

// ---
// exec utils
// ---

func ExecPodContainer(namespace string, pod string, container string, command string) error {
	clientset := defaulClientSet()
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(pod).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: container,
			Command:   []string{command},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	// 创建执行器
	executor, err := remotecommand.NewSPDYExecutor(defaultKubeConfig(), http.MethodPost, req.URL())
	if err != nil {
		return err
	}
	err = executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s/%s not found\n", namespace, pod)
		} else {
			fmt.Printf("Error executing command in container %s of pod %s/%s: %v\n", container, namespace, pod, err)
		}
		return err
	}

	return nil
}
