package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/apache/pulsar-client-go/pulsar"
)

func getenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, errors.New(fmt.Sprintf("env var %s is empty", key))
	}
	return v, nil
}

func getenvInt(key string) (int, error) {
	s, err := getenvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func produce(wg sync.WaitGroup, client pulsar.Client, queue int) {
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: fmt.Sprintf("persistent://public/default/my-topic%d", queue),
	})

	if err != nil {
		log.Fatal("Error while creating producer", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
				Payload: []byte("1"),
			})

			if err != nil {
				fmt.Println("Failed to publish message", err)
			}

			time.Sleep(time.Millisecond * 250)
		}
	}()
}

func consume(wg sync.WaitGroup, client pulsar.Client, queue int) {
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            fmt.Sprintf("persistent://public/default/my-topic%d", queue),
		SubscriptionName: "my-sub",
		Type:             pulsar.Shared,
	})

	if err != nil {
		log.Fatal("Error while subscribe: ", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; true; i++ {
			msg, err := consumer.Receive(context.Background())
			if err != nil {
				log.Print("Error while receive: ", err)
				continue
			}

			if i%10 == 0 {
				fmt.Printf("Queue: %d, LedgerID: %v, EntryID: %v, PartitionID: %v\n",
					queue, msg.ID().LedgerID(), msg.ID().EntryID(), msg.ID().PartitionIdx())
			}

			time.Sleep(time.Millisecond * 500)

			consumer.Ack(msg)
		}
	}()
}

func createRacks() {
	bookies := map[string][]string{
		"rack1": {"bookie1:3181", "bookie2:3181", "bookie3:3181"},
		"rack2": {"bookie4:3181", "bookie5:3181", "bookie6:3181"},
		"rack3": {"bookie7:3181", "bookie8:3181", "bookie9:3181"},
	}

	client := &http.Client{}
	for rack, bks := range bookies {
		body, err := json.Marshal(map[string]string{"rack": rack})
		if err != nil {
			log.Fatalln(errors.Wrap(err, "failed to marshal request body"))
		}

		for _, bk := range bks {
			buf := bytes.NewBuffer(body)
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://broker2:8080/admin/v2/bookies/racks-info/%s?group=default", bk), buf)
			if err != nil {
				log.Fatalln(errors.Wrap(err, "error creating post request"))
			}

			req.Close = true
			req.Header.Set("Content-type", "application/json")
			req.Header.Set("Accept", "application/json")

			res, err := client.Do(req)
			if err != nil {
				log.Fatalln(errors.Wrap(err, "error sending post request"))
			}
			defer res.Body.Close()

			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalln(errors.Wrap(err, "error reading response body"))
			}

			log.Println(string(resBody))

			if res.StatusCode >= 300 {
				log.Fatalln("status code >= 300")
			}
		}
	}

}

func main() {
	createRacks()

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		//URL:               "pulsar://qaas-pulsar-perf01:6650,qaas-pulsar-perf02:6650,qaas-pulsar-perf03:6650,qaas-pulsar-perf04:6650,qaas-pulsar-perf05:6650,qaas-pulsar-perf06:6650",
		URL:               "pulsar://broker1:6650,broker2:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}

	defer client.Close()

	queueCount, err := getenvInt("QUEUE_COUNT")
	if err != nil {
		log.Fatalf("Could not get queue count: %v", err)
	}

	var wg sync.WaitGroup

	for queue := 0; queue < queueCount; queue++ {
		produce(wg, client, queue)
		consume(wg, client, queue)
	}

	wg.Add(1)
	wg.Wait()
}
