# HTTP2: slow streams can potentially block other faster streams
Slow stream(s) consuming entire client conn flow control window can block and slowdown other faster streams.

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

# What did you do?
I'm using a Go HTTP2 (h2c) server as an Echo server, where handler writes response immediately for requests with priority header `request-type=priority` otherwise delays response by 5s. I use HTTP2 client to dispatch two concurrent requests with 2MB payload each, and one of them has priority header set.

# What did you expect to see?
Priority response must arrive immediately. Non-priority response must arrive after 5s delay.

# What did you see instead?
Both non-priority and priority responses arrived after 5s. This should not have happened as HTTP2 streams on same connections must not interfere with each other.

# Logs
```
H2c Server starting on :8081
low priority req response time: 5.015591042s 
high priority req response time: 5.023305906s   # expected this response immediately.
```
</details>

# Root Cause
Pre-requiste: https://httpwg.org/specs/rfc7540.html#FlowControl

HTTP2 client and server maintain a flow control window which is the maximum number of Data frame bytes each are willing to accept. Flow control window is applied separately at the connection and stream level. Data frames cannot be forwarded when flow control has been exhausted, client/server must wait for WINDOW_UPDATE frame from other side to replenish the window before trying to write DATA frames. 

Go HTTP2 server sets client flow control window size to [1MB](https://github.com/golang/net/blob/master/http2/server.go#L145). When dispatch two requests with 2MB payload in each, this exhausts the client side connection flow control window size of 1MB. HTTP2 server buffers the request payload internally and does not replenish client conn & stream flow control window until server handler has read the buffered data. 

Since our handler delays reading payload from the low-priority request, this holds server from sending WINDOW_UPDATE for both connection and stream. This blocks client from sending the payload of the high priority request as client connection flow control will not be replenished until low-priority request payload is read.

# Test Fix
Below patch fixes the issue in HTTP2 server:
```
diff --git a/http2/server.go b/http2/server.go
index 47524a6..bc8c6c1 100644
--- a/http2/server.go
+++ b/http2/server.go
@@ -1775,9 +1775,10 @@ func (sc *serverConn) processData(f *DataFrame) error {
 		// Return any padded flow control now, since we won't
 		// refund it later on body reads.
 		if pad := int32(f.Length) - int32(len(data)); pad > 0 {
-			sc.sendWindowUpdate32(nil, pad)
 			sc.sendWindowUpdate32(st, pad)
 		}
+
+		sc.sendWindowUpdate32(nil, int32(f.Length))
 	}
 	if f.StreamEnded() {
 		st.endStream()
@@ -2317,7 +2318,6 @@ func (sc *serverConn) noteBodyReadFromHandler(st *stream, n int, err error) {
 
 func (sc *serverConn) noteBodyRead(st *stream, n int) {
 	sc.serveG.check()
-	sc.sendWindowUpdate(nil, n) // conn-level
 	if st.state != stateHalfClosedRemote && st.state != stateClosed {
 		// Don't send this WINDOW_UPDATE if the stream is closed
 		// remotely.

```

Logs with fix:
```
H2c Server starting on :8081
high priority req response time: 28.067782ms # priority stream response arrived immediately.
low priority req response time: 5.013978621s
```