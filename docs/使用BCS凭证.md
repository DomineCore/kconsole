# 使用BCS凭证

kconsole支持使用bcs(蓝鲸容器管理平台)的个人密钥进行认证，自动获取到bcs项目和集群，无需配置~/.kconsole/config即可使用kconsole的所有能力。

要使用bcs个人密钥认证模式，只需按照如下的配置：

## 修改配置文件~/.kconsole/config.yaml
```
auth: bcs
bcsHost: "your bcs host root path"
bcsToke: "your bcs token"
```

