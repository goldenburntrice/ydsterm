# YDSterm 开发计划

## Phase 1: Host 管理

### 后端（已完成）
- [x] Host/Key/Group CRUD DAO (`internal/dbcore/dao/`)
- [x] HostService 实现（加密存储密码/私钥）
- [x] Wails bindings 自动生成

### 前端（已完成）
- [x] HomeView 主页面（分组树 + 搜索 + 操作按钮）
- [x] 主机表单弹窗（新增/编辑 Host）
- [x] 主机分组管理（新建/重命名/删除分组）
- [x] SSH 密钥管理面板（导入/查看/删除密钥）
- [x] 主机右键菜单（SSH 连接 / SFTP 连接）
- [x] 主机搜索过滤

---

## Phase 2: 终端连接

### 后端（已完成）
- [x] SSH 连接管理（connect/disconnect/write/resize）
- [x] PTY 会话（伪终端，支持 resize）
- [x] 通过 Wails Events 推送终端输出
- [x] 多会话管理

### 前端（已完成）
- [x] xterm.js 终端模拟器集成
- [x] 多 Tab 终端会话
- [x] Tab 切换/关闭/重命名
- [x] 终端输入转发到后端
- [x] 终端字体大小调节
- [x] 从 Host 列表快速连接

---

## Phase 3: SFTP 文件传输

- [x] 双面板布局（本地 ↔ 远程）
- [x] 文件树浏览
- [x] 拖拽上传/下载
- [x] 传输进度显示（分块传输 + 事件推送 + 取消传输）
- [x] 本地文件操作（新建文件夹、重命名、删除）— 后端已完成
- [ ] 本地文件操作 UI（新建文件夹、重命名、删除）

---

## Phase 4: 端口转发

### 后端（已完成）
- [x] PortForward CRUD DAO + Service
- [x] Local/Remote/Dynamic 三种模式数据模型
- [x] Toggle 启用/禁用

### 前端（待实现）
- [ ] 端口转发规则 CRUD UI
- [ ] 启用/禁用开关
- [ ] 连接状态监控

---

## Phase 5: Snippets 脚本管理

### 后端（已完成）
- [x] Snippet + SnippetFolder CRUD DAO + Service
- [x] 文件夹分类管理

### 前端（待实现）
- [ ] 代码片段 CRUD UI
- [ ] 快速搜索 + 一键插入终端
- [ ] 多语言语法高亮

---

## Phase 6: 打磨

- [ ] 单元测试 + 集成测试
- [ ] 跨平台构建（Windows/Linux/macOS）
- [ ] 安装包制作（NSIS/AppImage/DMG）
- [ ] 文档完善
