package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"statForMarket/internal/model"
	"statForMarket/internal/repository"
)

type Application struct {
	Repository *repository.Repository
}

func (a *Application) Run(ctx context.Context) {
	repoConn := repository.InitRepository()
	a.Repository = &repository.Repository{Conn: repoConn}
}

func (a *Application) TestEvents(w http.ResponseWriter, r *http.Request) {
	var count int
	if err := json.NewDecoder(r.Body).Decode(&count); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при парсинге тела запроса: %s", err), http.StatusBadRequest)
		return
	}

	events := generateEvents(count)

	if err := a.Repository.TestEvents(events); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при добавлении события: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *Application) Events(w http.ResponseWriter, r *http.Request) {
	event := new(model.Event)
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при парсинге тела запроса: %s", err), http.StatusBadRequest)
		return
	}

	if err := a.Repository.CreateEvent(event); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при добавлении события: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *Application) CreateEvent(w http.ResponseWriter, r *http.Request) {
	event := new(model.Event)
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при парсинге тела запроса: %s", err), http.StatusBadRequest)
		return
	}

	if err := a.Repository.CreateEvent(event); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при добавлении события: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func generateEvents(count int) []*model.Event {
	events := make([]*model.Event, count)

	for i := 0; i < count; i++ {
		event := new(model.Event)

		events[i] = event
	}

	return events

}
