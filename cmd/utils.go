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
	"archive/tar"
	"context"
	"fmt"
	"io"
	"kconsole/config"
	"kconsole/utils/bcs"
	"kconsole/utils/errorx"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/manifoldco/promptui"
	"github.com/pingcap/errors"
	"github.com/pterm/pterm"
	v1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
)

var (
	once      sync.Once
	clientSet *kubernetes.Clientset = &kubernetes.Clientset{}
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

// defaultKubeConfig used to configure the kubeclient by ~/.kube/config
func defaultKubeConfig() *rest.Config {
	home, err := os.UserHomeDir()
	errorx.CheckError(err)

	// 构建kubeconfig文件路径
	kubeconfig := filepath.Join(home, ".kube", "config")

	// 加载kubeconfig文件
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	errorx.CheckError(err)

	return config
}

func newKubeConfigForToken(host string, token string) *rest.Config {
	config := &rest.Config{
		Host:        host,
		BearerToken: token,
	}
	return config
}

func newBcsConfig(clusterid string) *rest.Config {
	kconsoleConfig := config.GetKconsoleConfig()
	return newKubeConfigForToken(
		fmt.Sprintf("%s/clusters/%s", kconsoleConfig.BCSHost, clusterid),
		kconsoleConfig.BCSToken,
	)
}

func bcsClientSet(clusterid string) *kubernetes.Clientset {
	config := newBcsConfig(clusterid)
	clientset, err := kubernetes.NewForConfig(config)
	errorx.CheckError(err)
	return clientset
}

func defaulClientSet() *kubernetes.Clientset {
	config := defaultKubeConfig()
	clientset, err := kubernetes.NewForConfig(config)
	errorx.CheckError(err)
	return clientset
}

func getClientSet() *kubernetes.Clientset {
	once.Do(func() {
		c := config.GetKconsoleConfig()
		switch c.Auth {
		case config.LocalConfigAuth:
			clientSet = defaulClientSet()
		case config.BcsAuth:
			// select cluster
			clusterid := selectBCSCluster()
			clientSet = bcsClientSet(clusterid)
		}
	})
	return clientSet
}

func selectBCSCluster() (clusterid string) {
	projs, err := bcs.UserBCSProjects(context.Background())
	errorx.CheckErrorWithCode(err, errorx.ErrorGetBCSUserProjErr)
	// build project list of []string
	projnames := make([]string, 0)
	// build name:id map for project
	nameid := make(map[string]string, 0)
	for _, proj := range projs.Data.Results {
		projnames = append(projnames, proj.Name)
		nameid[proj.Name] = proj.ProjectID
	}
	selectprojname := SelectUI(projnames, "select a bcs project")
	selectprojid := nameid[selectprojname]
	// ----
	clusters, err := bcs.UserBCSCluster(context.Background(), selectprojid)
	errorx.CheckErrorWithCode(err, errorx.ErrorGetBCSUserProjErr)
	// ----
	clusternames := make([]string, 0)
	clusternameid := make(map[string]string, 0)
	for _, cluster := range clusters.Data {
		clusternames = append(clusternames, cluster.ClusterName)
		clusternameid[cluster.ClusterName] = cluster.ClusterID
	}
	selectclustername := SelectUI(clusternames, "select a bcs cluster")
	selectclusterid := clusternameid[selectclustername]
	return selectclusterid
}

func allPodList() *v1.PodList {
	pods, err := getClientSet().CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	errorx.CheckError(err)

	return pods
}

func getPod(podname string, namespace string) (*v1.Pod, error) {
	pod, err := getClientSet().CoreV1().Pods(namespace).Get(context.Background(), podname, metav1.GetOptions{})
	errorx.CheckError(err)

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
	errorx.CheckError(err)

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
	errorx.CheckErrorWithCode(err, errorx.ErrorSelectExit)

	return result
}

func InputUI(title string, prefix string, defaultStr string) string {
	validate := func(input string) error {
		if prefix == "/" {
			if !strings.HasPrefix(input, "/") {
				return errors.New("please start with '/'")
			}
		} else {
			if !strings.HasPrefix(input, "/") && !strings.HasPrefix(input, "./") && !strings.HasPrefix(input, "../") {
				return errors.New("please start with '/' or './' or '../'")
			}
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    title,
		Validate: validate,
		Default:  defaultStr,
	}
	result, err := prompt.Run()
	if err != nil {
		panic(err.Error())
	}
	return result
}

// ---
// exec utils
// ---

func ExecPodContainer(namespace string, pod string, container string, command string) error {
	clientset := getClientSet()
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
	errorx.CheckError(err)

	err = executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		if k8serror.IsNotFound(err) {
			fmt.Printf("Pod %s/%s not found\n", namespace, pod)
		} else {
			fmt.Printf("Error executing command in container %s of pod %s/%s: %v\n", container, namespace, pod, err)
		}
		return err
	}

	return nil
}

func copyFromPod(namespace string, pod string, container string, srcPath string, destPath string) error {
	clientset := getClientSet()
	reader, outStream := io.Pipe()
	//todo some containers failed : tar: Refusing to write archive contents to terminal (missing -f option?) when execute `tar cf -` in container
	cmdArr := []string{"tar", "cf", "-", srcPath}
	req := clientset.CoreV1().RESTClient().
		Get().
		Namespace(namespace).
		Resource("pods").
		Name(pod).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: container,
			Command:   cmdArr,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(defaultKubeConfig(), http.MethodPost, req.URL())
	if err != nil {
		log.Fatalf("error %s\n", err)
		return err
	}
	go func() {
		defer outStream.Close()
		err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: outStream,
			Stderr: os.Stderr,
			Tty:    false,
		})
		cmdutil.CheckErr(err)
	}()
	prefix := getPrefix(srcPath)
	prefix = path.Clean(prefix)
	prefix = stripPathShortcuts(prefix)
	destPath = path.Join(destPath, path.Base(prefix))
	err = unTarAll(reader, destPath, prefix)
	return err
}

func copyToPod(namespace string, pod string, container string, srcPath string, destPath string) error {
	clientset := getClientSet()
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		cmdutil.CheckErr(makeTar(srcPath, destPath, writer))
	}()

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command:   []string{"tar", "-xmf", "-"},
			Container: container,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(defaultKubeConfig(), http.MethodPost, req.URL())
	if err != nil {
		return err
	}

	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  reader,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})
	if err != nil {
		return err
	}
	return nil
}

func getPrefix(file string) string {
	return strings.TrimLeft(file, "/")
}

// stripPathShortcuts removes any leading or trailing "../" from a given path
func stripPathShortcuts(p string) string {

	newPath := path.Clean(p)
	trimmed := strings.TrimPrefix(newPath, "../")

	for trimmed != newPath {
		newPath = trimmed
		trimmed = strings.TrimPrefix(newPath, "../")
	}

	// trim leftover {".", ".."}
	if newPath == "." || newPath == ".." {
		newPath = ""
	}

	if len(newPath) > 0 && string(newPath[0]) == "/" {
		return newPath[1:]
	}

	return newPath
}

func unTarAll(reader io.Reader, destDir, prefix string) error {
	tarReader := tar.NewReader(reader)
	for {
		fmt.Println("aa")
		header, err := tarReader.Next()
		fmt.Println("aa")
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		if !strings.HasPrefix(header.Name, prefix) {
			return fmt.Errorf("tar contents corrupted")
		}

		mode := header.FileInfo().Mode()
		destFileName := filepath.Join(destDir, header.Name[len(prefix):])
		baseName := filepath.Dir(destFileName)
		if err := os.MkdirAll(baseName, 0755); err != nil {
			return err
		}
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(destFileName, 0755); err != nil {
				return err
			}
			continue
		}

		evaledPath, err := filepath.EvalSymlinks(baseName)
		if err != nil {
			return err
		}

		if mode&os.ModeSymlink != 0 {
			linkname := header.Linkname

			if !filepath.IsAbs(linkname) {
				_ = filepath.Join(evaledPath, linkname)
			}

			if err := os.Symlink(linkname, destFileName); err != nil {
				return err
			}
		} else {
			outFile, err := os.Create(destFileName)
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			if err := outFile.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func makeTar(srcPath, destPath string, writer io.Writer) error {
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	srcPath = path.Clean(srcPath)
	destPath = path.Clean(destPath)
	return recursiveTar(path.Dir(srcPath), path.Base(srcPath), path.Dir(destPath), destPath, tarWriter)
}

func recursiveTar(srcBase, srcFile, destBase, destFile string, tw *tar.Writer) error {
	srcPath := path.Join(srcBase, srcFile)
	matchedPaths, err := filepath.Glob(srcPath)
	if err != nil {
		return err
	}
	for _, fpath := range matchedPaths {
		stat, err := os.Lstat(fpath)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			files, err := os.ReadDir(fpath)
			if err != nil {
				return err
			}
			if len(files) == 0 {
				//case empty directory
				hdr, _ := tar.FileInfoHeader(stat, fpath)
				hdr.Name = destFile
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
			}
			for _, f := range files {
				if err := recursiveTar(srcBase, path.Join(srcFile, f.Name()), destBase, path.Join(destFile, f.Name()), tw); err != nil {
					return err
				}
			}
			return nil
		} else if stat.Mode()&os.ModeSymlink != 0 {
			//case soft link
			hdr, _ := tar.FileInfoHeader(stat, fpath)
			target, err := os.Readlink(fpath)
			if err != nil {
				return err
			}

			hdr.Linkname = target
			hdr.Name = destFile
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
		} else {
			//case regular file or other file type like pipe
			hdr, err := tar.FileInfoHeader(stat, fpath)
			if err != nil {
				return err
			}
			hdr.Name = destFile

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			f, err := os.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
			return f.Close()
		}
	}
	return nil
}
