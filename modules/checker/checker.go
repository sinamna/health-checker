package checker

import (
	"fmt"
	"health_checker/pkg/model"
	"health_checker/pkg/repository"
	"net/http"
	"time"
)

type Checker struct {
	WorkerNum int
	taskChan  chan *model.EndpointResponse
	resetChan chan *model.EndpointResponse
	interval  int
}

func NewChecker(workerNum int, interval int) *Checker {
	checker := &Checker{
		WorkerNum: workerNum,
		taskChan:  make(chan *model.EndpointResponse, 100),
		resetChan: make(chan *model.EndpointResponse, 100),
		interval:  interval,
	}
	return checker
}

func (c *Checker) Start() {
	fmt.Println("starting ")
	ticker := time.NewTicker(time.Duration(c.interval) * time.Second)
	go c.retrieve(ticker)
	for i := 0; i < c.WorkerNum; i++ {
		go c.work()
	}

	go c.alerter(ticker)
	for i := 0; i < c.WorkerNum; i++ {
		go c.resetAndAlert()
	}
}
func (c *Checker) work() {
	timeout := 2 * time.Second
	client := http.Client{
		Timeout: timeout,
	}
	for endpoint := range c.taskChan {
		resp, err := client.Get("https://" + endpoint.Url)

		var result string
		if err == nil && (resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 202) {
			fmt.Println(endpoint.Url + " is healthy")
			result = repository.Success
		} else {
			fmt.Println(endpoint.Url + " is sick")
			result = repository.Failed
			fmt.Println(err.Error())
		}

		err = repository.Database.UpdateEndpointResultByOne(endpoint.Id, result)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (c *Checker) resetAndAlert() {
	for endpoint := range c.resetChan {
		fmt.Println("alerting on " + endpoint.Url)
		_, err := repository.Database.CreateAlert(endpoint.Id, "Endpoint wasn't healthy.")
		if err != nil {
			fmt.Println(err.Error())
		} else {
			err = repository.Database.ResetEndpointFailed(endpoint.Id)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
func (c *Checker) retrieve(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			endpoints, err := repository.Database.GetAllEndpoints()
			if err != nil {
				fmt.Println(err.Error())
			}
			for _, endpoint := range endpoints {
				c.taskChan <- endpoint
			}
		}
	}
}
func (c *Checker) alerter(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			endpoints, err := repository.Database.GetEndpointsByThresholdCrossed()
			if err != nil {
				fmt.Println(err.Error())
			}
			for _, endpoint := range endpoints {
				c.resetChan <- endpoint
			}
		}
	}
}
