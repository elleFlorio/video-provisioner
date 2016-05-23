package discovery

import (
	"errors"
	"log"
	"time"

	"github.com/elleFlorio/video-provisioner/Godeps/_workspace/src/github.com/coreos/etcd/client"
	"github.com/elleFlorio/video-provisioner/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/elleFlorio/video-provisioner/utils"
)

var (
	uuid              string
	myKey             string
	myAddress         string
	myName            string
	kAPI              client.KeysAPI
	ErrNoDestinations = errors.New("No destinations available")
)

func InitializeEtcd(uri string) error {
	cfg := client.Config{
		Endpoints: []string{uri},
	}

	etcd, err := client.New(cfg)
	if err != nil {
		return err
	}

	kAPI = client.NewKeysAPI(etcd)

	//This is needed to probe if the etcd server is reachable
	_, err = kAPI.Set(
		context.Background(),
		"/probe",
		"etcd",
		&client.SetOptions{TTL: time.Duration(1) * time.Millisecond},
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func RegisterToEtcd(name string, address string) error {
	var err error

	uuid, err = utils.GenerateUUID()
	if err != nil {
		log.Println(err)
		return err
	}

	myKey = "app/video/services/" + name + "/" + uuid
	myAddress = address
	myName = name

	_, err = kAPI.Set(
		context.Background(),
		myKey,
		myAddress,
		&client.SetOptions{TTL: time.Duration(5) * time.Second},
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func UnregisterFromEtcd() {
	_, err := kAPI.Delete(context.Background(), myKey, nil)
	if err != nil {
		log.Println(err.Error())
		log.Println("Cannot unregister from etcd")
	}
}

func KeepAlive(ch_stop chan struct{}) {
	var err error
	ticker := time.NewTicker(time.Duration(5) * time.Second)

	for {
		select {
		case <-ticker.C:
			_, err = kAPI.Set(
				context.Background(),
				myKey,
				myAddress,
				&client.SetOptions{TTL: time.Duration(5) * time.Second},
			)
			if err != nil {
				log.Println(err)
				log.Println("Cannot keep the agent Alive")
			}
		case <-ch_stop:
			return
		}
	}
}

func GetAvailableInstances(service string) ([]string, error) {
	key := "app/video/services/" + service + "/"
	available := []string{}
	resp, err := kAPI.Get(context.Background(), key, nil)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	for _, n := range resp.Node.Nodes {
		available = append(available, n.Value)
	}

	if len(available) < 1 {
		log.Println(ErrNoDestinations)
		return []string{}, ErrNoDestinations
	}

	return available, nil
}

func AddRequestToHistory(reqID string, start time.Time) error {
	key := "app/video/requests/" + myName + "/" + reqID
	value := start.Format(time.RFC3339Nano)

	_, err := kAPI.Set(
		context.Background(),
		key,
		value,
		nil,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetRequestStartFromHistory(reqID string) (time.Time, error) {
	key := "app/video/requests/" + myName + "/" + reqID

	resp, err := kAPI.Get(context.Background(), key, nil)
	if err != nil {
		log.Println(err)
		return time.Now(), err
	}

	value := resp.Node.Value
	start, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		log.Println(err)
		return time.Now(), err
	}

	return start, nil
}

func RemoveRequestFromHistory(reqID string) error {
	key := "app/video/requests/" + myName + "/" + reqID
	_, err := kAPI.Delete(context.Background(), key, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
