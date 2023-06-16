# kconsole

kconsole 是一个用于与提高 Kubernetes 集群容器交互效率的命令行工具。它使用client-go 来与 Kubernetes 集群交互，提供快速进入容器终端、下载容器文件、上传本地文件到容器等便捷功能。

## 使用
以 console 子命令为例，要使用 kconsole 进入 Kubernetes 集群中的容器终端，请按照以下步骤操作：

### 1 下载 kconsole 可执行文件。
您可以从 GitHub Releases 页面 下载最新版本的 kconsole。

### 2 将 kconsole 可执行文件添加到您的 PATH 环境变量中。
例如，如果您将 kconsole 可执行文件下载到 /usr/local/bin 目录中，则可以使用以下命令将其添加到 PATH 环境变量中：

```
export PATH=$PATH:/usr/local/bin
```
现在，您可以在命令行中使用 kconsole 命令了。

### 3 进入容器console

在命令行中运行以下命令：

```
kconsole console
```
这将显示一个交互式菜单，列出了 Kubernetes 集群中的所有 Pod以及它们所在的namespace。

选择要进入的 Pod。您可以使用上下箭头键来选择 Pod，也可以输入名称进行搜索，然后按 Enter 键确认选择。

选择要进入的容器。您可以使用上下箭头键来选择容器，然后按 Enter 键确认选择。

现在，您已经进入了容器终端。您可以在终端中运行任何命令，就像在本地终端中一样。

要退出容器终端，请输入 exit 命令。

## 选项
kconsole 支持以下选项：

-h, --help: 显示帮助信息。

## 子命令
kconsole 提供以下子命令:

console: 进入集群中的容器终端
download: 下载集群中的容器内文件
upload: 上传本地文件到集群中的容器

## 开发
如果您想要为 kconsole 做出贡献，或者想要构建自己的版本，请按照以下步骤操作：

克隆 kconsole 仓库：

```
git clone https://github.com/DomineCore/kconsole.git
```
进入 kconsole 目录：

```
cd kconsole
```
构建 kconsole 可执行文件：
```
go build .
```
这将在 bin 目录中生成一个名为 kconsole 的可执行文件。

运行 kconsole：

```
./kconsole
```
现在，您可以测试您的更改是否正常工作了。

## 许可证
kconsole 使用 MIT 许可证。有关更多信息，请参见 LICENSE 文件。
