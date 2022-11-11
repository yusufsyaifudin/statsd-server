package main

import (
	"fmt"
	"github.com/DataDog/datadog-go/v5/statsd"
	"log"
)

func main() {
	ddClient, err := statsd.New("127.0.0.1:8125")
	defer func() {
		if ddClient == nil {
			return
		}

		if _err := ddClient.Close(); _err != nil {
			err = fmt.Errorf("error closing statsd client: %w", err)
			log.Fatal(_err)
		}
	}()

	if err != nil {
		err = fmt.Errorf("error initiate statsd client: %w", err)
		log.Fatal(err)
		return
	}

	for i := 0; i < 1000; i++ {
		err = ddClient.Count("my_team.my_counter", int64(i), []string{"test_tag", "another:tag", "other"}, 0.5)
		if err != nil {
			err = fmt.Errorf("error counter %d: %w", i, err)
			log.Fatal(err)
			return
		}

		err = ddClient.Flush()
		if err != nil {
			err = fmt.Errorf("error flush %d: %w", i, err)
			log.Fatal(err)
			return
		}
	}

}
