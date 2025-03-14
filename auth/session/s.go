package session

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	gSessions "github.com/gorilla/sessions"
)

const (
	errorFormat = "[sessions] ERROR! %s\n"
)

// Implement the sessions.session
type session struct {
	name    string
	request *http.Request
	store   sessions.Store
	session *gSessions.Session
	written bool
	writer  http.ResponseWriter
}

func (s *session) ID() string    { return s.Session().ID }
func (s *session) Written() bool { return s.written }

func (s *session) Get(key any) any { return s.Session().Values[key] }

func (s *session) Set(key any, val any) {
	s.Session().Values[key] = val
	s.written = true
}

func (s *session) Delete(key any) {
	delete(s.Session().Values, key)
	s.written = true
}

func (s *session) Clear() {
	s.Session().Values = make(map[any]any)
}

func (s *session) AddFlash(value any, vars ...string) {
	s.Session().AddFlash(value, vars...)
	s.written = true
}

func (s *session) Flashes(vars ...string) []any {
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *session) Options(options sessions.Options) {
	s.written = true
	s.Session().Options = options.ToGorillaOptions()
}

func (s *session) Save() error {
	if s.Written() {
		e := s.Session().Save(s.request, s.writer)
		if e == nil {
			s.written = false
		}
		return e
	}
	return nil
}

func (s *session) Session() *gSessions.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.request, s.name)
		if err != nil {
			log.Printf(errorFormat, err)
		}
	}
	return s.session
}
