BINARY_NAME=git-cm
VERSION=1.0.0
BUILD_DIR=build

# 支持的平台列表
PLATFORMS := \
    linux-amd64 \
    linux-arm64 \
    darwin-amd64 \
    darwin-arm64

# 为每个平台设置环境变量
define build_platform
	@echo "Building for $(1)..."
	GOOS=$(word 1,$(subst -, ,$(1))) \
	GOARCH=$(word 2,$(subst -, ,$(1))) \
	go build -o $(BUILD_DIR)/$(BINARY_NAME)-$(1) -ldflags "-X main.version=$(VERSION)" .
endef

# 默认构建所有平台
all: clean $(PLATFORMS)

# 单个平台构建目标
$(PLATFORMS):
	$(call build_platform,$@)

# 本地构建（当前平台）
local:
	go build -o $(BINARY_NAME) -ldflags "-X main.version=$(VERSION)" .

# 清理构建产物
clean:
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@mkdir -p $(BUILD_DIR)

# 打包成zip文件（便于GitHub Release）
package: all
	@for platform in $(PLATFORMS); do \
		zip -j $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-$$platform.zip $(BUILD_DIR)/$(BINARY_NAME)-$$platform; \
	done
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)-*

.PHONY: all local clean package $(PLATFORMS)
