# devx-logs

A prototype tool to help forward logs of a Systemd Unit to Kinesis. Additional
instance logs (cloud-init-output.log) are also included. The output is a
[Fluentbit](https://docs.fluentbit.io/manual/) config. You should typically
store this under `/etc/td-agent-bit/td-agent-bit.conf` if using Ubuntu.

Use the `-h` flag for more info.

```
devx-logs outputs a Fluentbit config appropriate for Guardian EC2 applications.

Usage:
  devx-logs [flags]

Flags:
  -h, --help                       help for devx-logs
      --kinesisStreamName string   Set to a Kinesis log stream name. Your instance will need the following permissions for this stream: kinesis:DescribeStream, kinesis:PutRecord.
      --systemdUnit string         Set to the name of your app's systemd service. I.e. 'name' from [name].service
      --tags string                Set a comma-separated list of Key=Value pairs, to be included on log records. At the least, this should include App, Stack, and Stage. Eg. 'App=foo,Stage=PROD,Stack=bar'.
```



