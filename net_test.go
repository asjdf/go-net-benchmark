package go_net_benchmark

import (
	"fmt"
	"github.com/asjdf/smux"
	"github.com/hashicorp/yamux"
	"github.com/xtaci/kcp-go/v5"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
)

var basePort = uint32(10000) // base of echo server

func handleEcho(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		go func(conn net.Conn) {
			_, err := io.Copy(conn, conn)
			if err != nil {
				return
			}
		}(conn)
	}
}

func getTCPConnectionPair() (net.Listener, net.Conn, error) {
	port := int(atomic.AddUint32(&basePort, 1))
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "localhost", port))
	if err != nil {
		return nil, nil, err
	}

	cs, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		return nil, nil, err
	}

	return listener, cs, nil
}

func getKCPConnectionPair() (net.Listener, net.Conn, error) {
	port := int(atomic.AddUint32(&basePort, 1))
	block, _ := kcp.NewNoneBlockCrypt(nil)
	listener, err := kcp.ListenWithOptions(fmt.Sprintf("%s:%d", "localhost", port), block, 10, 0)
	if err != nil {
		return nil, nil, err
	}
	listener.SetReadBuffer(16 * 1024 * 1024)
	listener.SetWriteBuffer(16 * 1024 * 1024)

	cs, err := kcp.DialWithOptions(fmt.Sprintf("%s:%d", "localhost", port), block, 10, 3)
	if err != nil {
		return nil, nil, err
	}
	cs.SetWindowSize(2048, 2048)
	cs.SetNoDelay(1, 10, 2, 1)
	cs.SetReadBuffer(16 * 1024 * 1024)
	cs.SetWriteBuffer(16 * 1024 * 1024)
	cs.SetMtu(1400)

	return listener, cs, nil
}

func getSmuxV1StreamPair() (net.Listener, net.Conn, error) {
	listen, c2, err := getTCPConnectionPair()
	if err != nil {
		return nil, nil, err
	}
	c1, err := listen.Accept()
	if err != nil {
		return nil, nil, err
	}

	conf := smux.DefaultConfig()
	conf.Version = 1
	conf.MaxFrameSize = 65535
	conf.MaxReceiveBuffer = 100 * 1024 * 1024 // 100M
	conf.MaxStreamBuffer = 100 * 1024 * 1024  // 100M
	s, err := smux.Server(c1, conf)
	if err != nil {
		return nil, nil, err
	}
	c, err := smux.Client(c2, conf)
	if err != nil {
		return nil, nil, err
	}
	cs, err := c.OpenStream()
	if err != nil {
		return nil, nil, err
	}

	return s, cs, nil
}

func getSmuxV2StreamPair() (net.Listener, net.Conn, error) {
	listen, c2, err := getTCPConnectionPair()
	if err != nil {
		return nil, nil, err
	}
	c1, err := listen.Accept()
	if err != nil {
		return nil, nil, err
	}

	conf := smux.DefaultConfig()
	conf.Version = 2
	conf.MaxFrameSize = 65535
	conf.MaxReceiveBuffer = 100 * 1024 * 1024 // 100M
	conf.MaxStreamBuffer = 100 * 1024 * 1024  // 100M
	s, err := smux.Server(c1, conf)
	if err != nil {
		return nil, nil, err
	}
	c, err := smux.Client(c2, conf)
	if err != nil {
		return nil, nil, err
	}
	cs, err := c.OpenStream()
	if err != nil {
		return nil, nil, err
	}

	return s, cs, nil
}

func getYamuxStreamPair() (net.Listener, net.Conn, error) {
	listen, c2, err := getTCPConnectionPair()
	if err != nil {
		return nil, nil, err
	}
	c1, err := listen.Accept()
	if err != nil {
		return nil, nil, err
	}

	conf := yamux.DefaultConfig()
	conf.MaxStreamWindowSize = 100 * 1024 * 1024
	s, err := yamux.Server(c1, conf)
	if err != nil {
		return nil, nil, err
	}
	c, err := yamux.Client(c2, conf)
	if err != nil {
		return nil, nil, err
	}
	cs, err := c.OpenStream()
	if err != nil {
		return nil, nil, err
	}

	return s, cs, nil
}

func speedBenchmark(b *testing.B, server net.Listener, conn1 net.Conn, msgLen int) error {
	b.ReportAllocs()
	b.SetBytes(int64(msgLen))

	go handleEcho(server)

	b.ResetTimer()
	buf := make([]byte, msgLen)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			// send packet
			if _, err := conn1.Write(buf); err != nil {
				return
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			// receive packet
			nrecv := 0
			for {
				n, err := conn1.Read(buf)
				if err != nil {
					return
				} else {
					nrecv += n
					if nrecv >= msgLen {
						break
					}
				}
			}
			wg.Done()
		}()
		wg.Wait()
	}
	_ = server.Close()
	_ = conn1.Close()
	return nil
}

func BenchmarkTCPSpeed4K(b *testing.B) {
	conn0, conn1, err := getTCPConnectionPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 4*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkTCPSpeed64K(b *testing.B) {
	conn0, conn1, err := getTCPConnectionPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 64*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkTCPSpeed512K(b *testing.B) {
	conn0, conn1, err := getTCPConnectionPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 512*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkTCPSpeed1M(b *testing.B) {
	conn0, conn1, err := getTCPConnectionPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 1024*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkKCPSpeed4K(b *testing.B) {
	conn0, conn1, err := getKCPConnectionPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 4*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkKCPSpeed64K(b *testing.B) {
	conn0, conn1, err := getKCPConnectionPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 64*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkKCPSpeed512K(b *testing.B) {
	conn0, conn1, err := getKCPConnectionPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 512*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkKCPSpeed1M(b *testing.B) {
	conn0, conn1, err := getKCPConnectionPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 1024*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkSmuxV1Speed4K(b *testing.B) {
	conn0, conn1, err := getSmuxV1StreamPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 4*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkSmuxV1Speed64K(b *testing.B) {
	conn0, conn1, err := getSmuxV1StreamPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 64*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkSmuxV1Speed512K(b *testing.B) {
	conn0, conn1, err := getSmuxV1StreamPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 512*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkSmuxV1Speed1M(b *testing.B) {
	conn0, conn1, err := getSmuxV1StreamPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 1024*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkSmuxV2Speed4K(b *testing.B) {
	conn0, conn1, err := getSmuxV2StreamPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 4*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkSmuxV2Speed64K(b *testing.B) {
	conn0, conn1, err := getSmuxV2StreamPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 64*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkSmuxV2Speed512K(b *testing.B) {
	conn0, conn1, err := getSmuxV2StreamPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 512*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkSmuxV2Speed1M(b *testing.B) {
	conn0, conn1, err := getSmuxV2StreamPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 1024*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkYamuxSpeed4K(b *testing.B) {
	conn0, conn1, err := getYamuxStreamPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 4*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkYamuxSpeed64K(b *testing.B) {
	conn0, conn1, err := getYamuxStreamPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 64*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkYamuxSpeed512K(b *testing.B) {
	conn0, conn1, err := getYamuxStreamPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 512*1024)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkYamuxSpeed1M(b *testing.B) {
	conn0, conn1, err := getYamuxStreamPair()
	if err != nil {
		b.Error(err)
	}
	err = speedBenchmark(b, conn0, conn1, 1024*1024)
	if err != nil {
		b.Error(err)
	}
}
