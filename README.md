# Pbb

Protobuf Build system for multi remote git repo

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

clean pbb dep

```
pbb clean
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
  - proto/proto1.proto
  - git@github.com:labulakalia/pbb.git
  - proto_dir
lint:
  - opt:
    - a=a
    - b=b
```

## Depend

[protoc](https://github.com/protocolbuffers/protobuf)  
[protolint](https://github.com/yoheimuta/protolint)
