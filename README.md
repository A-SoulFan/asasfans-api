# asasfans-goapi

asasfans api支持

项目主页

- [国内访问](https://app.asf.ink/)
- [全球访问](https://app.a-soul.fans)

## 环境

- golang 1.17+
- mysql 8.0+

## 目录结构

```shell
├─cmd # 可执行文件
│  └─asasapi
├─config   # 配置文件
└─internal # 内部使用
   ├─app   # 具体应用
   │  └─asasapi 
   │      ├─handler # http handler
   │      ├─help    # 助手函数
   │      └─router  # http router
   ├─launcher  # 基础启动器
   └─pkg   # 第三方依赖封装
      ├─database
      ├─httpserver
      └─log
```

## 开始

### 本地运行

```shell
# clone repository
git clone git@github.com:A-SoulFan/asasfans-api.git

# download go mod 
cd asasfans-api && go mod download
```

### Docker builder

```shell
# clone repository
git clone git@github.com:A-SoulFan/asasfans-api.git

cd asasfans-api

# docker builder
# 如果在 CN 进行 build 请自行将 Dockerfile 中注释的镜像源相关行开启
docker build --rm -t asasfans-api:latest -f builder/asasapi/Dockerfile .

# copy config file 并自行修改相关配置
cp config/config.template.yml config/asasapi.yml

# docker run
docker run \
  --detach \
  --name asasfans-api \
  --volume $PWD/config:/config \
  asasfans-api:latest
```

## 开发者规范

### 命名规范

- 文件命名
  - 全小写
  - 尽可能短
  - 尽量避免使用 `_`，如果一定要分割，使用 `_`
- 变量命名
  - 驼峰
- 常量命名
  - 驼峰

### 开发流程

[Fork & Pull Request 流程](https://aaronflower.github.io/essays/github-fork-pull-workflow.html)

- fork repository
- checkout develop -> feature
- coding
- pull request
- review
- merge

### 相关工具

- [apifox](https://www.apifox.cn/)
  - 注册后联系管理者加组
