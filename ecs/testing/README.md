> **Warning**
> This is a work in progress.

# Testing

The scripts in this directory show how the `ecs` project can be tested in a local environment, without needing to 
deploy to AWS, and also show how the process is wired up.

The container listens on port 24224, for messages in the webpack format.

The container can be started locally by running the start-container.sh script from this directory.
It will build the image and run it, and you can then send messages to it using the `send-message.sh` script.

In normal operation, the process sends messages to an AWS Kinesis stream, which is then read by a Lambda function and 
sent to the Central ELK stack.

For testing purposes, add

```
[OUTPUT]
    Name stdout
    Match *
```

to the `guardian.conf` file and the messages will be printed to stdout.  Note that you 
will *also* see a failure to send to AWS, but this part of the process is not relevant to testing the container.