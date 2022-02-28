package service

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"log"
	"sync"
)

var pubsubLocker sync.Mutex
var stopChannels = make(map[*pb.SDB_SubscribeServer]chan bool)
var subscribeServers = make(map[*pb.SDB_SubscribeServer]map[string]bool)

func Subscribe(topic []byte, subscribeServer *pb.SDB_SubscribeServer) (bool, error) {
	pubsubLocker.Lock()
	defer pubsubLocker.Unlock()

	if stopChannels[subscribeServer] == nil {
		stopChannels[subscribeServer] = make(chan bool)
	}
	if subscribeServers[subscribeServer] == nil {
		subscribeServers[subscribeServer] = make(map[string]bool)
	}
	subscribeServers[subscribeServer][string(topic)] = true
	return true, nil
}

func Publish(request *pb.PublishRequest) (bool, error) {
	go func() {
		message := &pb.Message{Topic: request.Topic, Payload: request.Payload}
		for subscribeServer, topics := range subscribeServers {
			if topics[string(request.Topic)] == true {
				if err := (*subscribeServer).Send(message); err != nil {
					log.Printf("Send: %+v to: %+v error, so stop", (*subscribeServer).Context(), message)
					stopChannels[subscribeServer] <- true
				}
			}
		}
	}()
	return true, nil
}

func CleanSubscribeServer(subscribeServer *pb.SDB_SubscribeServer) {
	close(stopChannels[subscribeServer])
	delete(stopChannels, subscribeServer)
	delete(subscribeServers, subscribeServer)
}

func GetStopChannel(subscribeServer *pb.SDB_SubscribeServer) chan bool {
	return stopChannels[subscribeServer]
}
