package sikachal

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

var _ http.Handler = &Server{}

type Server struct {
	logger     *slog.Logger
	router     *httprouter.Router
	db         *DB
	corsOrigin string
}

func NewServer(logger *slog.Logger, db *DB, corsOrigin string) *Server {
	s := &Server{
		logger:     logger,
		db:         db,
		corsOrigin: corsOrigin,
		router:     httprouter.New(),
	}
	s.init()
	return s
}
func (s *Server) init() {
	s.router.GET("/user/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id, err := strconv.ParseUint(params.ByName("id"), 10, 64)
		if err != nil {
			s.logger.Debug("failed to parse user id", err, "id", id)
			http.Error(writer, "bad request", http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(request.Context(), time.Second*5)
		defer cancel()
		request = request.WithContext(ctx)
		user, err := s.db.GetUserByID(request.Context(), id)
		if err != nil {
			s.logger.Warn("failed to get user by id", err, "id", id)
			http.Error(writer, "internal error", http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(writer, "user not found", http.StatusNotFound)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(user)
	})
	s.router.PanicHandler = func(writer http.ResponseWriter, request *http.Request, err any) {
		s.logger.Error("unhandled error", err)
		http.Error(writer, "internal error", http.StatusInternalServerError)
	}
	s.router.NotFound = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		s.logger.Debug("unhandled path", request.Method, request.URL)
		http.Error(writer, "not found", http.StatusNotFound)
	})
	s.router.GlobalOPTIONS = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Headers", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "*")
		writer.Header().Set("Access-Control-Allow-Origin", s.corsOrigin)
		writer.WriteHeader(http.StatusNoContent)
	})
	s.router.MethodNotAllowed = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		s.logger.Debug("unhandled method", request.Method, request.URL)
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	})
}
func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(res, req)
}
