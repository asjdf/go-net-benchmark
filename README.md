# go-net-benchmark

## 目标

测试 Golang 的 net 标准库及 kcp 库的速度、效率、延迟以及在恶劣网络环境下的可用性

## 已测试特性

- [x] 数据传输速率
- [ ] 链路/流建立速率

## 已测试列表

- TCP
- [KCP](https://github.com/xtaci/kcp-go)
- [Smux](https://github.com/asjdf/smux) Fork版（增加net.TCPConn替代能力）
- [Yamux](https://github.com/hashicorp/yamux)

## 结果

> MacBook Air M1 8+256

```shell
go test -v -run=^$ -bench=. -benchtime=20s

goos: darwin
goarch: arm64
pkg: go-net-benchmark
BenchmarkTCPSpeed4K
BenchmarkTCPSpeed4K-8            1000000             20207 ns/op         202.70 MB/s         128 B/op          2 allocs/op
BenchmarkTCPSpeed64K
BenchmarkTCPSpeed64K-8            726085             36170 ns/op        1811.89 MB/s         128 B/op          2 allocs/op
BenchmarkTCPSpeed512K
BenchmarkTCPSpeed512K-8           187071            129617 ns/op        4044.90 MB/s         131 B/op          2 allocs/op
BenchmarkTCPSpeed1M
BenchmarkTCPSpeed1M-8             105732            268275 ns/op        3908.59 MB/s         138 B/op          2 allocs/op
BenchmarkKCPSpeed4K
BenchmarkKCPSpeed4K-8             352902             68119 ns/op         60.13 MB/s         3399 B/op        117 allocs/op
BenchmarkKCPSpeed64K
BenchmarkKCPSpeed64K-8             40116            619485 ns/op        105.79 MB/s        56890 B/op       1335 allocs/op
BenchmarkKCPSpeed512K
BenchmarkKCPSpeed512K-8             4670           5088493 ns/op        103.03 MB/s       468138 B/op      11629 allocs/op
BenchmarkKCPSpeed1M
BenchmarkKCPSpeed1M-8               2260          10716793 ns/op         97.84 MB/s       948663 B/op      23855 allocs/op
BenchmarkSmuxV1Speed4K
BenchmarkSmuxV1Speed4K-8          828534             28255 ns/op        144.97 MB/s          728 B/op         14 allocs/op
BenchmarkSmuxV1Speed64K
BenchmarkSmuxV1Speed64K-8         435939             53453 ns/op       1226.05 MB/s         1548 B/op         26 allocs/op
BenchmarkSmuxV1Speed512K
BenchmarkSmuxV1Speed512K-8        152894            153878 ns/op       3407.16 MB/s         6543 B/op        105 allocs/op
BenchmarkSmuxV1Speed1M
BenchmarkSmuxV1Speed1M-8           84783            306250 ns/op       3423.92 MB/s        12104 B/op        187 allocs/op
BenchmarkSmuxV2Speed4K
BenchmarkSmuxV2Speed4K-8          807070             28406 ns/op        144.20 MB/s          727 B/op         14 allocs/op
BenchmarkSmuxV2Speed64K
BenchmarkSmuxV2Speed64K-8         447172             53351 ns/op       1228.38 MB/s         1546 B/op         26 allocs/op
BenchmarkSmuxV2Speed512K
BenchmarkSmuxV2Speed512K-8        161622            157755 ns/op       3323.43 MB/s         6553 B/op        105 allocs/op
BenchmarkSmuxV2Speed1M
BenchmarkSmuxV2Speed1M-8           81782            295704 ns/op       3546.03 MB/s        12097 B/op        188 allocs/op
BenchmarkYamuxSpeed4K
BenchmarkYamuxSpeed4K-8           754335             30560 ns/op        134.03 MB/s          305 B/op          6 allocs/op
BenchmarkYamuxSpeed64K
BenchmarkYamuxSpeed64K-8          410257             58453 ns/op       1121.17 MB/s          395 B/op          8 allocs/op
BenchmarkYamuxSpeed512K
2023/07/10 15:52:48 [ERR] yamux: Failed to write header: write tcp 127.0.0.1:63337->127.0.0.1:10074: write: broken pipe
BenchmarkYamuxSpeed512K-8          92732            261870 ns/op       2002.09 MB/s         1688 B/op         36 allocs/op
BenchmarkYamuxSpeed1M
2023/07/10 15:53:20 [ERR] yamux: Failed to write header: write tcp 127.0.0.1:63368->127.0.0.1:10079: write: broken pipe
BenchmarkYamuxSpeed1M-8            49674            490873 ns/op       2136.15 MB/s         3282 B/op         68 allocs/op
PASS
ok      go-net-benchmark        549.487s

```
