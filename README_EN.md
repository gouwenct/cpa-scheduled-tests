# CPA Scheduled Tests

Native CLIProxyAPI (CPA) plugin for Sub2API-style scheduled Codex tests.

> 中文说明：[README.md](README.md)

## Screenshot

![Scheduled Tests panel](assets/scheduled-tests.png)

The screenshot shows an earlier UI; current releases also include bulk plan creation and a GitHub / Star link.

## Features

- Per-account scheduled tests: **Account + Model + Cron + Timezone + Prompt**
- **Bulk-create plans for all Codex accounts** with duplicate detection
- Per-plan **Run Now** and pause/resume controls, with 10 plans per page
- Plans are organized by **group**, with dropdown defaults for account and model; Cron remains free-form
- Timezones accept suggestions or any valid IANA timezone. New plans default to the browser timezone; editing preserves the saved timezone.
- **Send to all accounts now**, including records currently marked disabled/unavailable
- Model options synchronized from CPA Codex model definitions
- Persistent JSONL execution logs with trigger, account, model, HTTP status, latency and error category; the UI keeps only the latest two days
- Runs fully inside the CPA process; no Windows Task Scheduler or resident PowerShell/Python/Node service is required

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

Then fully restart CPA and open **Scheduled Tests** in the Management Center.

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
