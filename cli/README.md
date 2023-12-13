# DevX Logs CLI

A small tool to deep-link to Central ELK.

## Installation via homebrew

```bash
brew tap guardian/homebrew-devtools
brew install guardian/devtools/devx-logs

# update
brew upgrade devx-logs
```

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

## Releasing

Releasing is semi-automated. To release a new version, create a new tag with the
`cli-v` prefix:

```bash
git tag cli-v0.0.1
```

And then push the tag:

```bash
git push --tags
```

This will trigger [a GitHub Action](../.github/workflows/release-cli.yml),
publishing a new version to GitHub releases.

Once a new release is available, update the
[Homebrew formula](https://github.com/guardian/homebrew-devtools/blob/main/Formula/devx-logs.rb)
to point to the new version.
