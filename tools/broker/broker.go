package broker

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nrednav/cuid2"
)

var createClientId, _ = cuid2.Init(cuid2.WithLength(32))

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

type client struct {
	id     string
	userId string
	events chan Event
}

// NOTE(patrik): Based on: https://gist.github.com/Ananto30/8af841f250e89c07e122e2a838698246
type Broker struct {
	notifier chan Event

	newClients     chan *client
	closingClients chan *client
	clients        map[*client]struct{}

	producers []EventProducer
}

func NewBroker() *Broker {
	return &Broker{
		notifier:       make(chan Event, 1024),
		newClients:     make(chan *client),
		closingClients: make(chan *client),
		clients:        make(map[*client]struct{}),
	}
}

func (broker *Broker) RegisterProducer(p EventProducer) {
	broker.producers = append(broker.producers, p)
}

func (broker *Broker) Listen() {
	for {
		select {
		case c := <-broker.newClients:
			slog.Info("New client connected", slog.String("userId", c.userId))
			broker.clients[c] = struct{}{}
		case c := <-broker.closingClients:
			if _, ok := broker.clients[c]; ok {
				slog.Info("Client disconnected", slog.String("userId", c.userId))

				delete(broker.clients, c)
				close(c.events)
			}
		case event := <-broker.notifier:
			for c := range broker.clients {
				select {
				case c.events <- event:
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

func (broker *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request, userId string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	rc := http.NewResponseController(w)

	c := &client{
		id:     createClientId(),
		userId: userId,
		events: make(chan Event, 16),
	}
	broker.newClients <- c
	defer func() {
		broker.closingClients <- c
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

		case event := <-c.events:
			sendEvent(event)
		}
	}
}
