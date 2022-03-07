# devx-logs

A prototype tool to help forward logs of a SystemD Unit to Kinesis. The output
is a Fluentbit config, that reads from JournalD, and filters to the specified
SystemD Unit.

Use the `-h` flag for more info.
