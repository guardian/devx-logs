import argparse
import json
import socket
import time
import msgpack

parser = argparse.ArgumentParser()
parser.add_argument("message", help="JSON message to send")
args = parser.parse_args()

message = [
    "application-logs",
    [[
        time.time(),
        {
            "MESSAGE": args.message
        }
    ]]
]

data = msgpack.packb(message, use_bin_type=True)

with socket.create_connection(("127.0.0.1", 24224)) as s:
    s.sendall(data)
