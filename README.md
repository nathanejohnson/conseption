# conseption
consul inception

This is used to watch a prefix in a consul KV store where the values of keys represent consul api AgentServiceRegistration data.  This is useful for services that are not consul-aware but need to be registered as services within consul.  The way it works currently, a consul agent as well as this program would be deployed on each machine where registrations need to occur, and it would only register services to the local agent that match the hostname of where we live.  This can be overridden by changing the HOSTNAME environment variable.  CONSUL_HTTP_ADDR must be set to point to a consul agent, or there must be an agent at localhost, or the config must be specified in a configuration toml file.  If it notices health checks registered to another agent, it will deregister first.

Multiple registrations can exist in the same Value block in the consul KV, with the following formats accepted:

Serial representations of [AgentServiceRegistration](https://godoc.org/github.com/hashicorp/consul/api#AgentServiceRegistration)

```json
{
    "id": "cb01.labs.widget.co",
    "name": "couchbase",
    "tags": [
        "cache",
        "cloudstack",
        "rancid",
        "cloudythings"
    ],
    "address": "cb01.labs.widget.co",
    "port": 8091,
    "checks": [
        {
            "http": "http://cb01.labs.widget.co:8091/pools/",
            "interval": "30s"
        }
    ]
}
{
    "id": "cb02.labs.widget.co",
    "name": "couchbase",
    "tags": [
        "cache",
        "cloudstack",
        "rancid",
        "cloudythings"
    ],
    "address": "cb02.labs.widget.co",
    "port": 8091,
    "checks": [
        {
            "http": "http://cb02.labs.widget.co:8091/pools/",
            "interval": "30s"
        }
    ]
}
```

Comma separated representations of [AgentServiceRegistration](https://godoc.org/github.com/hashicorp/consul/api#AgentServiceRegistration)

```json
{
    "id": "cb01.labs.widget.co",
    "name": "couchbase",
    "tags": [
        "cache",
        "cloudstack",
        "rancid",
        "cloudythings"
    ],
    "address": "cb01.labs.widget.co",
    "port": 8091,
    "checks": [
        {
            "http": "http://cb01.labs.widget.co:8091/pools/",
            "interval": "30s"
        }
    ]
},
{
    "id": "cb02.labs.widget.co",
    "name": "couchbase",
    "tags": [
        "cache",
        "cloudstack",
        "rancid",
        "cloudythings"
    ],
    "address": "cb02.labs.widget.co",
    "port": 8091,
    "checks": [
        {
            "http": "http://cb02.labs.widget.co:8091/pools/",
            "interval": "30s"
        }
    ]
}
```

A list of [AgentServiceRegistration](https://godoc.org/github.com/hashicorp/consul/api#AgentServiceRegistration)s
```json
[
    {
        "id": "cb01.labs.widget.co",
        "name": "couchbase",
        "tags": [
            "cache",
            "cloudstack",
            "rancid",
            "cloudythings"
        ],
        "address": "cb01.labs.widget.co",
        "port": 8091,
        "checks": [
            {
                "http": "http://cb01.labs.widget.co:8091/pools/",
                "interval": "30s"
            }
        ]
    },
    {
        "id": "cb02.labs.widget.co",
        "name": "couchbase",
        "tags": [
            "cache",
            "cloudstack",
            "rancid",
            "cloudythings"
        ],
        "address": "cb02.labs.widget.co",
        "port": 8091,
        "checks": [
            {
                "http": "http://cb02.labs.widget.co:8091/pools/",
                "interval": "30s"
            }
        ]
    }
]
```

Or a json object that has a top level Services field

```json
{ "Services":
    [
        {
            "id": "cb01.labs.widget.co",
            "name": "couchbase",
            "tags": [
                "cache",
                "cloudstack",
                "rancid",
                "cloudythings"
            ],
            "address": "cb01.labs.widget.co",
            "port": 8091,
            "checks": [
                {
                    "http": "http://cb01.labs.widget.co:8091/pools/",
                    "interval": "30s"
                }
            ]
        },
        {
            "id": "cb02.labs.widget.co",
            "name": "couchbase",
            "tags": [
                "cache",
                "cloudstack",
                "rancid",
                "cloudythings"
            ],
            "address": "cb02.labs.widget.co",
            "port": 8091,
            "checks": [
                {
                    "http": "http://cb02.labs.widget.co:8091/pools/",
                    "interval": "30s"
                }
            ]
        }
    ]
}
```
