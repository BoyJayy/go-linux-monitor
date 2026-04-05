package main

import (
	"agent/internal/config"
	"agent/internal/sender"
	"agent/internal/snapshot"
	"log"
)

/*type snapshotResult struct {
	stat snapshot.Metrics
	err  error
}*/

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	sender := sender.NewHTTP(cfg.ServerURL, cfg.RequestTimeout)
	for {
		/*metricsCh := make(chan snapshotResult, 1)
		go func() {
			metric, err := snapshot.BuildSnapshot(interval)
			metricsCh <- snapshotResult{stat: metric, err: err}
		}()*/
		metrics, err := snapshot.BuildSnapshot(cfg.CollectionInterval) // _ while i want implement http/grpc
		//time.Sleep(interval) - as for me bad asf cuz when we receiving data - already spending interval so that's gonna be interval*2 time (for that version it could be fine but ill try to implement one that gonna fit the most there)
		if err != nil {
			log.Printf("Error while reading metrics: %v", err)
			continue
		}
		//fmt.Printf("%+v\n", metrics)
		err = sender.Send(metrics)
		if err != nil {
			log.Printf("Error while sending metrics: %v", err)
			continue
		}
		// any realization like sending with our sender (not done yet ofc :-) )
	}
}

// тест cборки метрик
/*
❯ docker run --rm -it \
  -v "$PWD":/usr/src/app \
  -w /usr/src/app \
  golang:1.25.5 \
  go run ./cmd/agent
go: downloading golang.org/x/sys v0.42.0
{0 0 0 100 0 0 0 0 0 map[cpu0:{0 0 0 100 0 0 0 0 0} cpu1:{0 0 0 100 0 0 0 0 0} cpu2:{0 0 0 100 0 0 0 0 0} cpu3:{0 0 0 100 0 0 0 0 0} cpu4:{0 0 0 100 0 0 0 0 0} cpu5:{0 0 0 100 0 0 0 0 0} cpu6:{0 0 0 100 0 0 0 0 0} cpu7:{0 0 0 100 0 0 0 0 0}]}
[{/ 485473984512 481838170112 3635814400 0.7489205428082273} {/etc/resolv.conf 485473984512 481838170112 3635814400 0.7489205428082273} {/etc/hostname 485473984512 481838170112 3635814400 0.7489205428082273} {/etc/hosts 485473984512 481838170112 3635814400 0.7489205428082273} {/usr/src/app 126562507685888 82537996091392 44024511594496 34.78479717212875}]
{2081261 21596 0 0 [{lo 0 0 0 0} {tunl0 0 0 0 0} {gre0 0 0 0 0} {gretap0 0 0 0 0} {erspan0 0 0 0 0} {ip_vti0 0 0 0 0} {ip6_vti0 0 0 0 0} {sit0 0 0 0 0} {ip6tnl0 0 0 0 0} {ip6gre0 0 0 0 0} {eth0 2081261 21596 0 0}]}
{8217731072 7539712000 678019072 8.250684599672429}
{0 0 0.4987531172069825 99.50124688279301 0 0 0 0 0.4987531172069825 map[cpu0:{0 0 1 99 0 0 0 0 1} cpu1:{0 0 0 100 0 0 0 0 0} cpu2:{0 0 3 97 0 0 0 0 3} cpu3:{0 0 0 100 0 0 0 0 0} cpu4:{0 0 0 100 0 0 0 0 0} cpu5:{0 0 0 100 0 0 0 0 0} cpu6:{0 0 0 100 0 0 0 0 0} cpu7:{0 0 0 100 0 0 0 0 0}]}
[{/ 485473984512 481838166016 3635818496 0.7489213865197609} {/etc/resolv.conf 485473984512 481838166016 3635818496 0.7489213865197609} {/etc/hostname 485473984512 481838166016 3635818496 0.7489213865197609} {/etc/hosts 485473984512 481838166016 3635818496 0.7489213865197609} {/usr/src/app 126562507685888 82513020059648 44049487626240 34.8045313194688}]
{2081303 21638 0 0 [{lo 0 0 0 0} {tunl0 0 0 0 0} {gre0 0 0 0 0} {gretap0 0 0 0 0} {erspan0 0 0 0 0} {ip_vti0 0 0 0 0} {ip6_vti0 0 0 0 0} {sit0 0 0 0 0} {ip6tnl0 0 0 0 0} {ip6gre0 0 0 0 0} {eth0 2081303 21638 0 0}]}
{8217731072 7542312960 675418112 8.219034014161519}
{0.12422360248447205 0 0.12422360248447205 99.75155279503105 0 0 0 0 0.2484472049689441 map[cpu0:{0 0 0 99.00990099009901 0 0 0.9900990099009901 0 0.9900990099009901} cpu1:{0 0 0 100 0 0 0 0 0} cpu2:{0 0 0 100 0 0 0 0 0} cpu3:{0 0 0 100 0 0 0 0 0} cpu4:{0 0 0 100 0 0 0 0 0} cpu5:{0 0 0 100 0 0 0 0 0} cpu6:{0 0 0 100 0 0 0 0 0} cpu7:{0 0 0 100 0 0 0 0 0}]}
[{/ 485473984512 481838166016 3635818496 0.7489213865197609} {/etc/resolv.conf 485473984512 481838166016 3635818496 0.7489213865197609} {/etc/hostname 485473984512 481838166016 3635818496 0.7489213865197609} {/etc/hosts 485473984512 481838166016 3635818496 0.7489213865197609} {/usr/src/app 126562507685888 82512878501888 44049629184000 34.80464316756868}]
{2081303 21638 0 0 [{lo 0 0 0 0} {tunl0 0 0 0 0} {gre0 0 0 0 0} {gretap0 0 0 0 0} {erspan0 0 0 0 0} {ip_vti0 0 0 0 0} {ip6_vti0 0 0 0 0} {sit0 0 0 0 0} {ip6tnl0 0 0 0 0} {ip6gre0 0 0 0 0} {eth0 2081303 21638 0 0}]}
{8217731072 7542439936 675291136 8.217488867467285}
{0 0 0 99.87562189054727 0 0 0.12437810945273632 0 0.12437810945273632 map[cpu0:{0 0 0 100 0 0 0 0 0} cpu1:{0 0 0 100 0 0 0 0 0} cpu2:{0 0 0 100 0 0 0 0 0} cpu3:{0 0 0 100 0 0 0 0 0} cpu4:{0 0 0 100 0 0 0 0 0} cpu5:{0 0 0 100 0 0 0 0 0} cpu6:{0 0 0 100 0 0 0 0 0} cpu7:{0 0 0 100 0 0 0 0 0}]}
[{/ 485473984512 481838166016 3635818496 0.7489213865197609} {/etc/resolv.conf 485473984512 481838166016 3635818496 0.7489213865197609} {/etc/hostname 485473984512 481838166016 3635818496 0.7489213865197609} {/etc/hosts 485473984512 481838166016 3635818496 0.7489213865197609} {/usr/src/app 126562507685888 82512834461696 44049673224192 34.80467796475531}]
{2081369 21704 0 0 [{lo 0 0 0 0} {tunl0 0 0 0 0} {gre0 0 0 0 0} {gretap0 0 0 0 0} {erspan0 0 0 0 0} {ip_vti0 0 0 0 0} {ip6_vti0 0 0 0 0} {sit0 0 0 0 0} {ip6tnl0 0 0 0 0} {ip6gre0 0 0 0 0} {eth0 2081369 21704 0 0}]}
{8217731072 7541489664 676241408 8.229052545953161}
*/
