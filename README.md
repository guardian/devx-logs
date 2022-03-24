# devx-logs

A prototype tool to help forward logs of a SystemD Unit to Kinesis. The output
is a Fluentbit config, that reads from JournalD, and filters to the specified
SystemD Unit.

Instance app/stage/stack tags are added to records. You can add additional tags
via a tag(!): `LogTags`.

Use the `-h` flag for more info.

By default, devx-logs ships the following:

* /var/log/cloud-init-output.log
* the specified systemd unit(s)
* ~memory statistics (see
  [here](https://docs.fluentbit.io/manual/pipeline/inputs/memory-metrics) for
  details)~
* it's own journald logs(!) and also fluentbit's logs

## Instace Metadata Service Version 2 (IMDSV2)

`devx-logs` requires instances to be running v2 of the Instance Metadata Service
with tags enabled. To enable this set MetadataOptions in your Launch Template:

https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata.html#cfn-ec2-launchtemplate-launchtemplatedata-metadataoptions



