# Pbm ProtoBuf Manage

Protobuf Manage for multi remote git repo

# Install

[Download](https://github.com/pbm-org/pbm/releases/tag/v0.0.1)

## How to use

init pbm project

```
pbm -init
```

build proto protobuf

```
pbm -build
```

update dep buf latest

```
pbm -update
```

clean pbm dep

```
pbm -clean
```

example pbm.yaml

```
version: v1
deps:
  - remote: https://cnb.cool/medianexapp/plugin_api
    ref: main
  - remote: https://github.com/googleapis/googleapis
    ref: master
gen:
  - plugin: go
    out: .
    opt:
      - paths=source_relative
input:
  - local: testdata/proto/proto1.proto
    desc_out: testdata/proto/proto1.pb
```

## WIP

- [x] Builder
- [ ] Lint
- [ ] Break
- [ ] Protobuf Lsp

## Depend

[protoc](https://github.com/protocolbuffers/protobuf)
