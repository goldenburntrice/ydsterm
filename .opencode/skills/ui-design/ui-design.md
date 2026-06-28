# Terminal 桌面应用开发 Skill（Wails v3 版）

已根据 **Wails v3 (Go + Web 前端)** 调整技术栈、后端集成和通信方式。下面是更新后的 `SKILL.md`。

```markdown
---
name: desktop-terminal-builder
description: 基于 Wails v3 (Go + Web) 构建跨平台桌面终端应用（对标 Termius/Tabby），前端使用 Tailwind CSS。涵盖布局结构、配色系统、xterm.js 集成、Go 后端 PTY/SSH 会话管理、Wails 绑定与事件通信。当用户需要开发终端类工具、SSH 客户端、命令行管理工具时使用。
---

# Desktop Terminal Builder (Wails v3)

## 概述

本 Skill 用于基于 **Wails v3** 构建简洁、专业的跨平台桌面终端应用，参考 Termius、Tabby、Warp 的设计语言。核心原则：**信息密度高但不拥挤、配色克制、终端为主角**。

架构特点：**Go 负责 PTY/SSH 等系统能力，Web 前端负责渲染**，两者通过 Wails 的 Bindings（方法调用）和 Events（流式数据）通信。

## 技术栈选型

| 层级 | 方案 | 说明 |
|------|------|------|
| 桌面框架 | **Wails v3** | Go 后端 + 系统 WebView，体积小 |
| 后端语言 | **Go** | 处理 PTY、SSH、文件传输、凭据存储 |
| 前端框架 | React / Vue / Svelte | 任选，下方以通用写法示例 |
| 样式 | **Tailwind CSS** | 配合 CSS 变量做 design tokens |
| 终端渲染 | **xterm.js** + addons | WebGL 渲染 |
| 本地 PTY | `github.com/creack/pty` (Unix) / `github.com/UserExistsError/conpty` (Windows) | 本地 shell |
| SSH | `golang.org/x/crypto/ssh` | 远程会话 |
| 凭据存储 | `github.com/zalando/go-keyring` | 系统钥匙串 |
| 状态管理 | Zustand / Pinia | 前端轻量状态 |
| 图标 | Lucide / Tabler Icons | 线性图标 |

> **跨平台 PTY 注意**：Windows 需走 ConPTY，Unix 走 creack/pty。建议在 Go 侧封装统一的 `PtySession` 接口，用 build tag 分文件实现。

## Wails v3 项目结构

```
myterm/
├── main.go                 # App 入口，注册 Services
├── wails.json
├── internal/
│   ├── terminal/           # PTY 会话服务
│   │   ├── service.go      # TerminalService（暴露给前端）
│   │   ├── pty_unix.go     // +build !windows
│   │   └── pty_windows.go  // +build windows
│   ├── ssh/                # SSH 会话服务
│   │   └── service.go
│   └── store/              # Host/凭据持久化
│       └── service.go
└── frontend/
    ├── src/
    │   ├── components/      # TitleBar, Sidebar, TerminalPane...
    │   ├── stores/
    │   └── styles/
    ├── tailwind.config.js
    └── ...
```

## Wails v3 通信模型（关键）

Wails v3 中，Go 的 **Service** 方法会自动生成 TS 绑定，前端可直接 `import` 调用；后端推送数据用 **Events**。

### Go 侧：注册 Service

```go
// main.go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "myterm/internal/terminal"
)

func main() {
    app := application.New(application.Options{
        Name: "MyTerm",
        Services: []application.Service{
            application.NewService(terminal.NewService()),
        },
    })

    app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
        Title:            "MyTerm",
        Frameless:        true, // 自绘 TitleBar
        BackgroundColour: application.NewRGB(13, 17, 23), // #0d1117
        Width:            1200,
        Height:           760,
        MinWidth:         800,
        MinHeight:        500,
    })

    app.Run()
}
```

### Go 侧：终端服务（PTY + 事件推流）

```go
// internal/terminal/service.go
package terminal

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct {
    sessions map[string]*PtySession
}

func NewService() *Service {
    return &Service{sessions: map[string]*PtySession{}}
}

// 前端调用：创建本地 shell 会话
func (s *Service) StartLocal(id string, cols, rows int) error {
    sess, err := newPty(cols, rows) // 平台相关实现
    if err != nil {
        return err
    }
    s.sessions[id] = sess

    // 读取 PTY 输出，通过事件推给前端
    go func() {
        buf := make([]byte, 4096)
        for {
            n, err := sess.Read(buf)
            if err != nil {
                application.Get().Event.Emit("term:exit:"+id, nil)
                return
            }
            // 二进制安全：转 string 或 base64
            application.Get().Event.Emit("term:data:"+id, string(buf[:n]))
        }
    }()
    return nil
}

// 前端调用：写入用户输入
func (s *Service) Write(id, data string) error {
    if sess, ok := s.sessions[id]; ok {
        _, err := sess.Write([]byte(data))
        return err
    }
    return nil
}

// 前端调用：窗口/容器 resize 同步
func (s *Service) Resize(id string, cols, rows int) error {
    if sess, ok := s.sessions[id]; ok {
        return sess.Resize(cols, rows)
    }
    return nil
}

func (s *Service) Close(id string) error {
    if sess, ok := s.sessions[id]; ok {
        sess.Close()
        delete(s.sessions, id)
    }
    return nil
}
```

### 前端：调用绑定 + 监听事件

```ts
// Wails v3 自动生成的绑定
import { TerminalService } from '../bindings/myterm/internal/terminal';
import { Events } from '@wailsio/runtime';

async function attachTerminal(term: Terminal, id: string) {
  // 启动后端 PTY
  await TerminalService.StartLocal(id, term.cols, term.rows);

  // 后端输出 → 写入 xterm
  Events.On(`term:data:${id}`, (e) => term.write(e.data));
  Events.On(`term:exit:${id}`, () => term.write('\r\n[process exited]\r\n'));

  // 用户输入 → 发给后端
  term.onData((data) => TerminalService.Write(id, data));

  // resize 同步
  term.onResize(({ cols, rows }) => TerminalService.Resize(id, cols, rows));
}
```

## 设计系统（核心）

### 配色：暗色优先的双主题

终端工具应**默认暗色**。使用语义化 token，不要硬编码颜色。

```css
:root[data-theme="dark"] {
  --bg-base:      #0d1117;  /* 最底层背景 */
  --bg-surface:   #161b22;  /* 侧栏/面板 */
  --bg-elevated:  #1c2128;  /* 卡片/悬浮 */
  --bg-hover:     #21262d;  /* hover 态 */
  --border:       #30363d;
  --border-muted: #21262d;

  --text-primary:   #e6edf3;
  --text-secondary: #8b949e;
  --text-muted:     #6e7681;

  --accent:       #2f81f7;
  --accent-hover: #4493f8;
  --success:      #3fb950;
  --warning:      #d29922;
  --danger:       #f85149;
}

:root[data-theme="light"] {
  --bg-base:      #ffffff;
  --bg-surface:   #f6f8fa;
  --bg-elevated:  #ffffff;
  --bg-hover:     #eaeef2;
  --border:       #d0d7de;
  --text-primary: #1f2328;
  --text-secondary:#656d76;
  --accent:       #0969da;
}
```

> Wails 窗口的 `BackgroundColour` 应与 `--bg-base` 一致，避免启动闪白。

```js
// tailwind.config.js
export default {
  theme: {
    extend: {
      colors: {
        base:     'var(--bg-base)',
        surface:  'var(--bg-surface)',
        elevated: 'var(--bg-elevated)',
        hover:    'var(--bg-hover)',
        border:   'var(--border)',
        primary:  'var(--text-primary)',
        secondary:'var(--text-secondary)',
        muted:    'var(--text-muted)',
        accent:   'var(--accent)',
      },
    },
  },
}
```

### 字体

- **UI 文字**：`Inter` / 系统 UI 字体，13-14px
- **终端**：`JetBrains Mono` / `Cascadia Code`，14px，启用连字
- 行高紧凑：UI `leading-tight`，终端 `1.4`

### 间距与圆角

- 基础单位 4px，密集布局多用 `p-2`、`gap-1.5`、`px-3`
- 圆角统一 `rounded-md`（6px），大面板 `rounded-lg`
- 边框 1px 用 `--border`，避免重阴影；悬浮用 subtle shadow

## 布局结构（黄金三栏）

```
┌─────────────────────────────────────────────────────────┐
│  TitleBar (frameless 自绘，含窗口控制 + Tab 栏)           │  40px
├──────────┬──────────────────────────────────────────────┤
│ Sidebar  │            Terminal Area                      │
│ - Hosts  │       (xterm.js, 可分屏 split panes)           │
│ - Groups │  ┌──────────────┬──────────────┐              │
│ - Snippets│  │  pane 1      │  pane 2      │              │
│  240px   │  └──────────────┴──────────────┘              │
├──────────┴──────────────────────────────────────────────┤
│  StatusBar (连接状态、延迟、编码)                          │  24px
└─────────────────────────────────────────────────────────┘
```

### Frameless 窗口与拖拽

Wails v3 frameless 模式下，需用 CSS 指定可拖拽区域：

```css
/* TitleBar 容器 */
.titlebar { --wails-draggable: drag; }
/* 按钮、Tab 等可点击元素要排除 */
.titlebar button, .titlebar .tab { --wails-draggable: no-drag; }
```

> 窗口控制按钮（最小化/最大化/关闭）调用 Wails runtime 的 `Window.Minimise()` / `Maximise()` / `Close()`。macOS 红绿灯在左，Windows 在右，需按平台调整布局。

### 关键组件清单

1. **自定义 TitleBar** — frameless + 自绘窗口按钮 + 内嵌 Tab
2. **Sidebar** — 搜索框（`Cmd/Ctrl+K` 命令面板）、Host 树形分组、可折叠
3. **Terminal Area** — Tab 切换 + 水平/垂直分屏
4. **StatusBar** — 左侧连接信息，右侧延迟/编码/行列数

## xterm.js 集成要点

```ts
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebglAddon } from '@xterm/addon-webgl';
import { WebLinksAddon } from '@xterm/addon-web-links';

const term = new Terminal({
  fontFamily: '"JetBrains Mono", monospace',
  fontSize: 14,
  lineHeight: 1.4,
  cursorBlink: true,
  cursorStyle: 'bar',
  allowProposedApi: true,
  theme: {
    background: '#0d1117',
    foreground: '#e6edf3',
    cursor: '#2f81f7',
    selectionBackground: '#2f81f740',
    black: '#484f58',  red: '#ff7b72',
    green: '#3fb950',  yellow: '#d29922',
    blue: '#58a6ff',   magenta: '#bc8cff',
    cyan: '#39c5cf',   white: '#b1bac4',
  },
});
const fit = new FitAddon();
term.loadAddon(fit);
term.loadAddon(new WebglAddon());
term.loadAddon(new WebLinksAddon());

// 容器/窗口 resize 时 fit，并把新尺寸同步给 Go 后端
new ResizeObserver(() => fit.fit()).observe(containerEl);
```

**必做**：resize → `fit.fit()` → `term.onResize` 触发 → 调 `TerminalService.Resize` 同步 PTY 尺寸。

## 核心功能模块（按优先级）

1. **会话管理** — 创建/编辑/删除 Host，SSH key / 密码 / 跳板机（Go store service）
2. **多 Tab + 分屏** — 每个 pane 一个会话 id
3. **命令面板** — `Cmd+K` 模糊搜索 host、命令、设置
4. **Snippets** — 常用命令片段，一键发送（前端调 `Write`）
5. **凭据管理** — Go `go-keyring` 加密存储
6. **主题切换 + 字体设置**
7. **SFTP 文件传输**（Go `sftp` 包，进阶）

## 交互细节（提升专业感）

- 快捷键：`Cmd+T` 新 Tab、`Cmd+W` 关闭、`Cmd+D` 分屏、`Cmd+1~9` 切 Tab
- hover/active 用 `--bg-hover`，过渡 `transition-colors duration-150`
- 不用花哨动画，工具类追求即时响应
- 连接状态小圆点：绿(在线)/灰(离线)/黄(连接中，`animate-pulse`)
- 滚动条细化：`scrollbar-thin scrollbar-thumb-border`

## 开发流程建议

1. `wails3 init` 创建项目，配置 frontend 用所选框架 + Tailwind
2. 实现 frameless 窗口 + 自绘 TitleBar + 窗口控制
3. 落地 design tokens + Tailwind 配置（先暗色主题）
4. 实现三栏布局静态结构
5. Go 侧封装跨平台 `PtySession`（Unix/Windows 分文件）+ `TerminalService`
6. 前端接入 xterm.js，通过 Bindings + Events 跑通本地 shell
7. 新增 `SSHService`（`x/crypto/ssh`）实现远程会话
8. 补齐 Tab、分屏、会话管理、命令面板
9. 最后做设置、主题、SFTP

## 参考产品

- **Termius** — 会话管理、Snippets 标杆
- **Tabby** — 开源、Web 技术栈，可读源码
- **Warp** — 现代 UI、命令面板交互
- **Ghostty / iTerm2** — 终端渲染体验基准

## 验收标准

- [ ] frameless 窗口拖拽正常，三平台窗口控制按钮位置正确
- [ ] Wails 窗口背景色与主题一致，启动无闪白
- [ ] 暗/亮主题切换无闪烁，对比度达 WCAG AA
- [ ] 窗口缩放时终端正确 reflow，PTY 尺寸同步无错位
- [ ] PTY 输出二进制安全（中文/特殊字符不乱码，注意 UTF-8）
- [ ] Windows ConPTY 与 Unix PTY 行为一致
- [ ] 终端 WebGL 渲染流畅，大量输出不卡顿
- [ ] 键盘可完成全部核心操作
```

