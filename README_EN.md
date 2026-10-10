# CPA Scheduled 5H

<img src="assets/logo.png" alt="5H OPEN" width="128">

Native CLIProxyAPI (CPA) plugin for Sub2API-style scheduled Codex tests.

> 中文说明：[README.md](README.md)

## Screenshot

![Scheduled Tests panel](assets/scheduled-tests.png)

The screenshot shows an earlier UI; current releases also include bulk plan creation and a GitHub / Star link.

## Features

- Per-account scheduled tests: **Account + Model + Cron + Timezone + Prompt**. Each Cron occurrence sends once in the configured minute and again in the following minute, regardless of the first result. Both results are logged; the second is marked `scheduled+1min`. Manual sends remain single requests. This improves activation reliability but cannot bypass upstream quota limits.
- **Bulk-create plans for all Codex accounts** with duplicate detection
- Per-plan **Run Now** and pause/resume controls, with 10 plans per page
- Plans are organized by **group**, with dropdown defaults for account and model; Cron remains free-form
- Timezones accept suggestions or any valid IANA timezone. New plans default to the browser timezone; editing preserves the saved timezone.
- **Send to all accounts now**, including records currently marked disabled/unavailable; the default prompt is `你好`
- Model options synchronized from CPA Codex model definitions
- Persistent JSONL execution logs with trigger, account, model, HTTP status, five-hour window status, usage percentage, reset time, latency and error category; the UI keeps only the latest two days
- Runs fully inside the CPA process; no Windows Task Scheduler or resident PowerShell/Python/Node service is required

The five-hour status is a snapshot after sending. The plugin reads the 300-minute quota window from response headers, with a per-account upstream usage query as fallback. Positive usage and a future reset time confirm an open window; 100% usage means the open window is exhausted. Missing evidence and old logs show Unknown. HTTP 200 alone does not confirm an open window.

## Install

Enable the plugin in CPA `config.yaml`:

```yaml
plugins:
  enabled: true
  dir: plugins
  configs:
    cpa-scheduled-tests:
      enabled: true
      priority: 1
```

On Windows x64, place the DLL at:

```text
plugins/windows/amd64/cpa-scheduled-tests.dll
```

For other platforms, use the matching ZIP and install the root library at:

| Platform | Architecture | Path |
| --- | --- | --- |
| Linux | amd64 | `plugins/linux/amd64/cpa-scheduled-tests.so` |
| Linux | arm64 | `plugins/linux/arm64/cpa-scheduled-tests.so` |
| macOS | amd64 | `plugins/darwin/amd64/cpa-scheduled-tests.dylib` |
| macOS | arm64 | `plugins/darwin/arm64/cpa-scheduled-tests.dylib` |

Download from the [latest release](https://github.com/gouwenct/cpa-scheduled-tests/releases/latest) and verify the ZIP against that release's `checksums.txt`.

Then fully restart CPA and open **CPA Scheduled 5H** in the Management Center.

## Agent-assisted installation

For a Chinese-first installation or update, copy the prompt in [README.md](README.md#让-agent-安装或更新) to your Agent. It checks the latest release, verifies the SHA-256 checksum, backs up the old DLL, replaces it only after CPA Core is stopped, and verifies plugin registration after restart. It must not expose management keys or OAuth tokens.

For Windows source builds and one-click local installation, see [README.md](README.md).

## Example cron

```text
0 6,11,16,21 * * *
```

Runs every day at 06:00, 11:00, 16:00 and 21:00 in the plan's configured timezone.

## Safety / behavior

`Send to all accounts now` performs real model requests. It intentionally attempts every Codex auth record returned by CPA, including records marked disabled/unavailable; invalid credentials are logged as failures. The plugin never changes CPA account enable/disable state and does not persist OAuth tokens or full upstream responses.

## License

MIT

## Build and release

On Linux or macOS, install Go 1.23+, Python 3 and a C compiler (GCC or Xcode Command Line Tools), then run `bash build-linux.sh`. It builds the native shared library, runs the fake-host ABI smoke test, and writes the platform ZIP and SHA-256 checksum. Windows can use `python build-release.py` with Go and GCC on PATH.

The [platform build workflow](.github/workflows/build-release.yml) runs tests, vet, native library builds and ABI smoke tests for Darwin amd64/arm64, Linux amd64/arm64 and Windows amd64. Ordinary pushes produce a `release-bundle` artifact. A matching version tag publishes the complete release only after all five platforms pass.
