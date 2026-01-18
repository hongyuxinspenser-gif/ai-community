# AI Community Backend (MVP)

这是一个基于 **Go (Golang)** 构建的高性能 AI 学习社区后端服务 MVP (Minimum Viable Product)。

该项目旨在提供一个轻量级、响应迅速的 API 服务，支持用户发布帖子、进行多级嵌套评论（盖楼），并高效地以树状结构返回评论数据。

## 📋 目录
- [项目简介](#-项目简介)
- [快速开始](#-快速开始)
- [核心 API 示例](#-核心-api-示例)
- [技术选型与设计思考](#-技术选型与设计思考)
  - [为什么选择 Gin 框架？](#1-为什么选择-gin-框架)
  - [数据库设计](#2-数据库设计)
  - [核心算法：评论层级查询](#3-核心算法评论层级查询)
- [当前局限性](#-当前局限性)
- [百万级用户架构演进](#-百万级用户架构演进)

---

## 🚀 项目简介

本项目实现了社区的核心互动功能，重点解决了**无限层级评论**的存储与检索问题。

**主要功能：**
* 📝 **帖子管理**：发布、查询、删除帖子。
* 💬 **多级评论**：支持无限层级的“楼中楼”回复（Reply to Reply）。
* 🌳 **树状返回**：后端自动组装 JSON 树，前端无需递归处理，直接渲染。
* 📄 **智能分页**：基于“顶级评论（楼主）”进行分页，保证对话上下文完整性。

**技术栈：**
* **语言**: Go 1.21+
* **Web 框架**: Gin
* **ORM**: GORM
* **数据库**: SQLite (轻量级，易于部署)

---

## 🛠 快速开始

### 环境要求
* Go 1.18 或更高版本
* Git

### 运行步骤

1.  **克隆项目**
    ```bash
    git clone [https://github.com/YourUsername/ai-community-backend.git](https://github.com/YourUsername/ai-community-backend.git)
    cd ai-community-backend
    ```

2.  **安装依赖**
    go mod tidy

3.  **运行服务**
    go run main.go
    *服务默认监听在 `:8080` 端口。项目会自动在根目录生成 `community.db` 数据库文件。*

---

## 📡 核心 API 示例

以下示例使用 `curl` 命令，你也可以使用 Postman 进行测试。

### 1. 发布帖子
curl -X POST http://localhost:8080/api/posts \
-H "Content-Type: application/json" \
-d '{"title":"AI 社区上线啦", "content":"欢迎大家讨论 Go 语言开发"}'

### 2. 发布顶级评论
curl -X POST http://localhost:8080/api/comments \
-H "Content-Type: application/json" \
-d '{"post_id": 1, "content": "前排支持！"}'

### 3. 回复某条评论 (构建层级)
curl -X POST http://localhost:8080/api/comments \
-H "Content-Type: application/json" \
-d '{"post_id": 1, "parent_id": 1, "content": "层主说得对 (二级回复)"}'

### 4. 获取树状评论列表 (核心功能)
curl "http://localhost:8080/api/posts/1/comments?page=1&page_size=10"


##  技术选型与设计思考
### 1. 为什么选择 Gin 框架？

高性能：Gin 基于 httprouter，速度极快，内存占用低，非常适合处理高并发的社区 API。

生态成熟：拥有丰富的中间件支持（日志、CORS、Recover 等）。

开发效率：JSON 绑定、路由分组等 API 设计简洁直观。

### 2. 数据库设计

我们采用了 邻接表 (Adjacency List) 模式设计 comments 表：

字段名	类型	说明
id	uint	主键
post_id	uint	归属的帖子
parent_id	*uint	关键字段。指向父评论 ID，为 NULL 时代表顶级评论。
content	string	内容
设计理由：这种结构最简单，写入极快（O(1)），且天然支持无限层级。虽然递归查询复杂，但通过应用层算法可以完美解决。

### 3. 核心算法：评论层级查询

为了避免数据库的 N+1 查询问题，我们采用 "一次查询 + 内存组装" 的策略：

Fetch: 根据 post_id 一次性取出该帖子下的所有评论。

Map: 建立 ID -> Comment指针 的哈希映射。

Attach: 遍历评论列表，如果 ParentID 存在，直接通过 Map 找到父节点，将自己追加到父节点的 Children 切片中。

Filter: 最后只返回 ParentID 为空的节点列表（即顶级评论）。

该算法的时间复杂度为 O(N)，效率极高。

## 当前局限性
作为 MVP 版本，目前存在以下局限：

### 数据库锁：SQLite 是文件级锁，高并发写入时性能会下降。

### 全量查询：构建评论树时，目前是加载帖子下所有评论到内存。如果单贴评论数超过 10万+，内存压力会变大。

### 缺失鉴权：尚未集成 JWT 用户登录系统，任何人都可以发帖。

## 百万级用户架构演进
如果要支撑百万级 DAU，系统架构需做如下升级：

### 数据库迁移：

将 SQLite 迁移至 MySQL 或 PostgreSQL。

实施读写分离，主库写，从库读。

### 缓存策略 (Redis)：

热点缓存：将热门帖子及其评论树结构缓存到 Redis 中（Key: post:comments:tree:1001）。

写回策略：新评论写入数据库后，异步更新或失效缓存。

### 评论存储优化：

对于超长评论楼层，改用 闭包表 (Closure Table) 或 物化路径 (Materialized Path) 存储，优化子树查询效率。

或者使用 MongoDB 这种文档型数据库，直接存储嵌套结构。

### 微服务拆分：

将 PostService 和 CommentService 拆分为独立服务，通过 gRPC 通信。

引入 Elasticsearch 提供全文检索功能。

### 基础设施：

使用 Docker + Kubernetes 进行容器化部署与自动扩容。

接入 Nginx 负载均衡。

