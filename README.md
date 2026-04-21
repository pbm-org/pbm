# Pbm ProtoBuf Manage

Protobuf Manage for multi remote git repo

[![Release](https://github.com/pbm-org/pbm/actions/workflows/release.yml/badge.svg)](https://github.com/pbm-org/pbm/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/pbm-org/pbm)](https://goreportcard.com/report/github.com/pbm-org/pbm)

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
  - remote: https://github.com/bufbuild/protoc-gen-validate
    ref: main
gen:
  - plugin: go-lite
    out: .
    opt:
      - paths=source_relative
  - plugin: validate-go
    out: .
    opt:
      - paths=source_relative
input:
  - local: testdata/proto
    desc_out: testdata/proto/proto1.pb
```

## WIP

- [x] Builder
- [ ] Lint
- [ ] Break
- [ ] Protobuf Lsp

## Depend

[protoc](https://github.com/protocolbuffers/protobuf)
