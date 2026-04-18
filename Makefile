BINARY_NAME=pbm
BUILD_DIR=build
BUILD_FILE=cmd/main.go

PLATFORMS := linux-amd64 linux-arm64 windows-amd64 windows-arm64 darwin-amd64 darwin-arm64

.PHONY: build-all $(PLATFORMS)

build-all: $(PLATFORMS)

$(PLATFORMS):
	$(eval OS := $(word 1,$(subst -, ,$@)))
	$(eval ARCH := $(word 2,$(subst -, ,$@)))
	@echo "Build $(OS) $(ARCH)..."
	GOOS=$(OS) GOARCH=$(ARCH) go build -x -o $(BUILD_DIR)/$(BINARY_NAME)-$(OS)-$(ARCH)$(if $(findstring windows,$(OS)),.exe,) $(BUILD_FILE)
	
clean:
	@rm -rf $(BUILD_DIR)
