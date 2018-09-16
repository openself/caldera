package provider

import (
	"context"

	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/contracts/events"
)

// Events defines data store Events provider methods
type Events interface {
	TransactProvider() (EventsTransact, error)
	Context(ctx context.Context) Events
	New(model *events.Event) (*events.Event, error)
	Find(id string) (*events.Event, error)
	FindByName(name string) ([]events.Event, error)
	List() ([]events.Event, error)
	Save(model *events.Event) (*events.Event, error)
	Delete(id string) error
	DeleteByName(name string) error
}

// EventsTransact allow transactions in provider
type EventsTransact interface {
	Transact
	Events
}
