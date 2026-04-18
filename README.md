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
pbm update
```

clean pbm dep

```
pbm clean
```

`pbm.yaml`

```
version: v1
deps:
  - remote: https://cnb.cool/medianexapp/plugin_api
    ref: main
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

## Depend

[protoc](https://github.com/protocolbuffers/protobuf)
