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
go test -run=^$ -bench=. -benchtime=20s

goos: darwin
goarch: arm64
pkg: go-net-benchmark
BenchmarkTCPSpeed4K-8            1000000             20365 ns/op         201.13 MB/s         128 B/op          2 allocs/op
BenchmarkTCPSpeed64K-8            692217             32615 ns/op        2009.37 MB/s         128 B/op          2 allocs/op
BenchmarkTCPSpeed512K-8           182823            131516 ns/op        3986.49 MB/s         131 B/op          2 allocs/op
BenchmarkTCPSpeed1M-8             103548            229517 ns/op        4568.62 MB/s         138 B/op          2 allocs/op
BenchmarkKCPSpeed4K-8             347572             67618 ns/op          60.58 MB/s        3396 B/op        117 allocs/op
BenchmarkKCPSpeed64K-8             40382            620412 ns/op         105.63 MB/s       56910 B/op       1336 allocs/op
BenchmarkKCPSpeed512K-8             4468           5148164 ns/op         101.84 MB/s      467656 B/op      11574 allocs/op
BenchmarkKCPSpeed1M-8               2167          12491852 ns/op          83.94 MB/s      935754 B/op      22852 allocs/op
BenchmarkSmuxV1Speed4K-8          758518             28346 ns/op         144.50 MB/s         727 B/op         14 allocs/op
BenchmarkSmuxV1Speed64K-8         447194             53337 ns/op        1228.71 MB/s        1549 B/op         26 allocs/op
BenchmarkSmuxV1Speed512K-8        160027            155783 ns/op        3365.51 MB/s        6553 B/op        105 allocs/op
BenchmarkSmuxV1Speed1M-8           70965            294330 ns/op        3562.59 MB/s       12007 B/op        188 allocs/op
BenchmarkSmuxV2Speed4K-8          793562             33719 ns/op         121.48 MB/s         727 B/op         14 allocs/op
BenchmarkSmuxV2Speed64K-8         415078             54009 ns/op        1213.42 MB/s        1547 B/op         26 allocs/op
BenchmarkSmuxV2Speed512K-8        137670            176861 ns/op        2964.41 MB/s        6587 B/op        104 allocs/op
BenchmarkSmuxV2Speed1M-8           78840            304719 ns/op        3441.12 MB/s       12080 B/op        188 allocs/op
BenchmarkYamuxSpeed4K-8           732481             31790 ns/op         128.85 MB/s         305 B/op          6 allocs/op
BenchmarkYamuxSpeed64K-8          401103             60932 ns/op        1075.57 MB/s         395 B/op          8 allocs/op
BenchmarkYamuxSpeed512K-8       2023/07/10 18:38:59 [ERR] yamux: Failed to write header: write tcp 127.0.0.1:52302->127.0.0.1:10076: write: broken pipe
   90397            277266 ns/op        1890.92 MB/s        1690 B/op         36 allocs/op
BenchmarkYamuxSpeed1M-8            46183            518529 ns/op        2022.21 MB/s        3302 B/op         68 allocs/op
PASS
ok      go-net-benchmark        578.657s
```
