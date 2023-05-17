package game

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"veriDrinkApi/gender"
)

// charset is the set of characters used when generating a session ID.
const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

// Session represents a game session.
type Session struct {
	Id         string     `json:"id"`         // Unique identifier for the session
	Difficulty int        `json:"difficulty"` // Difficulty of the session
	Owner      string     `json:"owner"`      // Name of the session owner
	Players    []*Player  `json:"players"`    // List of players in the session
	queue      []*Player  // Queue of players for random selection
	mu         sync.Mutex // Mutex to prevent data races
}

// SessionManager manages all game sessions.
type SessionManager struct {
	sessions []*Session // List of sessions
	mu       sync.Mutex // Mutex to prevent data races
}

// sessionManager is a global SessionManager instance.
var sessionManager *SessionManager

// newSessionManager creates a new SessionManager.
func newSessionManager() *SessionManager {
	return &SessionManager{
		sessions: []*Session{},
	}
}

// GetSessionManager returns the global SessionManager instance, creating it if necessary.
func GetSessionManager() *SessionManager {
	if sessionManager == nil {
		sessionManager = newSessionManager()
	}
	return sessionManager
}

// NewSession creates a new session with the given owner and adds it to the session manager.
func (sm *SessionManager) NewSession(owner string) (*Session, error) {
	id, err := generateId(6)
	if err != nil {
		return nil, err
	}

	session := &Session{
		Id:      id,
		Owner:   owner,
		Players: []*Player{},
	}

	sm.mu.Lock()
	sm.sessions = append(sm.sessions, session)
	sm.mu.Unlock()

	return session, nil
}

// generateId generates a random session ID of the given length.
func generateId(n int) (string, error) {
	result := make([]byte, n)
	for i := range result {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[index.Int64()]
	}
	return string(result), nil
}

// FindSessionById finds a session with the given ID.
func (sm *SessionManager) FindSessionById(sessionId string) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return sm.findSessionById(sessionId)
}

// findSessionById is an unexported helper method that finds a session with the given ID.
func (sm *SessionManager) findSessionById(sessionId string) (*Session, error) {
	for _, session := range sm.sessions {
		if session.Id == sessionId {
			return session, nil
		}
	}

	return nil, fmt.Errorf("session not found")
}

// AddPlayer adds a player to the session and the queue.
func (session *Session) AddPlayer(player *Player) {
	session.mu.Lock()
	defer session.mu.Unlock()

	if player.Preference == gender.Female {

	}
	session.Players = append(session.Players, player)
	session.queue = append(session.queue, player)
}

// RemovePlayer removes a player from the session and the queue.
func (session *Session) RemovePlayer(playerName string) error {
	session.mu.Lock()
	defer session.mu.Unlock()

	playerIndex := -1
	for i, player := range session.Players {
		if player.Name == playerName {
			playerIndex = i
			break
		}
	}

	if playerIndex == -1 {
		return errors.New("player not found")
	}

	session.Players = append(session.Players[:playerIndex], session.Players[playerIndex+1:]...)

	queueIndex := -1
	for i, player := range session.queue {
		if player.Name == playerName {
			queueIndex = i
			break
		}
	}

	if queueIndex != -1 {
		session.queue = append(session.queue[:queueIndex], session.queue[queueIndex+1:]...)
	}

	return nil
}

// randomPlayer returns a random player from the session, ensuring that each player is selected once before any player is selected again.
func (session *Session) randomPlayer() (*Player, error) {
	session.mu.Lock()
	defer session.mu.Unlock()

	// If the queue is empty, fill it with the players
	if len(session.queue) == 0 {
		if len(session.Players) == 0 {
			return nil, errors.New("no players in the session")
		}

		session.queue = make([]*Player, len(session.Players))
		copy(session.queue, session.Players)
	}

	// Choose a random player from the queue
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(session.queue))))
	if err != nil {
		return nil, err
	}
	i := index.Int64()

	// Remove the chosen player from the queue
	player := session.queue[i]
	session.queue = append(session.queue[:i], session.queue[i+1:]...)

	return player, nil
}

// getEligiblePlayers returns players who are not in the exclude list and satisfy a given condition.
func (session *Session) getEligiblePlayers(exclude []*Player, condition func(*Player) bool) ([]*Player, error) {
	session.mu.Lock()
	defer session.mu.Unlock()

	eligiblePlayers := make([]*Player, 0)
	for _, player := range session.Players {
		// Check if player is in the exclude list
		excluded := false
		for _, ex := range exclude {
			if player.Name == ex.Name {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		if condition(player) {
			eligiblePlayers = append(eligiblePlayers, player)
		}
	}

	if len(eligiblePlayers) == 0 {
		return nil, errors.New("no eligible players found")
	}

	return eligiblePlayers, nil
}

// randomPlayerFromList returns a random player from a given list of players.
func (session *Session) randomPlayerFromList(players []*Player) (*Player, error) {
	session.mu.Lock()
	defer session.mu.Unlock()

	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(players))))
	if err != nil {
		return nil, err
	}
	i := index.Int64()

	return players[i], nil
}

// NextRound replaces certain symbols in a string with player names.
func (session *Session) NextRound(input string) (string, error) {
	session.mu.Lock()
	defer session.mu.Unlock()

	currentPlayer, err := session.randomPlayer()
	if err != nil {
		return "", err
	}

	excludedPlayers := []*Player{currentPlayer}

	replacerFunc := func(symbol string, gender string) (string, error) {
		eligiblePlayers, err := session.getEligiblePlayers(excludedPlayers, func(p *Player) bool {
			return p.Gender == gender
		})
		if err != nil {
			return "", err
		}

		randomPlayer, err := session.randomPlayerFromList(eligiblePlayers)
		if err != nil {
			return "", err
		}

		excludedPlayers = append(excludedPlayers, randomPlayer)

		return strings.Replace(input, symbol, randomPlayer.Name, 1), nil
	}

	input, err = replacerFunc(":&", currentPlayer.Gender)
	if err != nil {
		return "", err
	}

	input, err = replacerFunc(":@", "")
	if err != nil {
		return "", err
	}

	input, err = replacerFunc(":@o", gender.Male)
	if currentPlayer.Gender == gender.Male {
		input, err = replacerFunc(":@o", gender.Female)
	}
	if err != nil {
		return "", err
	}

	input, err = replacerFunc(":@s", currentPlayer.Gender)
	if err != nil {
		return "", err
	}

	input, err = replacerFunc(":@a", currentPlayer.Preference)
	if err != nil {
		return "", err
	}

	return input, nil
}
