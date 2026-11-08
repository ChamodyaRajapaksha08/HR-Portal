package controllers

import (
	"strings"
	"sync"
	"time"
)

const (
	loginAttemptLimit  = 5
	loginAttemptWindow = 15 * time.Minute
)

var loginAttempts = struct {
	sync.Mutex
	byKey map[string]loginAttemptState
}{byKey: make(map[string]loginAttemptState)}

type loginAttemptState struct {
	failures []time.Time
	inFlight int
}

func loginAttemptKey(email, ip string) string {
	return strings.ToLower(strings.TrimSpace(email)) + "|" + ip
}

func reserveLoginAttempt(key string, now time.Time) bool {
	loginAttempts.Lock()
	defer loginAttempts.Unlock()
	state := loginAttempts.byKey[key]
	state.failures = activeAttempts(state.failures, now)
	if len(state.failures)+state.inFlight >= loginAttemptLimit {
		loginAttempts.byKey[key] = state
		return false
	}
	state.inFlight++
	loginAttempts.byKey[key] = state
	return true
}

func finishLoginAttempt(key string, now time.Time, failed, succeeded bool) {
	loginAttempts.Lock()
	defer loginAttempts.Unlock()
	state := loginAttempts.byKey[key]
	if state.inFlight > 0 {
		state.inFlight--
	}
	state.failures = activeAttempts(state.failures, now)
	if failed {
		state.failures = append(state.failures, now)
	}
	if succeeded {
		state.failures = nil
	}
	if len(state.failures) == 0 && state.inFlight == 0 {
		delete(loginAttempts.byKey, key)
	} else {
		loginAttempts.byKey[key] = state
	}
	for candidate, candidateState := range loginAttempts.byKey {
		if len(activeAttempts(candidateState.failures, now)) == 0 && candidateState.inFlight == 0 {
			delete(loginAttempts.byKey, candidate)
		}
	}
}

func activeAttempts(attempts []time.Time, now time.Time) []time.Time {
	cutoff := now.Add(-loginAttemptWindow)
	active := attempts[:0]
	for _, attempt := range attempts {
		if attempt.After(cutoff) {
			active = append(active, attempt)
		}
	}
	return active
}
