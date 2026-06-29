# 云配置同步功能交接文档

## 功能概述

云配置同步功能允许用户在多台设备间同步 SSH 主机配置、密钥、代码片段等数据。采用客户端-服务端架构，支持增量同步和多用户管理。

## 架构设计

### 客户端（桌面应用）
- **框架**：Wails v3 + Vue 3 + Go
- **数据库**：SQLite（本地）
- **核心改动**：
  - 所有业务表增加 `sync_user` 字段，实现数据隔离
  - DAO 层自动注入当前用户上下文
  - 新增设置面板（通用/云配置）
  - 新增用户管理功能

### 服务端（独立服务）
- **框架**：GoFrame v2
- **数据库**：SQLite
- **API 端点**：
  - `POST /api/auth/verify` - 用户验证/注册
  - `GET /api/sync/pull` - 拉取变更
  - `POST /api/sync/push` - 推送变更
- **认证**：所有请求需携带 `X-Server-Key` 请求头

## 实现状态

### 已完成

#### 1. 数据库迁移
- ✅ 迁移文件 `000002_add_sync_user.up.sql`
- ✅ 6 张业务表添加 `sync_user` 字段
- ✅ 新增 `ydsterm_sync_users` 表（用户管理）
- ✅ 新增 `ydsterm_settings` 表（本地设置）
- ✅ 默认创建 `LOCALUSER` 账户

#### 2. 数据隔离层
- ✅ `internal/dbcore/context.go` - 用户上下文管理
- ✅ 所有 DAO 方法自动过滤 `sync_user`
- ✅ 生成新的 DAO/Entity 代码

#### 3. 服务端实现
- ✅ GoFrame 框架搭建
- ✅ 数据库初始化（8 张表）
- ✅ 用户认证（bcrypt 密码哈希）
- ✅ 增量同步 API（pull/push）
- ✅ Server Key 认证中间件
- ✅ 自动生成 DAO 层代码

#### 4. 客户端服务层
- ✅ `SettingsService` - 设置管理
- ✅ `UserService` - 用户管理
- ✅ Wails 绑定注册

#### 5. 前端界面
- ✅ 设置面板组件 `SettingsPanel.vue`
- ✅ TabBar 添加设置按钮（齿轮图标）
- ✅ 云配置表单（服务器地址、密钥、用户名、密码）
- ✅ 验证并登录功能

### 待完成

#### 1. 同步引擎（优先级：高）
- [ ] 实现客户端 Push 逻辑（本地变更 → 服务端）
- [ ] 实现客户端 Pull 逻辑（服务端变更 → 本地）
- [ ] 敏感数据加解密（SSH 密码、私钥）
- [ ] 冲突解决策略（last-write-wins）

#### 2. 账户切换（优先级：高）
- [ ] 主页显示当前用户
- [ ] 账户切换 UI
- [ ] 切换时自动同步

#### 3. 自动同步（优先级：中）
- [ ] 本地修改后自动 push（防抖 2s）
- [ ] 同步状态指示器
- [ ] 网络异常处理

#### 4. 测试验证（优先级：高）
- [ ] 端到端测试
- [ ] LOCALUSER 数据隔离验证
- [ ] 多用户切换测试

## 关键文件清单

### 客户端
```
manifest/migrations/000002_add_sync_user.up.sql  # 数据库迁移
internal/dbcore/context.go                        # 用户上下文
internal/dbcore/dao/ydsterm_sync_users.go        # 用户 DAO
internal/dbcore/dao/ydsterm_settings.go          # 设置 DAO
internal/service/user.go                          # 用户服务
internal/service/settings.go                      # 设置服务
internal/types/sync.go                            # 同步类型定义
frontend/src/components/SettingsPanel.vue        # 设置面板
```

### 服务端
```
server/
├── main.go                                       # 入口
├── go.mod                                        # 依赖
├── hack/config.yaml                              # GoFrame CLI 配置
├── manifest/
│   ├── config/config.yaml                        # 应用配置
│   └── migrations/000001_init.up.sql            # 数据库迁移
└── internal/
    ├── cmd/cmd.go                                # 命令入口
    ├── consts/consts.go                          # 常量
    ├── controller/sync.go                        # API 控制器
    ├── dao/                                      # 自动生成的 DAO
    └── model/                                    # 自动生成的 Model
```

## 启动方式

### 服务端
```bash
cd server
make server-dev
# 或
cd server && go run .
```

首次启动会自动：
1. 创建 `ydsterm-server.db` 数据库
2. 初始化所有表
3. 生成 Server Key（打印到日志）

### 客户端
```bash
make dev
```

## 测试流程

1. **启动服务端**
   ```bash
   make server-dev
   ```
   记录日志中的 Server Key

2. **启动客户端**
   ```bash
   make dev
   ```

3. **配置云连接**
   - 点击 TabBar 右侧齿轮图标
   - 切换到"云配置"标签
   - 填写：
     - 服务器地址：`http://localhost:8080`
     - 服务器密钥：从服务端日志复制
     - 用户名：任意（如 `alice`）
     - 密码：任意
   - 点击"验证并登录"

4. **验证**
   - 服务端日志应显示用户创建成功
   - 客户端应显示"验证成功，已登录"

## 数据流

### 本地模式（LOCALUSER）
```
用户操作 → Service → DAO（自动注入 sync_user=LOCALUSER）→ SQLite
```

### 云同步模式（alice）
```
用户操作 → Service → DAO（自动注入 sync_user=alice）→ SQLite
                                                    ↓
                                              SyncEngine（待实现）
                                                    ↓
                                              Server API
```

## 注意事项

1. **数据隔离**：所有 DAO 查询自动过滤 `sync_user`，无需手动添加
2. **LOCALUSER**：默认本地用户，数据不会同步到云端
3. **敏感数据**：当前传输明文，需实现加解密（SSH 密码、私钥）
4. **冲突处理**：采用 last-write-wins 策略，以 `updated_at` 为准
5. **服务端部署**：建议配置 HTTPS，避免密码明文传输

## 后续优化建议

1. **增量同步优化**：当前全量拉取，可优化为只传输变更字段
2. **离线支持**：网络断开时缓存变更，恢复后自动同步
3. **冲突提示**：检测到冲突时提示用户选择
4. **同步历史**：记录同步日志，支持回滚
5. **多设备管理**：支持查看和管理已连接的设备

## 相关文档

- [产品需求文档](./PRD.md)
- [系统架构](./architecture.md)
- [开发计划](./development-plan.md)
