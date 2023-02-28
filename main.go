package main

import (
	"fmt"

	"go-micro/common/micro"
	"go-micro/common/micro/component"
)

func main() {
	server, err := micro.NewServer(
		"demo",
		"v100",
		[]micro.IComponent{
			&component.LoggingComponent{},
			&component.TracingComponent{},
			&component.GossipKVCacheComponent{
				ClusterName:   "platform-global",
				Port:          6666,
				InMachineMode: false,
			},
			&component.KafkaComponent{},
		},
	)
	if err != nil {
		panic(err)
	}
	err = server.Init()
	if err != nil {
		panic(err)
	}

	err = server.Run()
	if err != nil {
		fmt.Println(err)
	}
}
