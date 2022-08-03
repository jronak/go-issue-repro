# HTTP2 Drops Conn-Level Flow Control Update
Reproduce issue where HTTP2 drops connection-evel flow control update eventually starving client connection from writing data frames and hangs client connection forever.

Repro for https://github.com/golang/go/issues/54185 

# What version of Go are you using (go version)?
go version go1.18.3 darwin/amd64

# What operating system and processor architecture are you using (go env)?
<details>
<summary>go env</summary>

```
GO111MODULE=""
GOARCH="amd64"
GOBIN=""
GOCACHE="/Users/ronakj/Library/Caches/go-build"
GOENV="/Users/ronakj/Library/Application Support/go/env"
GOEXE=""
GOEXPERIMENT=""
GOFLAGS=""
GOHOSTARCH="amd64"
GOHOSTOS="darwin"
GOINSECURE=""
GOMODCACHE="/Users/ronakj/gocode/pkg/mod"
GONOPROXY="none"
GONOSUMDB="*"
GOOS="darwin"
GOPATH="/Users/ronakj/gocode"
GOPRIVATE=""
GOPROXY="https://proxy.golang.org,direct"
GOROOT="/usr/local/go"
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GOTOOLDIR="/usr/local/go/pkg/tool/darwin_amd64"
GOVCS=""
GOVERSION="go1.18.3"
GCCGO="gccgo"
GOAMD64="v1"
AR="ar"
CC="clang"
CXX="clang++"
CGO_ENABLED="1"
GOMOD="/Users/ronakj/project/http2-issue-repro/go.mod"
GOWORK=""
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -arch x86_64 -m64 -pthread -fno-caret-diagnostics -Qunused-arguments -fmessage-length=0 -fdebug-prefix-map=/var/folders/4d/2jw_2tc15x339gr53x6k64hm0000gn/T/go-build1453600270=/tmp/go-build -gno-record-gcc-switches -fno-common"
```
</details>

# What is unexpected?
Send unary request with 16KB payload and `Content-length` metadata set to invalid size `2`. After 64 requests, client requests hang forever as outbound control flow length has gone down to `0`. This happens after 64 requests as it takes 64 requests of 16KB each to drain initial receive window size of 1MB (1<<20). Client to server flow control is never updated when HTTP2 server encounters invalid `Content-length` header.

# What is expected behaviour?
Stream with invalid header must result in stream reset, and it must not affect other streams or connection flow control.

# Logs
<details>

```
022/08/02 00:54:23 Failed stream num:0 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:1 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:2 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:3 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:4 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:5 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:6 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:7 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:8 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:9 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:10 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:11 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:12 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:13 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:14 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:15 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:16 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:17 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:18 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:19 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:20 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:21 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:22 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:23 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:24 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:25 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:26 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:27 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:28 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:29 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:30 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:31 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:32 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:33 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:34 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:35 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:36 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:37 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:38 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:39 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:40 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:41 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:42 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:43 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:44 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:45 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:46 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:47 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:48 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:49 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:50 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:51 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:52 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:53 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:54 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:55 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:56 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:57 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:58 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:59 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:60 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:61 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:62 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
2022/08/02 00:54:23 Failed stream num:63 with err:rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
# Hangs forever
```
</details>