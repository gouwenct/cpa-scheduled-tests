# CPA Scheduled Tests

Native CLIProxyAPI (CPA) plugin for Sub2API-style scheduled Codex tests.

> 中文说明：[README_CN.md](README_CN.md)

## Features

- Per-account scheduled tests: **Account + Model + Cron + Timezone + Prompt**
- **Bulk-create plans for all Codex accounts** with duplicate detection
- Per-plan **Run Now**
- Quick send to one account/model
- **Send to all accounts now**, including records currently marked disabled/unavailable
- Model suggestions synchronized from CPA Codex model definitions, while still allowing manual model input
- Persistent JSONL execution logs with trigger, account, model, HTTP status, latency and error category
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

For Windows source builds and one-click local installation, see [README_CN.md](README_CN.md).

## Example cron

```text
0 6,11,16,21 * * *
```

Runs every day at 06:00, 11:00, 16:00 and 21:00 in the plan's configured timezone.

## Safety / behavior

`Send to all accounts now` performs real model requests. It intentionally attempts every Codex auth record returned by CPA, including records marked disabled/unavailable; invalid credentials are logged as failures. The plugin never changes CPA account enable/disable state and does not persist OAuth tokens or full upstream responses.

## License

MIT
