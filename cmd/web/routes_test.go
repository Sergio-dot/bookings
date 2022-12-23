package main

import (
	"testing"

	"github.com/Sergio-dot/bookings/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing: test passed
	default:
		t.Fatalf("unexpected type %T", v)
	}
}
