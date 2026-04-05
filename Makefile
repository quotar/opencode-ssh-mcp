# Makefile for opencode-ssh-mcp

# 项目变量
BINARY_NAME=opencode-ssh-mcp
VERSION=$(shell git describe --tags --always --dirty="-dev" 2>/dev/null || echo "v1.0.0")
BUILD_TIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# Go 相关变量
GO?=go
GOFMT?=gofmt
GOLINT?=golangci-lint

# 构建标志
LDFLAGS=-ldflags="-s -w -X main.version=${VERSION} -X main.commit=${GIT_COMMIT} -X main.buildTime=${BUILD_TIME}"

.PHONY: all build clean test format lint vet install uninstall generate

all: build

# 构建二进制文件
build:
	@echo "Building ${BINARY_NAME} version ${VERSION}"
	${GO} build ${LDFLAGS} -o ${BINARY_NAME}

# 构建并安装到 GOPATH
install:
	${GO} install ${LDFLAGS} .

# 从 GOPATH 卸载
uninstall:
	@echo "Removing ${BINARY_NAME} from $(GOPATH)/bin"
	rm -f $(GOPATH)/bin/${BINARY_NAME}

# 清理构建产物
clean:
	@echo "Cleaning build artifacts"
	rm -f ${BINARY_NAME}
	rm -rf dist/

# 运行测试
test:
	@echo "Running tests"
	${GO} test ./...

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "Running tests with coverage"
	${GO} test -coverprofile=coverage.out ./...
	${GO} tool cover -html=coverage.out -o coverage.html

# 格式化代码
format:
	@echo "Formatting code"
	${GOFMT} -s -w .

# 运行 vet 检查
vet:
	@echo "Running go vet"
	${GO} vet ./...

# 运行 lint 检查（如果有 golangci-lint）
lint:
	@if command -v ${GOLINT} >/dev/null 2>&1; then \
		echo "Running golangci-lint"; \
		${GOLINT} run; \
	else \
		echo "golangci-lint not found, skipping lint"; \
	fi

# 生成 Go 代码（如果有生成的代码）
generate:
	${GO} generate ./...

# 构建所有平台
build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64

build-linux-amd64:
	GOOS=linux GOARCH=amd64 ${GO} build ${LDFLAGS} -o ${BINARY_NAME}_linux_amd64

build-linux-arm64:
	GOOS=linux GOARCH=arm64 ${GO} build ${LDFLAGS} -o ${BINARY_NAME}_linux_arm64

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 ${GO} build ${LDFLAGS} -o ${BINARY_NAME}_darwin_amd64

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 ${GO} build ${LDFLAGS} -o ${BINARY_NAME}_darwin_arm64

build-windows-amd64:
	GOOS=windows GOARCH=amd64 ${GO} build ${LDFLAGS} -o ${BINARY_NAME}_windows_amd64.exe

# 打包发行版
release: build-all
	@echo "Creating release archives"
	mkdir -p dist
	tar -czvf dist/${BINARY_NAME}_${VERSION}_linux_amd64.tar.gz ${BINARY_NAME}_linux_amd64
	tar -czvf dist/${BINARY_NAME}_${VERSION}_linux_arm64.tar.gz ${BINARY_NAME}_linux_arm64
	tar -czvf dist/${BINARY_NAME}_${VERSION}_darwin_amd64.tar.gz ${BINARY_NAME}_darwin_amd64
	tar -czvf dist/${BINARY_NAME}_${VERSION}_darwin_arm64.tar.gz ${BINARY_NAME}_darwin_arm64
	tar -czvf dist/${BINARY_NAME}_${VERSION}_windows_amd64.tar.gz ${BINARY_NAME}_windows_amd64.exe
	@echo "Release archives created in dist/"

# 运行所有检查
check: vet lint test

# 打印项目信息
info:
	@echo "Project: ${BINARY_NAME}"
	@echo "Version: ${VERSION}"
	@echo "Commit: ${GIT_COMMIT}"
	@echo "Build Time: ${BUILD_TIME}"
	@echo "Go Version: $(shell ${GO} version)"