# YDSterm 架构文档

## 1. 项目概述

YDSterm 是基于 Wails v3 构建的跨平台 SSH 终端客户端。架构遵循 **Go 后端 + Web 前端** 分层模式，通过 Wails 的 Bindings（方法调用）和 Events（事件推送）实现前后端通信。

## 2. 项目目录结构

```
ydsterm/
├── main.go                       # 应用入口：DB 初始化、迁移、服务注册、Wails 启动
├── go.mod                        # Go 模块定义
├── Makefile                      # 构建命令集
├── Taskfile.yml                  # Wails 任务定义
│
├── internal/                     # ── Go 后端 ──
│   ├── boot/
│   │   ├── app.go                # App 结构体，聚合所有 Wails Service
│   │   └── migrations.go         # 数据库迁移引擎（gres + gdb）
│   ├── database/
│   │   └── database.go           # SQLite 初始化（GoFrame gdb）
│   ├── dbcore/                   # GoFrame DAO 层
│   │   ├── dao/                  # 自定义 DAO（业务方法）
│   │   │   ├── ydsterm_hosts.go
│   │   │   ├── ydsterm_keys.go
│   │   │   ├── ydsterm_host_groups.go
│   │   │   ├── ydsterm_snippets.go
│   │   │   ├── ydsterm_snippet_folders.go
│   │   │   ├── ydsterm_port_forwards.go
│   │   │   ├── ydsterm_schema_migrations.go
│   │   │   └── internal/         # 自动生成的基础 DAO（禁止编辑）
│   │   └── model/
│   │       ├── entity/           # 数据库表结构体（自动生成）
│   │       └── do/               # Data Objects（自动生成）
│   ├── service/                  # 服务层（Wails 绑定的入口）
│   │   ├── service.go            # 接口定义
│   │   ├── host.go               # 主机管理服务
│   │   ├── snippet.go            # 代码片段服务
│   │   ├── portforward.go        # 端口转发服务
│   │   ├── ssh.go                # SSH 终端服务
│   │   ├── sshconfig.go          # SSH 配置构建辅助函数
│   │   ├── sftp.go               # SFTP 文件传输服务
│   │   └── convert.go            # entity → types 类型转换
│   ├── types/                    # 共享类型/DTO
│   │   ├── id.go                 # ULID ID 生成器
│   │   ├── host.go
│   │   ├── snippet.go
│   │   └── portforward.go
│   └── crypto/
│       └── crypto.go             # AES-256-GCM 加密/解密
│
├── manifest/
│   └── migrations/               # SQL 迁移文件（gres 嵌入）
│       ├── 000001_init.up.sql
│       └── 000001_init.down.sql
│
├── packed/
│   └── packed.go                 # gf pack 生成的 gres 嵌入代码
│
├── hack/
│   └── config.yaml               # gf gen dao 数据库连接配置
│
├── frontend/                     # ── Vue 3.5 前端 ──
│   ├── index.html                # SPA 入口
│   ├── package.json
│   ├── vite.config.ts            # Vite + Vue + Tailwind + Wails 插件
│   ├── tsconfig.json             # TypeScript 配置
│   ├── bindings/                 # Wails 自动生成的 TS 绑定（禁止编辑）
│   └── src/
│       ├── main.ts               # Vue 应用入口
│       ├── App.vue               # 根组件（TabBar + 视图切换）
│       ├── styles/main.css       # Tailwind + 设计变量
│       ├── views/                # 页面级组件
│       │   ├── HomeView.vue      # 主页（主机管理 + 搜索 + 分组树）
│       │   ├── TerminalView.vue  # SSH 终端
│       │   ├── SFTPView.vue      # SFTP 双面板文件浏览器
│       │   ├── PortForwardView.vue
│       │   └── SnippetsView.vue
│       ├── components/           # 可复用组件
│       │   ├── TabBar.vue        # 顶部 Tab 导航栏（主题切换）
│       │   ├── HostFormDialog.vue
│       │   ├── KeyManagerDialog.vue
│       │   └── ContextMenu.vue
│       ├── composables/          # 组合函数
│       │   ├── useAppTabs.ts     # Tab 生命周期管理（打开/关闭/切换）
│       │   ├── useModal.ts       # 弹窗状态管理
│       │   └── useTheme.ts       # 暗色/亮色主题切换
│
└── docs/                         # 项目文档
    ├── architecture.md           # 本文档
    ├── development-plan.md
    ├── PRD.md
    └── modules/                  # 模块文档
```

## 3. 后端架构

### 3.1 启动流程

```
main()
  ├── boot.DefaultDBPath()        → 确定 SQLite 文件路径
  ├── boot.EnsureDataDir()        → 创建数据目录
  ├── database.Init()             → 通过 gdb.SetConfig() 配置 SQLite
  │   └── 驱动: sqlite::@file(path)
  │   └── 表前缀: ydsterm_
  ├── boot.RunMigrations()        → 从 gres 读取 .up.sql 并执行
  │   └── 追踪表: ydsterm_schema_migrations
  ├── boot.NewApp()               → 创建 App 实例（聚合所有 Service）
  └── app.Run()                   → 启动 Wails 应用
```

### 3.2 服务层 (Service Layer)

所有前端可调用的后端方法都在 `internal/service/` 中实现。每个 Service 是一个 Go struct，其导出方法自动生成 TypeScript 绑定。

```
internal/service/service.go   → IService 接口定义
internal/service/host.go      → HostServiceImpl    (主机 CRUD + 密钥管理)
internal/service/snippet.go   → SnippetServiceImpl (片段 CRUD)
internal/service/portforward.go → PortForwardServiceImpl (端口转发 CRUD)
internal/service/ssh.go       → TerminalServiceImpl (SSH 连接/终端)
internal/service/sshconfig.go → buildSSHConfig() 共享辅助函数
internal/service/sftp.go      → SFTPServiceImpl    (SFTP 文件传输)
internal/service/convert.go   → entity ↔ types 转换
```

**注册流程**：
1. `internal/boot/app.go` 创建各 Service 实例
2. `main.go` 将 Service 注册到 Wails:
   ```go
   Services: []application.Service{
       application.NewService(appInstance.HostService),
       application.NewService(appInstance.TerminalService),
       application.NewService(appInstance.SFTPService),
       application.NewService(appInstance.SnippetService),
       application.NewService(appInstance.PortForwardService),
       ...
   }
   ```

### 3.3 数据访问层 (DAO Layer)

使用 GoFrame `gf gen dao` 自动生成，遵循 DAO 模式：

```
internal/dbcore/
├── dao/           → 自定义 DAO（添加 List/Create/Update/Delete 方法）
└── dao/internal/  → 自动生成（Columns/Table/Ctx/Transaction，不可编辑）
```

**数据流**：Service → DAO → gdb.Model → SQLite

### 3.4 数据库迁移

采用类 aihelper 的迁移方案：
- 迁移 SQL 通过 `gf pack` 嵌入二进制（`packed/packed.go` 使用 gres）
- 启动时自动检测 `ydsterm_schema_migrations` 表，按版本号升序执行未应用的 `.up.sql`
- 每个迁移成功后记录版本号

## 4. 前端架构

### 4.1 导航

YDSterm 采用 **Tab 导航** 模式，不使用 Vue Router。所有视图通过 `useAppTabs` 组合函数管理，以 Tab 形式在 `App.vue` 中通过 `v-show` 切换。

```
App.vue
├── TabBar.vue              # 顶部 Tab 栏 + 主题切换
└── 视图区域（v-show 切换）
    ├── HomeView.vue (tab "主页", 不可关闭)
    │   ├── YDSterm Logo + 搜索栏
    │   ├── 操作按钮（新建主机 / 主题 / 密钥管理）
    │   ├── 分组树（展开/折叠）
    │   │   └── 主机条目（颜色标记 / hover 操作按钮）
    │   ├── 右键菜单（SSH 连接 / SFTP 连接）
    │   └── HostFormDialog.vue + KeyManagerDialog.vue (modal)
    ├── TerminalView.vue (动态 Tab, 可关闭)
    │   └── xterm.js 终端区域
    └── SFTPView.vue (动态 Tab, 可关闭)
        └── 双面板文件浏览器（本地 ↔ 远程）
```

### 4.2 组合函数 (Composables)

| 文件 | 函数 | 用途 |
|------|------|------|
| `useAppTabs.ts` | `useAppTabs()` | 全局 Tab 生命周期：打开/关闭/切换终端和 SFTP Tab，关闭时自动 Disconnect |
| `useModal.ts` | `useModal()` | visible/open/close 弹窗状态管理 |
| `useTheme.ts` | `useTheme()` | 暗色/亮色主题切换 + xterm.js 配色 |

参考 GitHub-dark 配色方案，使用 CSS 变量定义语义化 token：

```
--bg-base:      #0d1117   (主背景)
--bg-surface:   #161b22   (侧栏/面板)
--bg-elevated:  #1c2128   (弹窗/卡片)
--border:        #30363d
--accent:        #2f81f7
--text-primary:  #e6edf3
--text-secondary:#8b949e
--text-muted:    #6e7681
```

## 5. 通信模型

### 5.1 方法调用（Bindings）

```
前端 import { Create } from '...hostserviceimpl'
  → Create(input)
  → Wails protobuf IPC
  → Go HostServiceImpl.Create(input)
  → DAO 操作
  → 返回结果
  → 前端 Promise resolve
```

### 5.2 事件推送（Events）

```
Go goroutine 读取 SSH output
  → application.Get().Event.Emit("term-output", payload)
  → Wails protobuf IPC
  → 前端 Events.On("term-output", callback)
  → callback({ name: "term-output", data: { SessionId, Data, IsStderr } })
  → 前端 unwrap payload.data
  → xterm.write(atob(payload.data.Data))
```

**关键约定**：
- Go 中 `application.RegisterEvent[T](name)` 必须在前
- 前端收到的 payload 是 `{name, data}` 信封
- 事件 struct 的 Go 字段名直接映射到 JS 属性名（PascalCase）

已注册事件：

| 事件名 | Payload 类型 | 用途 |
|--------|-------------|------|
| `term-output` | `TermOutputPayload` | 终端输出流 |
| `term-disconnected` | `TermDisconnectedPayload` | 终端断开通知 |
| `sftp-transfer-progress` | `SftpTransferPayload` | SFTP 传输进度 |

## 6. 安全

- SSH 密码和私钥通过 AES-256-GCM 加密后存入 SQLite
- 加密密钥从本机 hostname 派生 (`internal/crypto/crypto.go`)
- 私钥和密码的明文仅存于内存，入库前即加密
