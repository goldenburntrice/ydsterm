# YDSterm - Product Requirements Document

## 1. 产品概述

**YDSterm** 是一个跨平台 SSH 终端客户端，对标 Termius，使用 Wails v3 + Vue 3.5 + Tailwind CSS 构建。支持 Windows、Linux、macOS。

### 1.1 技术栈

| 层级 | 技术 |
|------|------|
| 桌面框架 | Wails v3 (Go + WebView2/WebKitGTK) |
| 后端语言 | Go 1.25+ |
| 前端框架 | Vue 3.5 + Composition API + TypeScript |
| CSS 框架 | Tailwind CSS v4 |
| 数据库 | SQLite (modernc.org/sqlite - 纯 Go，无 CGO) |
| SSH 协议 | golang.org/x/crypto/ssh |
| SFTP | github.com/pkg/sftp |
| 终端模拟 | xterm.js |
| 路由 | 无（Tab 导航模式，useAppTabs 组合函数管理） |
| 工具库 | @vueuse/core |

### 1.2 数据库迁移策略

参考 aihelper 项目的迁移方案：
- SQL 迁移文件存放在 `manifest/migrations/` 目录下
- 文件命名：`{6位序号}_{描述}.up.sql` / `.down.sql`
- 通过 `gf pack` 将迁移文件嵌入二进制（gres 资源管理）
- 启动时自动检测并应用未执行的迁移
- 使用 `ydsterm_schema_migrations` 表追踪已应用的迁移版本

---

## 2. 功能需求

### 2.1 Host 管理

管理 SSH 连接配置和密钥。

**数据结构：**

```sql
ydsterm_hosts
├── id (TEXT PRIMARY KEY)          -- ULID
├── name (TEXT NOT NULL)           -- 主机名称
├── hostname (TEXT NOT NULL)       -- IP/域名
├── port (INTEGER DEFAULT 22)      -- SSH 端口
├── username (TEXT)                 -- 用户名
├── auth_method (TEXT)              -- password / private_key
├── password_enc (TEXT)             -- 加密后的密码
├── key_id (TEXT)                   -- 关联 SSH key
├── group_id (TEXT)                 -- 所属分组
├── color (TEXT)                    -- 标记颜色
├── created_at (DATETIME)
├── updated_at (DATETIME)
└── deleted_at (DATETIME)           -- 软删除

ydsterm_keys
├── id (TEXT PRIMARY KEY)
├── name (TEXT NOT NULL)            -- 密钥名称
├── private_key_enc (TEXT NOT NULL) -- 加密后的私钥
├── public_key (TEXT)               -- 公钥
├── passphrase_enc (TEXT)           -- 加密后的密钥密码
├── created_at (DATETIME)
└── updated_at (DATETIME)

ydsterm_host_groups
├── id (TEXT PRIMARY KEY)
├── name (TEXT NOT NULL)
├── parent_id (TEXT)                -- 父分组（支持嵌套）
├── sort_order (INTEGER DEFAULT 0)
├── created_at (DATETIME)
└── updated_at (DATETIME)
```

**功能：**
- 创建/编辑/删除 SSH 主机配置
- 支持密码和密钥两种认证方式
- SSH 密钥的导入/导出
- 密钥对生成（RSA/Ed25519）
- 主机分组管理（支持嵌套分组）
- 主机搜索与过滤

### 2.2 SFTP 文件传输

多窗口 SFTP，支持拖拽文件传输。

> SFTP 会话纯内存管理，不持久化到数据库。

**功能：**
- 多 Tab/窗口 SFTP 会话
- 双面板布局（本地文件系统 + 远程文件系统）
- 文件树浏览（展开/折叠目录）
- 文件上传/下载（支持拖拽）
- 传输进度显示
- 文件权限查看与修改（chmod）
- 文件操作：新建文件夹、重命名、删除
- 批量传输支持

### 2.3 端口转发

支持三种端口转发模式（同 Termius）。

**数据结构：**

```sql
ydsterm_port_forwards
├── id (TEXT PRIMARY KEY)
├── name (TEXT)                      -- 规则名称
├── host_id (TEXT NOT NULL)          -- 关联主机
├── type (TEXT NOT NULL)             -- local / remote / dynamic
├── local_address (TEXT)             -- 本地绑定地址 (local/remote模式)
├── local_port (INTEGER)             -- 本地端口 (local/remote模式)
├── remote_host (TEXT)               -- 远程目标主机 (local/remote模式)
├── remote_port (INTEGER)            -- 远程目标端口 (local/remote模式)
├── socks_host (TEXT)                -- SOCKS绑定地址 (dynamic模式)
├── socks_port (INTEGER)             -- SOCKS绑定端口 (dynamic模式)
├── enabled (INTEGER DEFAULT 0)      -- 是否启用
├── created_at (DATETIME)
└── updated_at (DATETIME)
```

**三种模式说明：**

| 模式 | 说明 | 示例场景 |
|------|------|----------|
| Local | 将本地端口转发到远程主机:端口 | 本地 3306 → 远程数据库 3306 |
| Remote | 将远程端口转发到本地主机:端口 | 远程 8080 → 本地开发服务器 8080 |
| Dynamic | SOCKS 代理 | 浏览器通过远程主机访问内网 |

**功能：**
- 创建/编辑/删除端口转发规则
- 启用/禁用开关
- 连接状态监控
- 支持绑定的地址配置（127.0.0.1 / 0.0.0.0）

### 2.4 Snippets 脚本管理

代码片段管理，快速插入终端。

**数据结构：**

```sql
ydsterm_snippets
├── id (TEXT PRIMARY KEY)
├── name (TEXT NOT NULL)             -- 片段名称
├── content (TEXT NOT NULL)          -- 脚本内容
├── language (TEXT)                  -- 语言类型标识
├── folder_id (TEXT)                 -- 所属文件夹
├── sort_order (INTEGER DEFAULT 0)
├── created_at (DATETIME)
└── updated_at (DATETIME)

ydsterm_snippet_folders
├── id (TEXT PRIMARY KEY)
├── name (TEXT NOT NULL)
├── parent_id (TEXT)                 -- 父文件夹（支持嵌套）
├── sort_order (INTEGER DEFAULT 0)
├── created_at (DATETIME)
└── updated_at (DATETIME)
```

**功能：**
- 创建/编辑/删除代码片段
- 文件夹分类管理（支持嵌套）
- 快速搜索片段
- 一键插入到终端
- 导入/导出片段
- 支持多行脚本

---

## 3. 非功能需求

### 3.1 安全性
- SSH 私钥和密码使用 AES-256-GCM 加密存储
- 加密密钥派生自系统设备指纹
- 内存中敏感数据使用后立即清除

### 3.2 跨平台
- Windows: WebView2 (系统自带)
- Linux: WebKitGTK 6.0
- macOS: WebKit (系统自带)

### 3.3 数据持久化
- SQLite 数据库文件存储在用户数据目录
- Windows: `%APPDATA%/ydsterm/data/ydsterm.db`
- Linux: `~/.local/share/ydsterm/data/ydsterm.db`
- macOS: `~/Library/Application Support/ydsterm/data/ydsterm.db`

### 3.4 升级兼容
- 每次启动自动执行数据库迁移
- 保证升级后数据库结构一致
- 迁移失败时回滚并阻止启动（防止数据损坏）

---

## 4. 项目结构

```
ydsterm/
├── main.go                    # 应用入口：初始化 DB、迁移、注册事件/服务、启动 Wails
├── go.mod
├── Makefile                   # 构建命令集 (pack, dev, build, gen-dao, tidy, clean)
├── Taskfile.yml               # Wails 任务定义
├── internal/
│   ├── boot/
│   │   ├── app.go             # App 结构体，聚合所有 Wails Service
│   │   └── migrations.go      # 数据库迁移引擎（gres + gdb）
│   ├── database/
│   │   └── database.go        # SQLite 初始化（GoFrame gdb）
│   ├── dbcore/                # GoFrame DAO 层
│   │   ├── dao/               # 自定义 DAO（业务方法）
│   │   │   ├── ydsterm_hosts.go
│   │   │   ├── ydsterm_keys.go
│   │   │   ├── ydsterm_host_groups.go
│   │   │   ├── ydsterm_snippets.go
│   │   │   ├── ydsterm_snippet_folders.go
│   │   │   ├── ydsterm_port_forwards.go
│   │   │   ├── ydsterm_schema_migrations.go
│   │   │   └── internal/      # 自动生成的基础 DAO（禁止编辑）
│   │   └── model/
│   │       ├── entity/        # 数据库表结构体（自动生成）
│   │       └── do/            # Data Objects（自动生成）
│   ├── service/               # 服务层（Wails 绑定的入口）
│   │   ├── service.go         # 接口定义 (IHostService 等)
│   │   ├── host.go            # HostServiceImpl (主机 CRUD + 密钥管理)
│   │   ├── ssh.go             # TerminalServiceImpl (SSH 连接/终端)
│   │   ├── sshconfig.go       # buildSSHConfig 共享辅助函数
│   │   ├── sftp.go            # SFTPServiceImpl (SFTP 文件传输)
│   │   ├── snippet.go         # SnippetServiceImpl (片段 CRUD)
│   │   ├── portforward.go     # PortForwardServiceImpl (端口转发 CRUD)
│   │   └── convert.go         # entity ↔ types 类型转换
│   ├── types/                 # 共享类型/DTO
│   │   ├── id.go              # ULID ID 生成器
│   │   ├── host.go
│   │   ├── snippet.go
│   │   └── portforward.go
│   └── crypto/
│       └── crypto.go          # AES-256-GCM 加密/解密
├── manifest/
│   └── migrations/            # SQL 迁移文件（gres 嵌入）
│       ├── 000001_init.up.sql
│       └── 000001_init.down.sql
├── packed/
│   └── packed.go              # gf pack 生成的 gres 嵌入代码
├── hack/
│   └── config.yaml            # gf gen dao 数据库连接配置
└── frontend/
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    ├── tsconfig.json
    ├── bindings/               # Wails 自动生成的 TS 绑定（禁止编辑）
    └── src/
        ├── main.ts             # Vue 应用入口
        ├── App.vue             # 根组件（TabBar + 视图切换）
        ├── styles/
        │   └── main.css        # Tailwind + 设计变量
        ├── views/              # 页面级组件
        │   ├── HomeView.vue    # 主页（主机管理 + 搜索 + 分组树）
        │   ├── TerminalView.vue # SSH 终端
        │   ├── SFTPView.vue    # SFTP 双面板文件浏览器
        │   ├── PortForwardView.vue
        │   └── SnippetsView.vue
        ├── components/         # 可复用组件
        │   ├── TabBar.vue      # 顶部 Tab 导航栏（主题切换）
        │   ├── HostFormDialog.vue
        │   ├── KeyManagerDialog.vue
        │   └── ContextMenu.vue
        ├── composables/        # 组合函数
        │   ├── useAppTabs.ts   # Tab 生命周期管理
        │   ├── useModal.ts     # 弹窗状态管理
        │   └── useTheme.ts     # 暗色/亮色主题切换
        └── types/              # TypeScript 类型定义（预留）
```

---

## 5. 开发计划

详见 [development-plan.md](./development-plan.md)。
