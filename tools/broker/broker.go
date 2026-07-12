package broker

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Event interface {
	GetEventType() string
}

type EventEmitter interface {
	EmitEvent(event Event)
}

type EventProducer interface {
	GetInitEvents() []Event
}

var _ EventEmitter = (*Broker)(nil)

// NOTE(patrik): Based on: https://gist.github.com/Ananto30/8af841f250e89c07e122e2a838698246
type Broker struct {
	notifier chan Event

	newClients     chan chan Event
	closingClients chan chan Event
	clients        map[chan Event]struct{}

	producers []EventProducer
}

func NewBroker() *Broker {
	return &Broker{
		notifier:       make(chan Event, 1024),
		newClients:     make(chan chan Event),
		closingClients: make(chan chan Event),
		clients:        make(map[chan Event]struct{}),
	}
}

func (broker *Broker) RegisterProducer(p EventProducer) {
	broker.producers = append(broker.producers, p)
}

func (broker *Broker) Listen() {
	for {
		select {
		case c := <-broker.newClients:
			slog.Info("New client connected")
			broker.clients[c] = struct{}{}
		case c := <-broker.closingClients:
			if _, ok := broker.clients[c]; ok {
				slog.Info("Client disconnected")

				delete(broker.clients, c)
				close(c)
			}
		case event := <-broker.notifier:
			for c := range broker.clients {
				select {
				case c <- event:
				default:
					// Drop event for slow client instead of blocking
				}
			}
		}
	}
}

func (broker *Broker) Start() {
	go broker.Listen()
}

func (broker *Broker) EmitEvent(event Event) {
	broker.notifier <- event
}

var _ (Event) = (*ConnectedEvent)(nil)

type ConnectedEvent struct{}

func (c ConnectedEvent) GetEventType() string {
	return "connected"
}

func (broker *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	rc := http.NewResponseController(w)

	eventChan := make(chan Event, 16)
	broker.newClients <- eventChan
	defer func() {
		broker.closingClients <- eventChan
	}()

	sendEvent := func(eventData Event) {
		fmt.Fprintf(w, "event: %s\n", eventData.GetEventType())
		fmt.Fprintf(w, "data: ")

		encode := json.NewEncoder(w)
		encode.Encode(eventData)

		fmt.Fprintf(w, "\n\n")
		rc.Flush()
	}

	sendEvent(ConnectedEvent{})

	for _, producer := range broker.producers {
		for _, event := range producer.GetInitEvents() {
			sendEvent(event)
		}
	}

	for {
		select {
		case <-r.Context().Done():
			return

		case event := <-eventChan:
			sendEvent(event)
		}
	}
}
