# CPA Scheduled 5H

<img src="assets/logo.png" alt="5H OPEN" width="128">

English version: [README_EN.md](README_EN.md)

## 界面截图

![Scheduled Tests 插件界面](assets/scheduled-tests.png)

截图来自早期界面；当前版本另有“一键添加全部账号”和 GitHub / Star 入口。

一个原生 CLIProxyAPI（CPA）插件，按 Sub2API Scheduled Tests 的思路提供：

- **账号 + Model + Cron** 的定时测试计划：每个 Cron 时刻在设定分钟和下一分钟各发送一次，无论第一次结果如何；
- **一键为全部 Codex 账号创建定时计划**：统一指定 Model、Cron、时区和 Prompt，自动一账号一计划，并跳过完全重复计划；
- 每个计划的 **立即发送、暂停/恢复**，计划表每页显示 10 条；
- 计划按 **分组** 管理，编辑和新建计划提供账号、Model 的下拉选择及默认值，Cron 保留自由输入；
- 时区支持常用建议和任意有效 IANA 时区输入；新建计划默认使用浏览器时区，编辑时保留计划已保存的时区；
- **一键向全部 Codex 账号发送**，遍历 `host.auth.list` 中的所有 Codex auth 记录，不因 `disabled` / `unavailable` 字段而跳过；
- JSONL 持久化日志：时间、触发方式、账号、Model、HTTP 状态、5 小时窗口状态、使用比例、重置时间、延迟和错误分类；界面只显示最近两天；
- 全部功能运行在 CPA 进程内，无需 Windows Task Scheduler、PowerShell 常驻或外部服务。

## 重要说明

“全部账号立即发送”是**强制真实请求**，默认 Prompt 为“你好”。发送后会读取响应头中的五小时窗口数据；缺少明确数据时，使用同一账号凭据查询上游额度接口。日志显示发送后的窗口快照：确认 300 分钟窗口已有用量且重置时间未到，为“已开启”；用量达到 100% 时为“已开启 · 额度耗尽”；缺少数据或旧日志为“未知”。HTTP 200 本身不作为窗口开启的证据。该操作会产生极小的真实模型请求消耗。对于 `disabled` / `unavailable` 的 auth 文件，插件仍会尝试读取凭据并请求；若 token 已失效，会记录 `auth_error`，不会偷偷修改 CPA 账号状态。

当前 v0.1.6 专门面向 **Codex / ChatGPT OAuth auth**。UI 会通过 CPA 的 `model-definitions/codex` 自动同步当前模型列表作为下拉选项。

## UI

安装后在 CPA Management Center 侧栏打开：

**CPA Scheduled 5H**

或者直接访问：

```text
http://127.0.0.1:8317/v0/resource/plugins/cpa-scheduled-tests/panel
```

资源页面本身不携带敏感数据。首次打开请输入 CPA Management Key；它只保存到当前浏览器标签页的 `sessionStorage`。

## Cron

只需填写一条 Cron。例如 `0 6,11,16,21 * * *` 会在每天 06:00/06:01、11:00/11:01、16:00/16:01、21:00/21:01 各发送一次。两次结果分别写入日志；第二次标为 `scheduled+1min`。跨小时、跨天仍按一分钟间隔处理。重复发送提高五小时窗口开启成功率，但不能保证绕过上游限流。手动立即发送仍只发送一次。

使用标准 5 段数字 Cron：

```text
分钟 小时 日 月 星期
```

示例：

```text
0 6,11,16,21 * * *
```

代表每天 06:00、11:00、16:00、21:00。

支持：

- `*`
- `*/5`
- `1,2,3`
- `1-5`
- `1-10/2`

星期使用 `0-7`，其中 `0` 和 `7` 都代表星期日。

## 安装配置

CPA `config.yaml`：

```yaml
plugins:
  enabled: true
  dir: plugins
  configs:
    cpa-scheduled-tests:
      enabled: true
      priority: 1
```

各平台动态库放到对应目录（`amd64` 表示 Intel/AMD 64 位，`arm64` 包括 Apple Silicon）：

| 平台 | 架构 | 安装位置 |
| --- | --- | --- |
| Windows | amd64 | `plugins/windows/amd64/cpa-scheduled-tests.dll` |
| Linux | amd64 | `plugins/linux/amd64/cpa-scheduled-tests.so` |
| Linux | arm64 | `plugins/linux/arm64/cpa-scheduled-tests.so` |
| macOS | amd64 | `plugins/darwin/amd64/cpa-scheduled-tests.dylib` |
| macOS | arm64 | `plugins/darwin/arm64/cpa-scheduled-tests.dylib` |

Windows 示例：

```text
plugins/windows/amd64/cpa-scheduled-tests.dll
```

也可以直接放：

```text
plugins/cpa-scheduled-tests.dll
```

然后完整重启 CPA。

## 让 Agent 安装或更新

把下面的提示词直接发给你的 Agent。它会优先使用中文，并在替换 DLL 前检查版本、备份旧文件和校验下载包：

```text
请使用中文帮助我安装或更新 CPA Scheduled 5H 插件到最新版本。

目标：
- 仓库：https://github.com/gouwenct/cpa-scheduled-tests
- 先读取 GitHub 最新 Release 和 checksums.txt，不要猜版本号。
- 先检查本机 EasyCLIProxyAPI/CPA Core 是否运行、插件目录、当前 DLL 版本和文件哈希。
- 如果 CPA Core 正在运行，请先提示我在 EasyCLIProxyAPI 中停止 CPA Core；不要强制结束进程。
- 停止后，把旧的 cpa-scheduled-tests.dll 备份到同目录的 OLD 文件夹，再复制新 DLL。
- 校验下载包中的 SHA-256；校验失败就停止，不要覆盖旧文件。
- 保持文件名为 cpa-scheduled-tests.dll，不要修改现有计划、日志或 config.yaml 中无关的配置。
- 启动 CPA Core 后验证插件的 registered=true、effective_enabled=true，并报告安装路径、版本和哈希。
- 不要输出、保存或传播 CPA Management Key、OAuth token 或账号凭据。
- 没有权限时，明确告诉我需要以管理员身份执行的具体步骤。

完成后请用中文报告：已执行操作、验证结果、注意事项、未执行操作。
```

安装包请从 [最新 Release](https://github.com/gouwenct/cpa-scheduled-tests/releases/latest) 下载与系统、架构一致的 ZIP，并使用同一 Release 的 `checksums.txt` 校验 SHA-256。

## macOS / Linux 构建与完整平台发布

原生构建需要 Go 1.23+、Python 3 和 C 编译器。macOS 使用 Xcode Command Line Tools；Linux 使用 GCC。运行：

```sh
bash build-linux.sh
```

该脚本也支持 macOS，按当前系统和架构输出 `.so` 或 `.dylib`、平台 ZIP 与 `checksums.txt`；会使用模拟 CPA 宿主验证动态库，不读取真实账号凭据。

[GitHub Actions 构建流程](.github/workflows/build-release.yml) 在五个平台分别运行 Go 测试、`go vet` 和 ABI 模拟加载，汇总后生成含五个 ZIP 与完整校验清单的 `release-bundle`。

普通推送只生成构建产物。推送与源码版本一致的 `v<版本号>` 标签后，五个平台全部通过且安装包校验成功，才会发布 GitHub Release。


## Windows 一键构建并安装

源码包附带 `install-windows.ps1`。它不会常驻，也不负责定时任务；它只在第一次安装/升级时下载便携 Go + Zig、编译 CPA 原生 DLL，并复制到 CPA 插件目录。插件安装完成后所有调度都在 CPA 进程内运行。

例如 CPA 在 `D:\Program Files\CLIProxyAPI`：

```powershell
Set-ExecutionPolicy -Scope Process Bypass
.\install-windows.ps1 -CpaDir 'D:\Program Files\CLIProxyAPI'
```

如果同时提供 Management Key，脚本会尝试通过 CPA Management API 自动创建/启用插件配置：

```powershell
.\install-windows.ps1 `
  -CpaDir 'D:\Program Files\CLIProxyAPI' `
  -BaseUrl 'http://127.0.0.1:8317' `
  -ManagementKey '你的ManagementKey'
```

安装后完整重启一次 CPA，然后在侧栏进入 **CPA Scheduled 5H**。

## 数据位置

插件优先写入 CPA 提供的插件目录下 `.data/cpa-scheduled-tests/`；如果该目录不可写（例如受保护的 Program Files），自动回退到用户配置目录中的 `CLIProxyAPI/plugin-data/cpa-scheduled-tests/`。具体实际路径会显示在插件状态数据中。

主要文件：

```text
state.json
runs.jsonl
```

其中：

- `state.json`：计划和 UI 设置；
- `runs.jsonl`：持久化运行日志。

不会写入 OAuth token、Access Token 或完整上游响应。

## Management API

均位于 `/v0/management` 下，需要 CPA Management Key。

```text
GET  /plugins/cpa-scheduled-tests/state
POST /plugins/cpa-scheduled-tests/plans/save
POST /plugins/cpa-scheduled-tests/plans/bulk-create
POST /plugins/cpa-scheduled-tests/plans/delete
POST /plugins/cpa-scheduled-tests/settings
POST /plugins/cpa-scheduled-tests/run
POST /plugins/cpa-scheduled-tests/run-all
GET  /plugins/cpa-scheduled-tests/logs?limit=200
POST /plugins/cpa-scheduled-tests/logs/clear
```

## Windows 构建

此插件只依赖 Go 标准库，但 CPA 原生动态库需要 CGO。

在 Windows PowerShell 中：

```powershell
.\build-windows.ps1
```

要求：

- Go 1.23+；
- 一个可用的 C 编译器（推荐 MSYS2 UCRT64 / MinGW-w64）。

输出：

```text
dist\cpa-scheduled-tests.dll
dist\cpa-scheduled-tests_0.1.6_windows_amd64.zip
```

## 安全设计

- 敏感动作全部放在 CPA Management API 路由后面；
- Resource UI 不直接暴露 auth JSON；
- 不记录 token、原始 auth 文件或完整请求体；
- `host.http.do` 由 CPA 宿主执行出站请求，从而复用 CPA 的宿主网络/代理策略；
- 插件不会自动启用、禁用、恢复或修改任何 CPA auth 文件。
