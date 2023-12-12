# DevX Logs CLI

A small tool to deep-link to Central ELK.

## Installation

TODO.

## Usage

- Open the logs for Riff-Raff in PROD
  ```bash
  devx-logs --space devx --stage PROD --app riff-raff
  ```
- Display the URL for logs from Riff-Raff in PROD
  ```bash
  devx-logs --space devx --stage PROD --app riff-raff --no-follow
  ```
- Open the logs for Riff-Raff in PROD, where the level is INFO, and show the
  message and logger_name columns
  ```bash
  devx-logs --space devx --stage PROD --app riff-raff --filter level=INFO --filter region=eu-west-1 --column message --column logger_name
  ```
- Open the logs for the repository 'guardian/prism':
  ```bash
  devx-logs --filter gu:repo.keyword=guardian/prism --column message --column gu:repo
  ```

See all options via the `--help` flag:

```bash
devx-logs --help
```
