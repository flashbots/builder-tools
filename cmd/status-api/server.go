package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/flashbots/builder-tools/common"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var MaxEvents = 100

type HTTPServerConfig struct {
	ListenAddr   string
	Log          *slog.Logger
	PipeFilename string
	EnablePprof  bool

	DrainDuration            time.Duration
	GracefulShutdownDuration time.Duration
	ReadTimeout              time.Duration
	WriteTimeout             time.Duration
}

type Event struct {
	ReceivedAt time.Time `json:"received_at"`
	Message    string    `json:"message"`
}

type Server struct {
	cfg *HTTPServerConfig
	log *slog.Logger

	srv *http.Server

	events     []Event
	eventsLock sync.RWMutex
}

func NewServer(cfg *HTTPServerConfig) (srv *Server, err error) {
	srv = &Server{
		cfg:    cfg,
		log:    cfg.Log,
		srv:    nil,
		events: make([]Event, 0),
	}

	if cfg.PipeFilename != "" {
		os.Remove(cfg.PipeFilename)
		err := syscall.Mknod(cfg.PipeFilename, syscall.S_IFIFO|0o666, 0)
		if err != nil {
			return nil, err
		}

		go srv.readPipeInBackground()
	}

	mux := chi.NewRouter()
	mux.With(srv.httpLogger).Get("/", srv.handleLivenessCheck)
	mux.With(srv.httpLogger).Get("/livez", srv.handleLivenessCheck)
	mux.With(srv.httpLogger).Get("/api/v1/new_event", srv.handleNewEvent)
	mux.With(srv.httpLogger).Get("/api/v1/events", srv.handleGetEvents)

	if cfg.EnablePprof {
		srv.log.Info("pprof API enabled")
		mux.Mount("/debug", middleware.Profiler())
	}

	srv.srv = &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return srv, nil
}

func (s *Server) readPipeInBackground() {
	file, err := os.OpenFile(s.cfg.PipeFilename, os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		s.log.Error("Open named pipe file error:", "error", err)
		return
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadBytes('\n')
		if err == nil {
			msg := strings.Trim(string(line), "\n")

			s.log.Info("Received new event via pipe", "message", msg)
			s.addEvent(Event{
				ReceivedAt: time.Now().UTC(),
				Message:    msg,
			})
		}
	}
}

func (s *Server) httpLogger(next http.Handler) http.Handler {
	return common.LoggingMiddlewareSlog(s.log, next)
}

func (s *Server) Start() {
	s.log.Info("Starting HTTP server", "listenAddress", s.cfg.ListenAddr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Error("HTTP server failed", "err", err)
	}
}

func (s *Server) handleLivenessCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) addEvent(event Event) {
	// Add event to the list and prune if necessary
	s.eventsLock.Lock()
	s.events = append(s.events, event)
	if len(s.events) > MaxEvents {
		s.events = s.events[1:]
	}
	s.eventsLock.Unlock()
}

func (s *Server) handleNewEvent(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("message")
	s.log.Info("Received new event", "message", msg)
	s.addEvent(Event{
		ReceivedAt: time.Now().UTC(),
		Message:    msg,
	})
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetEvents(w http.ResponseWriter, r *http.Request) {
	s.eventsLock.RLock()
	defer s.eventsLock.RUnlock()

	// write events as JSON response
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(s.events)
	if err != nil {
		s.log.Error("Failed to encode events", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
