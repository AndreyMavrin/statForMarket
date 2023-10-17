package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"statForMarket/internal/model"
	"statForMarket/internal/repository"
	"statForMarket/internal/utils"
	"time"
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
		log.Println(err)
		http.Error(w, fmt.Sprintf("Ошибка при парсинге тела запроса: %s", err), http.StatusBadRequest)
		return
	}

	events := utils.GenerateEvents(count)
	if err := a.Repository.TestEvents(events); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Ошибка при добавлении событий: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *Application) Events(w http.ResponseWriter, r *http.Request) {
	filter := &model.EventFilter{
		EventType: r.URL.Query().Get("eventType"),
		From:      r.URL.Query().Get("from"),
		To:        r.URL.Query().Get("to"),
	}
	if filter.From != "" && filter.To == "" {
		filter.To = time.Now().Format("2006-01-02T15:04:05")
	} else if filter.From == "" && filter.To != "" {
		filter.From = ("1970-01-01T00:00:00")
	}
	if filter.From > filter.To {
		log.Println(errors.New("ошибка при указании времени фильтра"))
		http.Error(w, "Ошибка при указании времени фильтра", http.StatusBadRequest)
		return
	}

	events, err := a.Repository.Events(filter)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Ошибка при выборке событий: %s", err), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(events)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Ошибка при сериализации данных: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (a *Application) CreateEvent(w http.ResponseWriter, r *http.Request) {
	event := new(model.Event)
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при парсинге тела запроса: %s", err), http.StatusBadRequest)
		return
	}

	event.EventID = utils.GenerateEventID()
	if err := a.Repository.CreateEvent(event); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при добавлении события: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
