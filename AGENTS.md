# YDSterm - 项目规则

## 技术栈
- **桌面框架**: Wails v3 (Go + WebView)
- **前端**: Vue 3.5 + TypeScript + Tailwind CSS v4 + Vue Router
- **后端**: Go 1.25+ 配合 GoFrame v2 (gdb ORM, gres 资源管理)
- **数据库**: SQLite，使用 GoFrame sqlite 驱动

## 文档结构（强制）

编写任何代码前，必须先阅读以下文档：

```
docs/
├── architecture.md          # 项目整体架构
├── development-plan.md      # 开发计划
├── PRD.md                   # 产品需求文档
└── modules/
    ├── README.md            # 模块文档索引
    ├── host-management.md   # 主机管理模块
    └── terminal.md          # 终端/SSH 模块
```

**规则**：
- 修改或新增功能时，必须先阅读对应的模块文档
- 如需超出既有架构或修改架构，必须先向用户确认
- 模块文档应保持与代码同步更新

## 构建与开发

```bash
make gen-dao      # 数据库表结构变更后，重新生成 DAO 代码
make pack         # 将迁移文件嵌入二进制（迁移文件变更后必须运行）
make dev          # 启动 wails3 开发服务器
make build        # 生产构建 (pack + wails3 build)
make tidy         # go mod tidy + npm install
make clean        # 清理构建产物
```

## 数据库规则

### DAO 模式（强制）
始终使用 `dao.Xxx.方法()` 进行数据库操作，严禁直接使用 `g.DB().Model(...)`。

```go
// 正确:
hosts, err := dao.YdstermHosts.List(ctx, groupID)
dao.YdstermHosts.Create(ctx, input)

// 错误:
g.DB().Model("ydsterm_hosts").Ctx(ctx).Scan(&hosts)
```

### DAO 目录结构
```
internal/dbcore/
├── dao/                    # 自定义 DAO（在此添加业务方法）
│   ├── ydsterm_hosts.go    # 公开 DAO，全局变量 YdstermHosts
│   └── internal/           # 自动生成的基础 DAO（禁止编辑）
├── model/
│   ├── entity/             # 数据库表结构体（自动生成，禁止编辑）
│   └── do/                 # 查询用 Data Objects（自动生成，禁止编辑）
```

### 新增表的流程
1. 在 `manifest/migrations/` 中创建迁移 SQL 文件
2. 运行 `make pack` 嵌入迁移文件
3. 启动一次应用以应用迁移（创建表）
4. 运行 `make gen-dao` 生成 DAO/model 代码
5. 在 `dao/` 中为生成的 DAO 添加自定义方法

### 迁移规则
- 文件命名: `manifest/migrations/{6位序号}_{描述}.up.sql`
- 使用 `CREATE TABLE IF NOT EXISTS` 保证幂等性
- 迁移记录在 `ydsterm_schema_migrations` 表中
- 每次应用启动时自动执行未应用的迁移

## 服务层

```
internal/service/service.go   # 接口定义 (IHostService 等)
internal/service/host.go      # 接口实现（委托给 dao.Xxx）
internal/service/ssh.go       # SSH/终端服务
internal/service/convert.go   # entity -> types 转换
```

服务在 `internal/boot/app.go` 中注册，通过 Wails bindings 暴露给前端。

## 开发计划

| 阶段 | 内容 | 说明 |
|------|------|------|
| **Phase 1** | Host 管理 | SSH 主机 CRUD + 分组管理 + 密钥对管理/生成 |
| **Phase 2** | 终端连接 | SSH 连接建立 + xterm.js 终端模拟器 + 多 Tab 会话 |
| **Phase 3** | SFTP | 双面板文件浏览器 + 拖拽传输 + 进度显示 |
| **Phase 4** | 端口转发 | Local/Remote/Dynamic 三种模式 + 开关控制 |
| **Phase 5** | Snippets | 代码片段 CRUD + 文件夹管理 + 快速插入终端 |
| **Phase 6** | 打磨 | 测试 + 打包 + 文档 |

## 加密
- SSH 密码和私钥使用 AES-256-GCM 加密存储
- 加密密钥从本机 hostname 派生
- `internal/crypto/crypto.go` 提供 Encrypt/Decrypt 方法

## ID 生成
- 基于 ULID 的 ID，通过 `internal/types/id.go` (`types.NewID()`) 生成
- 依赖包: `github.com/oklog/ulid/v2`

---

## Wails v3 开发踩坑记录

### 1. Event 必须先 RegisterEvent 再 Emit
Wails v3 中，事件必须在 `main.go` 的 `init()` 中通过 `application.RegisterEvent[T]("event-name")` 注册，否则 `Event.Emit()` 静默丢弃数据。

```go
// main.go init()
func init() {
    application.RegisterEvent[TermOutputPayload]("term-output")
    application.RegisterEvent[TermDisconnectedPayload]("term-disconnected")
}
```

### 2. Event Payload 被 {name, data} 信封包裹
前端 `Events.On("event-name", callback)` 收到的不是裸 payload，而是：
```json
{"name": "event-name", "data": { /* 实际 payload */ }}
```
必须通过 `payload.data.SessionId` 访问字段，而非 `payload.SessionId`。

### 3. protobuf 序列化忽略 JSON tag
Wails v3 使用 protobuf 序列化，**JSON tag 在事件序列化层被忽略**。Go struct 字段名直接作为 wire format 的 key（PascalCase）。

```go
// 正确：Go 字段名 = JS 访问的 key
type TermOutputPayload struct {
    SessionId string  // → frontend: inner.SessionId
    Data      string  // → frontend: inner.Data
}
// 不要依赖 json:"sessionId" tag
```

### 4. 绑定必须与运行时二进制匹配
Wails bindings ( `frontend/bindings/` 目录) 中的方法 ID 是 Go 方法签名的哈希。Go 代码变更后必须重新 `wails3 build` 或 `wails3 dev` 以重新生成绑定，否则出现 `unknown bound method id` 错误。

### 5. xterm.js 在 WebView 中需要显式 focus
在 Wails WebView2 中，xterm 的隐藏 textarea 不会自动获得焦点。必须在 `terminal.open()` 后调用 `terminal.focus()`，并在容器上绑定 click/mousedown 事件重新聚焦。

### 6. Vue v-for 内部 ref 只绑定最后一个元素
`<div v-for="item in list" ref="myRef" />` 中，`myRef` 只引用循环的最后一个元素。多实例场景必须用函数 ref：
```html
<div v-for="tab in tabs" :ref="(el) => containerEls[tab.id] = el" />
```

### 7. TypeScript moduleResolution 必须用 Bundler
`tsconfig.json` 中 `moduleResolution` 必须设为 `"Bundler"`（非 `"Node"`），否则无法解析 Wails bindings 的无扩展名 import。

### 8. make dev 必须依赖 make pack
`Makefile` 中 `dev` 目标必须先执行 `pack`，确保迁移文件通过 `gres` 嵌入：
```makefile
dev: pack
    wails3 dev
```
