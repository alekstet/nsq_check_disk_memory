# Description
This service allows you to monitor the memory filling with NSQ messages, this may be caused by the absence of a client reading the queue, or because reading from the queue is slower than writing to it.
If this fact is detected, a message of the format is thrown into a certain topic:

```
[address, port, topic, channel]
```

For the program to work in the ___config.yml___ file, you need to specify the address of _Nsqlookupd_, the number of messages in memory and on disk, the address and subject where messages from the service will be sent.
You also need to set the _Nsqlookupd_ polling period, the program also sends messages to this channel that the service is working and functioning, the period of these test messages is also set in the configuration file.
