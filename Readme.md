# JSON RPC(?) HTTP Server

## Installation

Fluent(td-agent)

```bash
curl -L https://toolbelt.treasuredata.com/sh/install-redhat-td-agent2.sh | sh

/opt/td-agent/embedded/bin/fluent-gem install fluent-plugin-s3
```

Fluent(td-agent) test

```bash
curl -X POST -d 'json={"json":"message"}' http://localhost:8888/td.log.server
```