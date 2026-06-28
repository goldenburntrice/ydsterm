# 主机管理模块 (Host Management)

## 概述

Host 管理模块负责 SSH 主机配置的增删改查、分组管理、SSH 密钥对管理。对应开发计划 Phase 1，状态：**已完成**。

## 数据表

### ydsterm_hosts

| 字段 | 类型 | 说明 |
|------|------|------|
| id | TEXT PK | ULID |
| name | TEXT | 主机名称 |
| hostname | TEXT | IP/域名 |
| port | INTEGER | SSH 端口，默认 22 |
| username | TEXT | 登录用户名 |
| auth_method | TEXT | 认证方式: `password` / `private_key` |
| password_enc | TEXT | AES-256-GCM 加密的密码 |
| key_id | TEXT | 关联的 SSH 密钥 ID |
| group_id | TEXT | 所属分组 ID |
| color | TEXT | 标记颜色 (#hex) |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |
| deleted_at | DATETIME | 软删除时间 |

### ydsterm_keys

| 字段 | 类型 | 说明 |
|------|------|------|
| id | TEXT PK | ULID |
| name | TEXT | 密钥名称 |
| private_key_enc | TEXT | AES-256-GCM 加密的私钥 |
| public_key | TEXT | 公钥 |
| passphrase_enc | TEXT | 加密的 passphrase |
| created_at | DATETIME | |
| updated_at | DATETIME | |

### ydsterm_host_groups

| 字段 | 类型 | 说明 |
|------|------|------|
| id | TEXT PK | ULID |
| name | TEXT | 分组名称 |
| parent_id | TEXT | 父分组 ID（支持嵌套） |
| sort_order | INTEGER | 排序 |
| created_at | DATETIME | |
| updated_at | DATETIME | |

## 后端实现

### 服务层

`internal/service/host.go` — `HostServiceImpl`

| 方法 | 签名 | 说明 |
|------|------|------|
| Create | `(HostCreateInput) → Host` | 创建主机配置，密码自动加密 |
| Update | `(HostUpdateInput) → Host` | 部分更新，nil 字段不修改 |
| Delete | `(id) → error` | 删除主机 |
| Get | `(id) → Host` | 获取单个主机 |
| List | `(groupID) → []Host` | 列出主机，可按分组过滤 |
| CreateKey | `(KeyCreateInput) → Key` | 导入/创建密钥，私钥加密存储 |
| UpdateKey | `(KeyUpdateInput) → Key` | 更新密钥 |
| DeleteKey | `(id) → error` | 删除密钥 |
| GetKey | `(id) → Key` | 获取单个密钥 |
| ListKeys | `() → []Key` | 列出所有密钥（不含私钥） |
| CreateGroup | `(HostGroupCreateInput) → HostGroup` | 创建分组 |
| UpdateGroup | `(HostGroupUpdateInput) → HostGroup` | 更新分组 |
| DeleteGroup | `(id) → error` | 删除分组 |
| ListGroups | `() → []HostGroup` | 列出所有分组 |

### DAO 层

`internal/dbcore/dao/ydsterm_hosts.go` — `YdstermHosts`
- `List(ctx, groupID)` — WHERE deleted_at IS NULL, 可选 groupID 过滤
- `Get(ctx, id)` — 按主键查询
- `Create(ctx, input)` — 生成 ULID + 插入
- `Update(ctx, id, data)` — 按主键更新
- `Delete(ctx, id)` — 按主键删除

`internal/dbcore/dao/ydsterm_keys.go` — `YdstermKeys`
- `List(ctx)`, `Get(ctx, id)`, `Create`, `Update`, `Delete`

`internal/dbcore/dao/ydsterm_host_groups.go` — `YdstermHostGroups`
- `List(ctx)`, `Get`, `Create`, `Update`, `Delete`

### 加密

密码和私钥在 Service 层调用 `crypto.Encrypt()` 后存入数据库，读取时调用 `crypto.Decrypt()` 解密。DAO 层不感知加密逻辑。

## 前端实现

### 页面

`views/HomeView.vue` — 主页视图（包含主机管理功能）：
- **顶部**：YDSterm Logo + 搜索栏 + 操作按钮（新建主机 / 主题切换 / 密钥管理）
- **主体**：分组树视图（展开/折叠），分组内主机条目（颜色标记 / hover 操作按钮）
- **右键菜单**：SSH 连接 / SFTP 连接
- **弹窗**：HostFormDialog.vue + KeyManagerDialog.vue

### 组件

| 组件 | 文件 | 说明 |
|------|------|------|
| HostFormDialog | `components/HostFormDialog.vue` | 新建/编辑主机弹窗，含颜色选择、认证切换、分组选择（支持内联新建分组） |
| KeyManagerDialog | `components/KeyManagerDialog.vue` | 密钥管理弹窗，支持导入私钥、查看/删除已有密钥 |
| ContextMenu | `components/ContextMenu.vue` | 通用右键菜单组件 |

### 状态管理

主机管理状态直接在 `HomeView.vue` 中通过 `ref` 管理，无独立 composable：

```ts
const hosts = ref<Host[]>([])
const groups = ref<HostGroup[]>([])
const keys = ref<any[]>([])
```

所有 HostService 的 Wails 绑定调用直接内联在 HomeView.vue 的函数中。

### Wails 绑定

前端调用自动生成的绑定文件：
```ts
import * as HostService from '../../bindings/ydsterm/internal/service/hostserviceimpl'
import type { Host, HostCreateInput, ... } from '../../bindings/ydsterm/internal/types/models'
```

## 使用流程

1. 点击「添加主机」→ 填写名称、地址、端口、用户名
2. 选择认证方式（密码 / 密钥）
3. 如需密钥认证，先在「密钥」管理中导入私钥
4. 可选：分配到分组、选择标记颜色
5. 双击主机卡片 → 跳转终端页发起 SSH 连接

## 扩展点

- 密钥对生成（RSA/Ed25519）— Go `crypto/ssh` 可实现
- 主机导入/导出（JSON 格式）
- 分组拖拽排序
