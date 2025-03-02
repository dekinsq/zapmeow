#!/bin/bash
# CentOS 7自动安装Go语言环境脚本（国际版）
# 用法：bash install_go.sh [版本号] [GOPATH路径]
# 示例：bash install_go.sh 1.21.6 /data/go

set -e

# 参数处理
GO_VERSION=${1:-"1.23.6"}  # 默认版本1.23.6
GOPATH_DIR=${2:-"/data/go"} # 默认GOPATH路径

# 清理旧版本
echo ">>> 清理旧版本..."
rm -rf /usr/local/go 2>/dev/null || true

# 下载安装包（使用官方源）
echo ">>> 下载Go ${GO_VERSION}..."
DOWNLOAD_URL="https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"

if ! wget -q --show-progress --progress=bar:force "$DOWNLOAD_URL" -O /tmp/go.tar.gz; then
    echo "下载失败，请检查："
    echo "1. 版本号是否存在（https://go.dev/dl/）"
    echo "2. 网络是否能访问dl.google.com"
    exit 1
fi

# 安装到系统目录
echo ">>> 安装到/usr/local..."
tar -C /usr/local -xzf /tmp/go.tar.gz
rm -f /tmp/go.tar.gz

# 配置环境变量
echo ">>> 配置环境变量..."
cat <<EOF | tee /etc/profile.d/golang.sh >/dev/null
export GOROOT=/usr/local/go
export GOPATH=${GOPATH_DIR}
export PATH=\$PATH:\$GOROOT/bin:\$GOPATH/bin
export GO111MODULE=on
EOF

# 创建工作目录
echo ">>> 创建GOPATH目录结构..."
mkdir -p ${GOPATH_DIR}/{bin,pkg,src}

# 应用配置
source /etc/profile.d/golang.sh

# 验证安装
echo ">>> 验证安装..."
go version || {
    echo "安装验证失败！";
    exit 1;
}

# 测试HelloWorld
echo ">>> 运行HelloWorld测试..."
cat <<EOF > ${GOPATH_DIR}/src/hello.go
package main
import "fmt"
func main() {
    fmt.Println("自动安装脚本测试成功！")
}
EOF

go run ${GOPATH_DIR}/src/hello.go

echo "√√√ Go ${GO_VERSION} 安装完成！"
echo "GOROOT: /usr/local/go"
echo "GOPATH: ${GOPATH_DIR}"
