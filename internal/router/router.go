// Package router provides topic-based message routing, allowing messages
// to be dispatched to one or more named handlers based on topic patterns.
package router

import (
	"fmt"
	"regexp"
	"sync"
)

// Handler is a function that receives a routed message.
type Handler func(topic string, payload []byte)

// route holds a compiled pattern and its associated handler.
type route struct {
	name    string
	pattern *regexp.Regexp
	handler Handler
}

// Router dispatches messages to registered handlers whose topic pattern matches.
type Router struct {
	mu     sync.RWMutex
	routes []route
}

// New returns a new, empty Router.
func New() *Router {
	return &Router{}
}

// Register adds a named handler for topics matching the given regex pattern.
// Returns an error if the pattern is invalid or the name is already registered.
func (r *Router) Register(name, pattern string, h Handler) error {
	if name == "" {
		return fmt.Errorf("router: name must not be empty")
	}
	if h == nil {
		return fmt.Errorf("router: handler must not be nil")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("router: invalid pattern %q: %w", pattern, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, rt := range r.routes {
		if rt.name == name {
			return fmt.Errorf("router: handler %q already registered", name)
		}
	}

	r.routes = append(r.routes, route{name: name, pattern: re, handler: h})
	return nil
}

// Dispatch sends the message to all handlers whose pattern matches topic.
// Returns the number of handlers that were invoked.
func (r *Router) Dispatch(topic string, payload []byte) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, rt := range r.routes {
		if rt.pattern.MatchString(topic) {
			rt.handler(topic, payload)
			count++
		}
	}
	return count
}

// Len returns the number of registered routes.
func (r *Router) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.routes)
}
