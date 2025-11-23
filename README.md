# ChartProvider

[![ReleaseCard]][Release]![ReleaseDataCard]![LastCommitCard]  
![BuildStateCard]![ProjectLanguageCard]![ProjectLicense]

> ⚠️ **免责声明：本项目仅供学习和研究使用。
> Navigraph的相关服务和数据可能受到版权限制，请用户在使用时遵守当地法律法规及服务条款。
> 本项目不鼓励或支持任何侵犯知识产权的行为。**

ChartProvider 是一个基于 Go 语言开发的开源软件，用于学习请求代理的实现原理。  
为了学习本项目，你需要一个订阅的 Navigraph 账号。  
如果你还没有 Navigraph 账号，可以访问 [Navigraph官网](https://www.navigraph.com/) 了解更多信息。

## 如何使用

### ***(推荐)*** 使用Docker部署

1. ***(推荐)*** 使用docker-compose部署  
   i. 克隆或下载本项目到本地，并进入`docker`目录  
   ii. 按需编辑配置文件或`docker-compose.yml`文件  
   iii. 运行`docker-compose up -d`命令  
   iv. 访问[http://127.0.0.1:8080](http://127.0.0.1:8080)查看是否部署成功  
   v. 如果需要添加命令行参数
   ```yml
   services:
     fsd:
       image: halfnothing/chart-provider:latest
       # 省略部分字段
       command:
         - "-thread 32"
   ```
   推荐使用环境变量代替命令行参数
   ```yml
   services:
     fsd:
       image: halfnothing/chart-provider:latest
       # 省略部分字段
       environment:
         QUERY_THREAD: 32
   ```

2. 使用docker命令部署  
   命令示例如下
   ```shell
   docker run -d --name chart-provider -p 8080:8080 -p 8081:8081 -v $(pwd)/config.yaml:/chart-provider/config.yaml halfnothing/chart-provider:latest
   ``` 
   如果需要添加命令行参数, 则在命令的最后添加
   ```shell
   docker run -d ... halfnothing/chart-provider:latest -thread 32
   ```

3. 通过Dockerfile构建  
   i. 手动构建
   ```shell
   # 克隆本仓库
   git clone https://github.com/Flyleague-Collection/chart-provider.git
   # 进入项目目录
   cd chart-provider
   # 运行docker构建
   docker build -t chart-provider:latest .
   # 运行docker容器
   docker run -d --name chart-provider -p 8080:8080 -p 8081:8081 -v $(pwd)/config.yaml:/chart-provider/config.yaml chart-provider:latest
   ```
   ii. 自动构建
   ```shell
   # 克隆本仓库
   git clone https://github.com/Flyleague-Collection/chart-provider.git
   # 进入项目目录
   cd chart-provider
   # 进入docker目录并且修改docker-compose.yml文件
   cd docker
   vi docker-compose.yml
   ```
   将`image: halfnothing/chart-provider:latest`这一行替换为`build: ".."`    
   然后在同目录运行
   ```shell
   docker compose up -d
   ```

### 普通部署

1. 获取项目可执行文件
    - 前往[Release]页面下载最新版本
    - 前往[Action]页面下载最新开发版本
    - 手动[编译](#手动构建)本项目
2. [可选]下载[`config.yaml`](./docker/config.yaml)配置文件放置于可执行文件同级目录中
3. 运行可执行文件，如果配置文件存在，则使用配置文件，否则创建默认配置文件

## 首次运行

1. 运行项目
2. 查看项目日志  
   如果看到下面的日志，说明需要手动授权
    ```shell
    2025-11-24T01:26:07 | MAIN  | INFO  | TokenManager | Device authorization, please visit https://identity.api.navigraph.com/code/default.aspx?user_code=XXXXXXX to manual authorization
    ```
   请访问日志中提供的链接，并手动授权
3. 授权成功后，会看到下面的日志
    ```shell
    2025-11-24T01:26:17 | MAIN  | INFO  | TokenManager | Device authorization passed
   ```

一般来说只需要手动授权一次  
后续运行项目时, 会使用刷新令牌自动获取授权  
但如果长时间不使用接口, 刷新令牌可能会失效，此时需要重新手动授权    
如果自动授权成功，则会看到下面的日志

```shell
2025-11-24T01:27:49 | MAIN  | INFO  | TokenManager | Use cached flush token
```

如果自动授权失败，则会看到下面的日志

```shell
2025-11-24T01:47:28 | MAIN  | ERROR | TokenManager | refreshAccessToken StatusCode: 400
```

同时服务器也会自动重新申请手动授权

## 手动构建

```shell
# 克隆本仓库
git clone https://github.com/Flyleague-Collection/chart-provider.git
# 进入项目目录
cd chart-provider
# 确认安装了go编译器并且版本>=1.24.6
go version
# 运行go build命令
go build -ldflags="-w -s" .
# 对于windows系统, 可执行文件为chart-provider.exe
# 对于linux系统, 可执行文件为chart-provider
# [可选]使用upx压缩可执行文件
# windows
upx.exe -9 chart-provider.exe
# linux
upx -9 chart-provider
```

## 命令行参数与环境变量一览

| 命令行参数           | 环境变量             | 描述        | 默认值           |
|:----------------|:-----------------|:----------|:--------------|
| no_logs         | NO_LOGS          | 禁用日志输出到文件 | false         |
| config          | CONFIG_FILE_PATH | 配置文件路径    | "config.yaml" |
| signing_method  | SIGNING_METHOD   | 签名方法      | "HS512"       |
| request_timeout | REQUEST_TIMEOUT  | 请求超时时间    | "30s"         |
| gzip_level      | GZIP_LEVEL       | gzip压缩级别  | 5             |

## 贡献指南

1. 开一个 Issue 与我们讨论
2. Fork 本项目并完成你的修改
3. 不要修改任何除了你创建以外的源代码的版权信息
4. 遵守良好的代码编码规范
5. 开一个 Pull Request

## 开源协议

MIT License

Copyright © 2025 Half_nothing

无附加条款。

## 免责声明

1. 本软件仅供个人学习和研究目的使用，严禁用于任何商业用途。
2. 通过本软件获取或生成的任何内容，使用者应确保其使用方式符合相关法律法规，不得用于商业目的。
3. 使用者应尊重所有相关的知识产权，包括但不限于版权、商标权等，并严格遵守第三方服务的最终用户许可协议（EULA）及服务条款。
4. 本项目的开发者明确声明，提供此软件不构成对任何第三方软件许可限制、EULA 或其他合同条款的豁免或绕过授权。
5. 任何因滥用本软件、非法使用本软件或其衍生作品而导致的法律责任和后果均由使用者自行承担，本软件的作者及维护者对此不承担任何责任。
6. 本软件不对数据的准确性、完整性或适用性做出任何明示或暗示的保证，使用者应自行承担使用风险。
7. 使用者理解并同意，使用本软件即表示接受本免责声明的所有条款和条件，如有异议请立即停止使用本软件。

[ReleaseCard]: https://img.shields.io/github/v/release/Flyleague-Collection/chart-provider?style=for-the-badge&logo=github

[ReleaseDataCard]: https://img.shields.io/github/release-date/Flyleague-Collection/chart-provider?display_date=published_at&style=for-the-badge&logo=github

[LastCommitCard]: https://img.shields.io/github/last-commit/Flyleague-Collection/chart-provider?display_timestamp=committer&style=for-the-badge&logo=github

[BuildStateCard]: https://img.shields.io/github/actions/workflow/status/Flyleague-Collection/chart-provider/go-build.yml?style=for-the-badge&logo=github&label=Full-Build

[ProjectLanguageCard]: https://img.shields.io/github/languages/top/Flyleague-Collection/chart-provider?style=for-the-badge&logo=github

[ProjectLicense]: https://img.shields.io/badge/License-MIT-blue?style=for-the-badge&logo=github

[Release]: https://www.github.com/Flyleague-Collection/chart-provider/releases/latest

[Action]: https://github.com/Flyleague-Collection/chart-provider/actions/workflows/go-build.yml

[Release]: https://www.github.com/Flyleague-Collection/chart-provider/releases/latest

