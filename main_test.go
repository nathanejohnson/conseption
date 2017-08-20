package main

import (
	"testing"
)

var payload = []byte(`{
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
},
{
    "id": "cb01-dev.davidson.tn.widget.co",
    "name": "couchbase",
    "tags": [
        "throb",
        "radius",
        "satchel",
        "reportingResults",
        "catawampus"
    ],
    "address": "cb01-dev.davidson.tn.widget.co",
    "port": 8091,
    "checks": [
        {
            "http": "http://cb01-dev.davidson.tn.widget.co:8091/pools/",
            "interval": "30s"
        }
    ]
},
{
    "id": "cb02-dev.davidson.tn.widget.co",
    "name": "couchbase",
    "tags": [
        "throb",
        "radius",
        "satchel",
        "reportingResults",
        "catawampus"
    ],
    "address": "cb02-dev.davidson.tn.widget.co",
    "port": 8091,
    "checks": [
        {
            "http": "http://cb02-dev.davidson.tn.widget.co:8091/pools/",
            "interval": "30s"
        }
    ]
},
{
    "id": "cb03-dev.davidson.tn.widget.co",
    "name": "couchbase",
    "tags": [
        "throb",
        "radius",
        "satchel",
        "reportingResults",
        "catawampus"
    ],
    "address": "cb03-dev.davidson.tn.widget.co",
    "port": 8091,
    "checks": [
        {
            "http": "http://cb03-dev.davidson.tn.widget.co:8091/pools/",
            "interval": "30s"
        }
    ]
}`)

var payload2 = []byte(`{
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
{
    "id": "cb01-dev.davidson.tn.widget.co",
    "name": "couchbase",
    "tags": [
        "throb",
        "radius",
        "satchel",
        "reportingResults",
        "catawampus"
    ],
    "address": "cb01-dev.davidson.tn.widget.co",
    "port": 8091,
    "checks": [
        {
            "http": "http://cb01-dev.davidson.tn.widget.co:8091/pools/",
            "interval": "30s"
        }
    ]
}
{
    "id": "cb02-dev.davidson.tn.widget.co",
    "name": "couchbase",
    "tags": [
        "throb",
        "radius",
        "satchel",
        "reportingResults",
        "catawampus"
    ],
    "address": "cb02-dev.davidson.tn.widget.co",
    "port": 8091,
    "checks": [
        {
            "http": "http://cb02-dev.davidson.tn.widget.co:8091/pools/",
            "interval": "30s"
        }
    ]
}
{
    "id": "cb03-dev.davidson.tn.widget.co",
    "name": "couchbase",
    "tags": [
        "throb",
        "radius",
        "satchel",
        "reportingResults",
        "catawampus"
    ],
    "address": "cb03-dev.davidson.tn.widget.co",
    "port": 8091,
    "checks": [
        {
            "http": "http://cb03-dev.davidson.tn.widget.co:8091/pools/",
            "interval": "30s"
        }
    ]
}
`)

func TestProcessor(t *testing.T) {
	for i, p := range [][]byte{
		payload,
		[]byte("[" + string(payload) + "]"),
		[]byte("{\"services\": [" + string(payload) + "]}"),
		payload2,
	} {
		regs, err := parseServiceRegs(p)
		if err != nil {
			t.Errorf("parseServiceRegs failed in pass %d: %s\n", i err)
			t.FailNow()
		}
		if len(regs) != 5 {
			t.Errorf("Invalid length: %d in pass %d\n", len(regs), i)
			t.FailNow()
		}
		if regs[0].Address != "cb01.labs.widget.co" {
			t.Errorf("address mismatch in pass %d", i)
			t.Fail()
		}
		if regs[4].Name != "couchbase" {
			t.Errorf("name mismatch in pass %d", i)
			t.Fail()
		}

	}
}
