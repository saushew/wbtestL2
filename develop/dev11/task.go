package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например)
	сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться
	на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

// Event .
type Event struct {
	UserID  int       `json:"user_id"`
	EventID int       `json:"event_id"`
	Date    time.Time `json:"date"`
	Content string    `json:"content"`
}

type (
	// EventI .
	EventI interface {
		Create(*Event) error
		Update(map[string]interface{}) (*Event, error)
		Get(int, time.Duration) ([]Event, error)
		Delete(int) error
	}

	// EventRepoI .
	EventRepoI interface {
		Store(*Event) error
		Modify(map[string]interface{}) (*Event, error)
		Delete(int) error
		Get(int, time.Duration) ([]Event, error)
	}
)

// EventUseCase .
type EventUseCase struct {
	repo EventRepoI
}

// NewUseCase .
func NewUseCase(r EventRepoI) *EventUseCase {
	return &EventUseCase{
		repo: r,
	}
}

// Create .
func (uc *EventUseCase) Create(event *Event) error {
	err := uc.repo.Store(event)
	if err != nil {
		return err
	}

	return nil
}

// Update .
func (uc *EventUseCase) Update(changes map[string]interface{}) (*Event, error) {
	event, err := uc.repo.Modify(changes)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// Get .
func (uc *EventUseCase) Get(userID int, period time.Duration) ([]Event, error) {
	events, err := uc.repo.Get(userID, period)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// Delete .
func (uc *EventUseCase) Delete(eventID int) error {
	err := uc.repo.Delete(eventID)
	if err != nil {
		return err
	}

	return nil
}

//~~~~~~~~~~~~~~~~~~~

// EventRepo .
type EventRepo struct {
	DB *sql.DB
}

// NewRepo .
func NewRepo(db *sql.DB) *EventRepo {
	return &EventRepo{db}
}

// Store .
func (r *EventRepo) Store(event *Event) error {
	if err := r.DB.QueryRow(
		`INSERT INTO events (user_id, date, content)
		VALUES ($1, $3, $4)
		RETURNING event_id`,
		event.UserID, event.Date, event.Content,
	).Scan(
		&event.EventID,
	); err != nil {
		return fmt.Errorf("EventRepo - Store - r.DB.QueryRow: %w", err)
	}

	return nil
}

// Modify .
func (r *EventRepo) Modify(changes map[string]interface{}) (*Event, error) {
	eventID := int(changes["event_id"].(float64))

	rows, err := r.DB.Query(
		`SELECT user_id, date, content
		FROM events
		WHERE event_id = $1`,
		eventID,
	)
	if err != nil {
		return nil, fmt.Errorf("EventRepo - Modify - r.DB.Query: %w", err)
	}

	if !rows.Next() {
		return nil, fmt.Errorf("Event does not exist, event_id = %d", eventID)
	}

	var _event Event
	if err := rows.Scan(
		&_event.UserID,
		&_event.Date,
		&_event.Content,
	); err != nil {
		return nil, fmt.Errorf("EventRepo - Modify - rows.Next - rows.Scan: %w", err)
	}

	_event.EventID = eventID

	rows.Close()

	if v, ok := changes["date"]; ok {
		_event.Date = v.(time.Time)
	}
	if v, ok := changes["user_id"]; ok {
		_event.UserID = int(v.(float64))
	}
	if v, ok := changes["content"]; ok {
		_event.Content = fmt.Sprintf("%v", v)
	}

	_, err = r.DB.Exec(
		`UPDATE events
		SET user_id = $1, date = $2, content = $3
		WHERE event_id = $4;`,
		_event.UserID, _event.Date, _event.Content, _event.EventID,
	)
	if err != nil {
		return nil, fmt.Errorf("EventRepo - Modify - r.DB.Exec: %w", err)
	}

	return &_event, nil
}

// Delete .
func (r *EventRepo) Delete(eventID int) error {
	res, err := r.DB.Exec(
		`DELETE FROM events
		WHERE event_id = $1`,
		eventID,
	)
	if err != nil {
		return fmt.Errorf("EventRepo - Delete - r.DB.Exec: %w", err)
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("Event does not exist, event_id = %d", eventID)
	}

	return nil
}

// Get .
func (r *EventRepo) Get(userID int, period time.Duration) ([]Event, error) {
	begin := time.Now().Truncate(24 * time.Hour)
	end := begin.Add(period)

	row := r.DB.QueryRow(
		`SELECT COUNT(*)
		FROM events
		WHERE user_id = $1 AND date >= $2 AND date < $3`,
		userID, begin, end,
	)

	var numOfRows int
	err := row.Scan(&numOfRows)
	if err != nil {
		return nil, fmt.Errorf("EventRepo - Get - row.Scan: %w", err)
	}

	rows, err := r.DB.Query(
		`SELECT *
		FROM events
		WHERE user_id = $1 AND date >= $2 AND date < $3`,
		userID, begin, end,
	)
	if err != nil {
		return nil, fmt.Errorf("EventRepo - Get - r.DB.Query: %w", err)
	}

	result := make([]Event, numOfRows)

	var event Event
	for i := 0; rows.Next(); i++ {
		if err := rows.Scan(
			&event.UserID,
			&event.EventID,
			&event.Date,
			&event.Content,
		); err != nil {
			return nil, fmt.Errorf("EventRepo - Get - rows.Next - rows.Scan: %w", err)
		}

		result[i] = event
	}

	return result, nil
}

//~~~~~~~~~~~~~~~~~~~

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 5 * time.Second
	shutdownTimeout = 3 * time.Second
)

// Server .
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// NewServer .
func NewServer(handler http.Handler, port string) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Addr:         port,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: shutdownTimeout,
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify .
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown .
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

//~~~~~~~~~~~~~~~~~~~

// ResponseRecorder .
type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func httpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		logger := log.Info()

		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}

		start := time.Now()
		handler.ServeHTTP(rec, req)

		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.Msgf("%s %s completed with %d %s in %v",
			req.Method,
			req.RequestURI,
			rec.StatusCode,
			http.StatusText(rec.StatusCode),
			time.Now().Sub(start),
		)
	})
}

//~~~~~~~~~~~~~~~~~~~

var (
	errorRequired          = errors.New("Required field missed")
	errorInvalidDateFormat = errors.New("Invalid date format")
)

// ValidateRequest .
func ValidateRequest(path string, req map[string]interface{}) (*Event, error) {
	var event Event
	switch path {
	case "/create_event":
		if _, ok := req["user_id"]; !ok {
			return nil, fmt.Errorf("%w: %s", errorRequired, "user_id")
		}
		if _, ok := req["date"]; !ok {
			return nil, fmt.Errorf("%w: %s", errorRequired, "date")
		}
	case "/update_event":
		if _, ok := req["event_id"]; !ok {
			return nil, fmt.Errorf("%w: %s", errorRequired, "event_id")
		}
	case "/delete_event":
		if _, ok := req["event_id"]; !ok {
			return nil, fmt.Errorf("%w: %s", errorRequired, "event_id")
		}
	default:
		if _, ok := req["user_id"]; !ok {
			return nil, fmt.Errorf("%w: %s", errorRequired, "user_id")
		}
	}

	if v, ok := req["user_id"]; ok {
		event.UserID = int(v.(float64))
	}

	if v, ok := req["event_id"]; ok {
		event.EventID = int(v.(float64))
	}

	if v, ok := req["date"]; ok {
		t, err := time.Parse("2006-01-02", fmt.Sprint(v))
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errorInvalidDateFormat, err)
		}
		event.Date = t
	}

	if v, ok := req["content"]; ok {
		event.Content = fmt.Sprint(v)
	}

	return &event, nil
}

//~~~~~~~~~~~~~~~~~~~

var (
	errorMethodNotAllowed = errors.New("Method not allowed")
)

// NewRouter .
func NewRouter(handler *http.ServeMux, e EventI) {
	newEventRoutes(handler, e)
}

type eventRoutes struct {
	eventUseCase EventI
}

func newEventRoutes(handler *http.ServeMux, e EventI) {
	er := &eventRoutes{e}

	handler.HandleFunc("/create_event", er.handleCreate())
	handler.HandleFunc("/update_event", er.handleUpdate())
	handler.HandleFunc("/delete_event", er.handleDelete())

	handler.HandleFunc("/events_for_day", er.handleGetEvents())
	handler.HandleFunc("/events_for_week", er.handleGetEvents())
	handler.HandleFunc("/events_for_month", er.handleGetEvents())
}

func (er *eventRoutes) handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			er.error(w, r, http.StatusMethodNotAllowed, errorMethodNotAllowed)
			return
		}

		req := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			er.error(w, r, http.StatusBadRequest, err)
			return
		}

		event, err := ValidateRequest(r.URL.Path, req)
		if err != nil {
			er.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := er.eventUseCase.Create(event); err != nil {
			er.error(w, r, http.StatusServiceUnavailable, err)
			return
		}

		er.success(w, r, http.StatusOK, event)
	}
}

func (er *eventRoutes) handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			er.error(w, r, http.StatusMethodNotAllowed, errorMethodNotAllowed)
			return
		}

		req := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			er.error(w, r, http.StatusBadRequest, err)
			return
		}

		event, err := er.eventUseCase.Update(req)
		if err != nil {
			er.error(w, r, http.StatusServiceUnavailable, err)
			return
		}

		er.respond(w, r, http.StatusOK, event)
	}
}

func (er *eventRoutes) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			er.error(w, r, http.StatusMethodNotAllowed, errorMethodNotAllowed)
			return
		}

		req := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			er.error(w, r, http.StatusBadRequest, err)
			return
		}

		event, err := ValidateRequest(r.URL.Path, req)
		if err != nil {
			er.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := er.eventUseCase.Delete(event.EventID); err != nil {
			er.error(w, r, http.StatusServiceUnavailable, err)
			return
		}

		er.success(w, r, http.StatusOK, fmt.Sprintf("event event_id = %d deleted", event.EventID))
	}
}

func (er *eventRoutes) handleGetEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			er.error(w, r, http.StatusMethodNotAllowed, errorMethodNotAllowed)
			return
		}

		var req map[string]interface{}
		req["user_id"] = r.URL.Query().Get("user_id")

		event, err := ValidateRequest(r.URL.Path, req)
		if err != nil {
			er.error(w, r, http.StatusBadRequest, err)
			return
		}

		var period time.Duration
		switch r.URL.Path {
		case "/events_for_day":
			period = 24 * time.Hour
		case "/events_for_week":
			period = 7 * 24 * time.Hour
		case "/events_for_month":
			period = 30 * 7 * 24 * time.Hour
		}

		events, err := er.eventUseCase.Get(event.UserID, period)
		if err != nil {
			er.error(w, r, http.StatusServiceUnavailable, err)
			return
		}

		er.success(w, r, http.StatusOK, events)
	}
}

func (er *eventRoutes) success(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	er.respond(w, r, code, map[string]interface{}{"result": data})
}

func (er *eventRoutes) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	er.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (er *eventRoutes) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger := log.Info()

	// sqlite
	schemaSQL := `CREATE TABLE IF NOT EXISTS events (
		event_id	INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id     INTEGER,
		date 		TIMESTAMP,
		content 	VARCHAR(128)
	);`

	sqlDB, err := sql.Open("sqlite3", "db-test")
	if err != nil {
		logger.Err(err)
		os.Exit(1)
	}

	if _, err = sqlDB.Exec(schemaSQL); err != nil {
		logger.Err(err)
		os.Exit(1)
	}

	// use case
	eventUseCase := NewUseCase(
		NewRepo(sqlDB),
	)

	mux := http.NewServeMux()
	NewRouter(mux, eventUseCase)
	handler := httpLogger(mux)
	httpServer := NewServer(handler, ":8080")

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Msg("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		logger.Err(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Err(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
