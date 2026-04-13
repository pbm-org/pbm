# Pbm ProtoBuf Manage

Protobuf Manage for multi remote git repo

## How to use

init pbm project

```
pbm init
```

build proto protobuf

```
pbm gen
```

update dep buf latest

```
pbm dep
```

clean pbm dep

```
pbm clean
```

`pbm.yaml`

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
  - git@github.com:labulakalia/pbm.git
  - proto_dir
lint:
  - opt:
    - a=a
    - b=b
```

## DEV

- [ ] Builder
- [ ] Lint
- [ ] Break

## Depend

[protoc](https://github.com/protocolbuffers/protobuf)  
[protolint](https://github.com/yoheimuta/protolint)
