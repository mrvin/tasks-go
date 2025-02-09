package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/mrvin/tasks-go/pinger/internal/pinger/httpclient"
	"github.com/mrvin/tasks-go/pinger/internal/storage"
	probing "github.com/prometheus-community/pro-bing"
)

type Conf struct {
	PingPeriod       time.Duration
	PingTimeout      time.Duration
	PingCountPackets int
}

func Run(ctx context.Context, conf *Conf, confHTTP *httpclient.Conf) {
	ticker := time.Tick(conf.PingPeriod)
	for {
		select {
		case <-ticker:
			listHost, err := httpclient.ListHost(confHTTP.Addr)
			if err != nil {
				slog.Error("Get list host: " + err.Error())
				continue
			}
			slog.Info("Ping starting...")
			var wg sync.WaitGroup
			chResult := make(chan *storage.Ping)
			for _, host := range listHost {
				wg.Add(1)
				go func(ip net.IP) {
					defer wg.Done()
					pinger, err := probing.NewPinger(ip.String())
					if err != nil {
						slog.Error("New pinger: " + err.Error())
						return
					}
					pinger.Timeout = conf.PingTimeout
					pinger.Count = conf.PingCountPackets
					if err = pinger.Run(); err != nil {
						slog.Error("Run pinger: " + err.Error())
						return
					}
					stats := pinger.Statistics()

					if stats.PacketsRecv != 0 {
						chResult <- &storage.Ping{ip, stats.AvgRtt, time.Now()}
					}
				}(host.IP)
			}

			go func() {
				wg.Wait()
				slog.Info("Ping stopped")
				close(chResult)
			}()

			for ping := range chResult {
				fmt.Println(ping)
				if err := httpclient.CreatePing(confHTTP.Addr, ping); err != nil {
					slog.Error("Create ping: " + err.Error())
				}
			}
		case <-ctx.Done():
			slog.Info("Stop pinger")
			return
		}
	}
}
