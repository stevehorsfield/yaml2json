# yaml2json

`yaml2json` provides a one-way translation from arbitrary (but
JSON-compatible) YAML format data to an equivalent JSON
structure. The motivation is to support translation from
Kubernetes manifests to a form that can be easily queried from
scripts, but may be used wherever YAML is used to represent
a true JSON data structure. (Note that YAML supports constructs
that cannot be represented directly in JSON).

## Usage and example

Usage:

```
yaml2json - Convert YAML to equivalent JSON
(c) 2021 Stephen Horsfield

  -c	compact JSON output
  -s	process multiple YAML inputs as a JSON array
```

Compact output is particularly useful in scripts that need
one-line outputs.

Pipe data to `yaml2json` and then, recommended, to `jq` for
pretty-printing, manipulation or analysis. Here's an example
of raw output:

```
$ yaml2json <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: coredns
  namespace: kube-system
data:
  Corefile: |
    .:53 {
        errors
        log
        health
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
          pods insecure
          fallthrough in-addr.arpa ip6.arpa
        }
        forward . 8.8.8.8:53
        cache 30
        loop
        reload
        loadbalance
    }
EOF
```

Output:

```
{
  "apiVersion": "v1",
  "data": {
    "Corefile": ".:53 {\n    errors\n    log\n    health\n    ready\n    kubernetes cluster.local in-addr.arpa ip6.arpa {\n      pods insecure\n      fallthrough in-addr.arpa ip6.arpa\n    }\n    forward . 8.8.8.8:53\n    cache 30\n    loop\n    reload\n    loadbalance\n}\n"
  },
  "kind": "ConfigMap",
  "metadata": {
    "name": "coredns",
    "namespace": "kube-system"
  }
}
```

## Multi-document files

YAML files, particularly in Kubernetes, can contain multiple
documents split by lines with three consecutive hyphens:

```
---
```

`yaml2json` provides native support for multi-document files. 
Use the `-s` option to handle this. By default, only the first
document is processed. This default mode passes the content
verbatim to the `yaml` library and is less likely to encounter
problems. The `-s` option requires `yaml2json` to operate a
basic state machine over the input to split documents at
appropriate locations.