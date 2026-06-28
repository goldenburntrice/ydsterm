# SSH 终端模块 (Terminal)

## 概述

SSH 终端模块负责建立 SSH 连接、管理 PTY 会话、在 xterm.js 终端模拟器中渲染。对应开发计划 Phase 2，状态：**已完成**。

## 后端实现

### 终端服务

`internal/service/ssh.go` — `TerminalServiceImpl`

| 方法 | 签名 | 说明 |
|------|------|------|
| Connect | `(hostID) → sessionID` | 从 DB 加载主机配置 → 建立 SSH 连接 → 分配 PTY → 返回会话 ID |
| Write | `(sessionID, data) → error` | 将前端输入写入 PTY stdin |
| Resize | `(sessionID, cols, rows) → error` | 同步终端窗口大小 |
| Disconnect | `(sessionID) → error` | 关闭 PTY 和 SSH 连接 |
| ListSessions | `() → []map` | 列出当前活跃会话 |

### SSH 连接流程

```
Connect(hostID)
  ├── dao.YdstermHosts.Get(hostID)          → 加载主机配置
  ├── buildSSHConfig(host)                  → 构建 SSH 配置
  │   ├── crypto.Decrypt(passwordEnc)       → 解密密码
  │   └── ssh.ParsePrivateKey(...)          → 解析密钥
  ├── ssh.Dial("tcp", addr, config)         → TCP 连接
  ├── client.NewSession()                    → 创建 SSH 会话
  ├── session.RequestPty("xterm-256color")  → 分配 PTY
  ├── session.StdinPipe()                    → 输入管道
  ├── session.StdoutPipe() / StderrPipe()   → 输出管道
  ├── session.Shell()                        → 启动 Shell
  └── go readPipe(stdout) / go readPipe(stderr)
      └── bufio.Read() → Event.Emit("term-output")
```

### 事件通信

**Go → 前端**（输出流）：

```go
// readPipe goroutine
n, _ := reader.Read(buf)
application.Get().Event.Emit("term-output", TermOutputPayload{
    SessionId: sessionID,
    Data:      base64.StdEncoding.EncodeToString(buf[:n]),
    IsStderr:  false,
})
```

**前端 → Go**（用户输入）：

```ts
// xterm.onData 回调
term.onData((data) => {
    TerminalService.Write(tab.sessionId, data)
})
```

### 事件注册

```go
// main.go init() — 必须在 Emit 前注册
func init() {
    application.RegisterEvent[TermOutputPayload]("term-output")
    application.RegisterEvent[TermDisconnectedPayload]("term-disconnected")
}
```

### Payload 结构

```go
type TermOutputPayload struct {
    SessionId string  // session ID (PascalCase → JS: payload.data.SessionId)
    Data      string  // base64 编码的终端输出
    IsStderr  bool    // 是否 stderr
}

type TermDisconnectedPayload struct {
    SessionId string
}
```

## 前端实现

### 页面

`views/TerminalView.vue` — 多 Tab 终端视图：

**功能**：
- Tab 栏管理（新建/切换/关闭）
- 主机选择下拉菜单
- xterm.js 终端渲染（256 色 + FitAddon + WebLinksAddon）
- 窗口 resize → 终端 reflow → PTY 同步

**数据流**：
```
用户按键 → term.onData → Write(sessionId, data) → Go PTY stdin
Go PTY stdout → Event.Emit("term-output") → Events.On → term.write(data)
```

**事件解包**：
```ts
Events.On('term-output', (payload) => {
    const inner = payload.data                // Wails 信封: {name, data}
    const sid  = inner.SessionId              // Go 字段名 PascalCase
    const dat  = inner.Data
    tab.terminal.write(atob(dat))
})
```

### xterm.js 配置

```ts
const term = new Terminal({
    cursorBlink: true,
    cursorStyle: 'bar',
    fontSize: 14,
    fontFamily: '"JetBrains Mono", "Cascadia Code", Consolas, monospace',
    theme: { /* GitHub-dark 配色 */ },
    allowProposedApi: true,
})
term.loadAddon(new FitAddon())        // 自适应大小
term.loadAddon(new WebLinksAddon())   // 链接识别
```

### WebView 焦点管理

Wails WebView2 中 xterm 的 hidden textarea 需手动聚焦：

```ts
term.open(container)
term.focus()  // 初始聚焦
container.addEventListener('click', () => term.focus())    // 点击重新聚焦
container.addEventListener('mousedown', () => term.focus())
```

### 多 Tab 容器管理

`v-for` 内使用函数 ref 避免只绑定最后一个元素：

```ts
const containerEls: Record<string, HTMLDivElement> = {}
const containerRef = (id: string) => (el: unknown) => {
    if (el) containerEls[id] = el as HTMLDivElement
}

// template
<div v-for="tab in tabs" :ref="containerRef(tab.id)" />
```

## 依赖

| 层级 | 依赖 | 用途 |
|------|------|------|
| Go | `golang.org/x/crypto/ssh` | SSH 协议实现 |
| Go | `github.com/wailsapp/wails/v3` | Event 推送 |
| 前端 | `@xterm/xterm` | 终端模拟器 |
| 前端 | `@xterm/addon-fit` | 自适应容器 |
| 前端 | `@xterm/addon-web-links` | 链接点击 |

## 使用流程

1. 在主机列表双击主机卡片（或终端页点击 + 选择主机）
2. 后端建立 SSH 连接，返回 sessionID
3. 前端创建 Tab + 初始化 xterm.js
4. goroutine 开始读取 SSH 输出并推送到前端
5. 用户在终端中输入 → `Write()` 发送到 PTY
6. 关闭 Tab 或断开连接 → `Disconnect()`

## 扩展点

- 本地终端（PTY，无需 SSH）
- 分屏功能（水平/垂直 split）
- 连接状态状态栏（延迟、编码）
- 终端字体大小热键调节
