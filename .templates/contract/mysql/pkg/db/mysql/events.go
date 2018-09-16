package mysql

import (
	"context"
	"database/sql"

	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/contracts/events"
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/db/provider"

	"github.com/satori/go.uuid"
)

type eventsProvider struct {
	*provider.Provider
}

func newEventsProvider(db *sql.DB) *eventsProvider {
	return &eventsProvider{Provider: provider.New(db)}
}

// Transaction returns provider with transaction
func (ep *eventsProvider) TransactProvider() (provider.EventsTransact, error) {
	p, err := ep.Provider.TransactProvider()
	if err != nil {
		return ep, err
	}
	return &eventsProvider{Provider: p}, nil
}

// Context returns provider with context
func (ep *eventsProvider) Context(ctx context.Context) provider.Events {
	return &eventsProvider{Provider: ep.Provider.Context(ctx)}
}

// New creates new Event object
func (ep *eventsProvider) New(model *events.Event) (*events.Event, error) {
	if model.Name == "" {
		return nil, provider.ErrNotDefinedName
	}
	model.ID = uuid.NewV4().String()
	stmt, err := ep.Prepare(queryInsertEvent)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(model.ID, model.Name)
	return model, err
}

// Find returns Event requested by ID
func (ep *eventsProvider) Find(id string) (*events.Event, error) {
	event := new(events.Event)
	row := ep.QueryRow(queryEventByID, id)
	return event, row.Scan(&event.ID, &event.Name)
}

// FindByName returns Events requested by Event name
func (ep *eventsProvider) FindByName(name string) ([]events.Event, error) {
	return ep.find(queryEventsByName, name)
}

// List returns all Event objects
func (ep *eventsProvider) List() ([]events.Event, error) {
	return ep.find(queryEvents)
}

// Save updates Event object
func (ep *eventsProvider) Save(model *events.Event) (*events.Event, error) {
	if model.ID == "" {
		return nil, provider.ErrNotDefinedID
	}
	if model.Name == "" {
		return nil, provider.ErrNotDefinedName
	}
	stmt, err := ep.Prepare(queryUpdateEvent)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(model.ID, model.Name)
	return model, err
}

// Delete removes Event object by ID
func (ep *eventsProvider) Delete(id string) error {
	if id == "" {
		return provider.ErrNotDefinedID
	}
	stmt, err := ep.Prepare(queryDeleteEventByID)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

// DeleteByName removes Event objects by Event name
func (ep *eventsProvider) DeleteByName(name string) error {
	if name == "" {
		return provider.ErrNotDefinedName
	}
	stmt, err := ep.Prepare(queryDeleteEventsByName)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(name)
	return err
}

func (ep *eventsProvider) find(query string, args ...interface{}) ([]events.Event, error) {
	items := make([]events.Event, 0)
	rows, err := ep.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := events.Event{}
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

const (
	queryEventByID    = `SELECT id, name FROM events WHERE id = $1`
	queryEventsByName = `SELECT id, name FROM events WHERE name = $1`
	queryEvents       = `SELECT id, name FROM events`
	queryInsertEvent  = `INSERT INTO events (id, name) VALUES ($1, $2)`
	queryUpdateEvent  = `INSERT INTO events (id, name) VALUES ($1, $2)
		ON DUPLICATE KEY UPDATE name = $2`
	queryDeleteEventByID    = `DELETE FROM events WHERE id = $1`
	queryDeleteEventsByName = `DELETE FROM events WHERE name = $1`
)
