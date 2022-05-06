# devx-logs

A prototype tool to help forward logs of a Systemd Unit to Kinesis. The output
is a [Fluentbit](https://docs.fluentbit.io/manual/) config written to
`/etc/td-agent-bit/td-agent-bit.conf`. Use the `dry-run` flag to write to stdout
instead.

The `cdk-base` Amigo role uses this to provide log shipping out of the box.

The following logs are supported:

- [x] cloud-init
- [ ] application logs (currently being beta tested with DevX applications - we hope to add full support for this soon!)

Use the `-h` flag for more info.

```
devx-logs outputs a Fluentbit config appropriate for Guardian EC2 applications.

Configuration is typically provided by tags on the instance, but flags are also supported to customise behaviour.

Usage:
  devx-logs [flags]

Flags:
      --dry-run                    Set to true to print config to stdout rather than write to file.
  -h, --help                       help for devx-logs
      --kinesisStreamName string   Typically configured via a 'LogKinesisStreamName' tag on the instance, but you can override using this flag. To write to Kinesis, your instance will need the following permissions for this stream: kinesis:DescribeStream, kinesis:PutRecord.
      --tags string                Typically read from /etc/config/tags.json (see Amigo's cdk-base role here for more info), but you can override using this flag. Pass a comma-separated list of Key=Value pairs, to be included on log records.
```

# Deployment

This tools deploys via Riffraff whenever a PR is merged into `main`.
