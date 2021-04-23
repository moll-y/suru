package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const timeToSleep = 10 * time.Second

type item struct {
	Name   string `json:"name"`
	OnSale bool   `json:"onSale"`
}

type bestBuyService struct {
	apiURL string
	broker *MessageBroker
}

func NewBestBuyService(apiURL string, messageBroker *MessageBroker) *bestBuyService {
	return &bestBuyService{
		broker: messageBroker,
		apiURL: apiURL,
	}
}

func (b *bestBuyService) Run() {
	go b.broker.Dispatcher()
	var item item
	for {
		func() {
			resp, err := http.Get(b.apiURL)
			if err != nil {
				b.broker.SendMessage("something was wrong")
				log.Println(err.Error())
				return
			}
			defer resp.Body.Close()
			err = json.NewDecoder(resp.Body).Decode(&item)
			if err != nil {
				b.broker.SendMessage("something was wrong")
				log.Println(err.Error())
				return
			}
			if item.OnSale {
				b.broker.SendMessage("itme on stock")
			}
			log.Printf("name: %s onSale: %t\n", item.Name, item.OnSale)
		}()
		time.Sleep(timeToSleep)
	}
}