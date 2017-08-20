# conseption
consul inception

This is used to watch a prefix (currently hard coded to /services) in a consul KV store where the values of keys represent consul api AgentServiceRegistration data.  Multiple registrations can exist in the same Value block, with the following formats accepted:

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

Or a json object that has a top level Structures tag

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
