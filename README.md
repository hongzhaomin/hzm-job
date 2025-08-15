# hzm-job 分布式任务调度平台

[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

## 1. 概述

hzm-job 是一个轻量级的分布式任务调度平台，专注于解决分布式环境下的任务调度与执行问题。它提供了任务注册、调度、执行、监控等核心功能，支持多种任务类型和调度策略。

### 核心特性

- **分布式架构**：支持多执行器节点部署，实现任务的分布式执行
- **灵活调度**：支持Cron表达式等多种调度方式
- **任务管理**：提供任务的注册、修改、删除、暂停等管理功能
- **实时监控**：可视化界面展示任务执行状态和日志
- **高可用性**：支持执行器节点的自动注册与发现
- **易于集成**：提供简洁的API接口，便于业务系统集成

## 2. 快速开始

### 2.1 环境要求

- Go 1.18+
- MySQL 5.7+
- 网络互通的多台服务器（分布式部署时）

### 2.2 项目结构

```
hzm-job/
├── admin/          # 管理端服务
├── client/         # 客户端SDK
├── core/           # 核心组件
└── example-executor/ # 执行器示例
```

1. 配置数据库连接信息（`admin/hzm-job.yaml`）
2. 初始化数据库表结构（参考`admin/system.sql`）


## 3. 基础配置

### 3.1 管理端配置

管理端配置文件位于`admin/hzm-job.yaml`：

```yaml
hzm:
  job:
    admin:
      port: 8888  # 管理端服务端口
      mysql:
        host: localhost
        port: 3306
        username: root
        password: password
        dbname: hzm_job
```

### 3.2 执行器配置

执行器配置文件位于`client/hzm-job.yaml`：

```yaml
hzm:
  job:
    client:
      port: 7777                    # 客户端健康检查端口
      adminAddress: http://localhost:8888  # 管理服务器地址
      appKey: your-app-key           # 应用程序唯一标识符，各执行器自己定义
      appSecret: your-app-secret     # 认证密钥，在管理端创建执行器时由管理端生成
```

## 4. 安装指南

### 4.1 获取源码

```bash
  git clone https://github.com/hongzhaomin/hzm-job.git
```

### 4.2 依赖管理

项目使用Go Modules进行依赖管理，构建时会自动下载所需依赖。

### 4.3 构建可执行文件

```bash
    # 构建管理端
    cd admin
    go build -o hzm-job-admin
    
    # 构建执行器示例
    cd example-executor
    go build -o hzm-job-executor
```

### 4.4 启动管理端

```bash
    # 可以任意更改配置文件位置，传入 f 参数指定即可 
    ./hzm-job-executor -f hzm-job.yaml
```

### 4.5 访问管理端web界面

- 访问调度中心web界面 `http://localhost:8888`
- 使用默认账户登录
  > 用户名：admin
  >
  > 密码：1

## 5. 创建首个任务

### 5.1 下载执行器sdk依赖

在您的 Go 项目中，添加 hzm-job 客户端依赖：

```bash
  go get github.com/hongzhaomin/hzm-job/client
```

### 5.2 在管理端配置执行器

1. 访问管理端Web界面（默认端口8888）
2. 进入执行器管理页面
3. 点击"新增执行器"
4. 填写执行器信息：
   - AppSecret：点击生成，系统会自动生成密钥；也可以不生成，那么执行器就无需配置密钥（客户端配置中的 `appSecret` 与之相同）
   - AppKey：选择对应的执行器（客户端配置中的 `appKey` 与之相同）
   - 执行器名称：执行器的描述性名称
   - 注册类型：设置为"自动"以进行自动注册；也可以手动注册，此时需要填写手动注册节点地址
5. 保存执行器配置

### 5.3 配置客户端配置文件

在项目根目录创建配置文件 hzm-job.yaml：

```yaml
  hzm:
  job:
    client:
      port: 7777                    # 客户端健康检查端口
      adminAddress: http://localhost:8888  # 管理服务器地址
      appKey: your-app-key           # 应用程序唯一标识符，各执行器自己定义
      appSecret: your-app-secret     # 认证密钥，在管理端创建执行器时由管理端生成

    common:
      log:
        level: debug
        type: json
```

### 5.4 定义任务

在执行器中定义任务处理逻辑：
 - 基于函数的任务（JobFunc）
   基于函数的作业是定义任务最简单的方式。它们非常适合快速、直接的操作，特别是当您不需要复杂的参数处理或希望避免创建单独的结构体时。JobFunc 本质上是一个遵循特定签名的函数：`func(ctx context.Context, param *string) error`。

```go
   package main

import (
   "context"
   "errors"
   "fmt"
   "github.com/hongzhaomin/hzm-job/client/hzmjob"
   "time"
)

func main() {

   // 注册任务
   hzmjob.AddJob("cancelableJobFuncTest", func(ctx context.Context, param *string) error {
      fmt.Println("====== cancelableJobFuncTest ========> 任务开始执行:", param)
      var count int
      for {
         select {
         case <-ctx.Done():
            fmt.Println("====== cancelableJobFuncTest ========> 任务取消:", param)
            return errors.New("任务[cancelableJobFuncTest]被调度中心终止")
         default:
            count++
            time.Sleep(time.Second * 5)
            fmt.Println(fmt.Sprintf("=== cancelableJobFuncTest ===> 模拟数据库操作，执行 %d 次", count))
         }
      }
   })
   // 启动
   hzmjob.Enable()
}
```

 - 基于结构体的任务（annotation.HzmJob[T]）

基于结构体的作业为定义任务提供了更强大且类型安全的方式。它们实现了 `annotation.HzmJob[Param any]` 接口，为您提供了强大的参数解析能力和更好的代码组织。这种方法推荐用于生产环境，其中可维护性和类型安全是优先考虑的因素。

```go
package main

import (
   "context"
   "fmt"
   "github.com/hongzhaomin/hzm-job/client/annotation"
   "github.com/hongzhaomin/hzm-job/client/hzmjob"
)

func main() {

   // 注册任务
   hzmjob.AddJobs(&MyJob{})
   // 启动
   hzmjob.Enable()
}

type RequestParam struct {
   Name string `json:"name"`
}

// 定义任务结构体
type MyJob struct {
   annotation.HzmJob[RequestParam] `name:"myJob"`
}

// 实现任务处理方法
func (j *MyJob) DoHandle(ctx context.Context, param *RequestParam) error {
   // 任务处理逻辑
   fmt.Println("执行任务参数:", param)
   return nil
}
```

### 5.5 在管理端配置任务

1. 访问管理端Web界面（默认端口8888）
2. 进入任务管理页面
3. 点击"新增任务"
4. 填写任务信息：
   - 任务名称：myJob
   - 执行器：选择对应的执行器
   - 调度类型：选择Cron或固定频率
   - 调度配置：填写Cron表达式或时间间隔
5. 保存任务配置

### 5.6 启动任务

在任务列表中找到刚创建的任务，点击"启动"按钮即可开始调度执行。

## 6. 架构设计

### 6.1 系统架构

hzm-job采用中心化管理架构，由管理端和多个执行器节点组成：

```
+------------------+     +------------------+
|   管理端服务      |<--->|   执行器节点1     |
+------------------+     +------------------+
        ^                        |
        |________________________|
        |                        |
+------------------+     +------------------+
|   Web管理界面    |     |   执行器节点2     |
+------------------+     +------------------+
```

### 6.2 核心组件

1. **管理端（Admin）**：负责任务调度、执行器管理、任务监控等
2. **执行器（Executor）**：负责任务的实际执行
3. **客户端SDK（Client）**：提供任务注册和执行器通信功能
4. **Web界面**：提供任务管理和监控的可视化界面

## 7. API接口

### 7.1 管理端API

管理端提供RESTful API接口供执行器调用：

- `POST /api/registry`：执行器注册接口
- `POST /api/offline`：执行器下线接口
- `POST /api/callback`：任务执行结果回调接口

### 7.2 客户端接口

客户端SDK提供简洁的接口供业务系统使用：

- `hzmjob.AddJob(name string, job Job)`：注册任务
- `hzmjob.Enable()`：启用任务调度
- `hzmjob.Close()`：关闭任务调度

## 8. 监控与日志

### 8.1 任务监控

管理端提供Web界面实时展示：

- 任务执行状态
- 执行器节点状态
- 任务执行日志
- 调度统计信息

### 8.2 日志配置

通过配置文件可以设置日志级别和输出格式：

```yaml
hzm:
  job:
    common:
      log:
        level: info  # 日志级别：debug/info/warn/error
        type: json   # 输出格式：text/json
```

## 9. 最佳实践

### 9.1 任务设计

- 任务处理逻辑应尽量简洁，避免长时间阻塞
- 合理设置任务超时时间
- 实现任务的幂等性，避免重复执行导致的问题

### 9.2 部署建议

- 执行器节点根据业务负载情况动态扩展
- 数据库建议使用主从复制提高可用性

## 10. 故障排查

### 10.1 常见问题

1. **执行器无法注册**：检查网络连接和配置信息是否正确
2. **任务不执行**：检查任务状态和调度配置
3. **数据库连接失败**：检查数据库配置和网络连接

### 10.2 日志分析

通过查看管理端和执行器的日志可以定位大部分问题。

## 11. 贡献指南

欢迎提交Issue和Pull Request来改进hzm-job项目。

## 12. 许可证

hzm-job采用Apache 2.0许可证，详见[LICENSE](LICENSE)文件。

## 13. 赞赏作者

![Image text](admin/web/static/images/laud-author.png)
