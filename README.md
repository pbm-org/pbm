# Pbb Protobuf Build system for multi git repo

pbb for

## How to use

init pbb project

```
pbb init
```

build proto protobuf

```
pbb gen
```

update dep buf latest

```
pbb dep
```

`pbb.yaml`

```
version: v0.0.1
deps:
  - repo: "https://github.com/google/common-protos.git"
    rev: "master"
  - repo: "git@github.com:my-org/base-proto.git"
    rev: "v1.0.2"
gen:
  - language: "go"
    out: "."
    opt:
      - paths=source_relative
  - language: "dart"
    out: "./gen_dart"
input:
  - dir: src
  - path: proto/proto1.proto
lint:
  - opt:
    - a=a
    - b=b
```
