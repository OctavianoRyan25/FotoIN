package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"fotoin/backend/internal/config"
)

// GoogleUserInfo is the data returned from Google's userinfo endpoint.
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// oauthState stores a random state string and its expiry to prevent CSRF.
type oauthState struct {
	value     string
	expiresAt time.Time
}

// Google wraps OAuth2 configuration and an in-memory CSRF state store.
type Google struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string

	mu     sync.Mutex
	states map[string]oauthState // state -> expiry
}

func NewGoogle(ctx context.Context, cfg config.GoogleConfig) *Google {
	g := &Google{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		states:       make(map[string]oauthState),
	}
	g.StartCleanupWorker(ctx, 10*time.Minute)
	return g
}

// CleanUp expired states periodically
func (g *Google) StartCleanupWorker(ctx context.Context, interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                return // Keluar secara rapi saat server shutdown
            case <-ticker.C:
                g.mu.Lock()
                now := time.Now()
                for state, data := range g.states {
                    if now.After(data.expiresAt) {
                        delete(g.states, state)
                    }
                }
                g.mu.Unlock()
            }
        }
    }()
}

// AuthURL builds the Google OAuth consent-screen URL and registers the state.
func (g *Google) AuthURL(state string) string {
	g.mu.Lock()
	g.states[state] = oauthState{value: state, expiresAt: time.Now().Add(10 * time.Minute)}
	g.mu.Unlock()

	params := fmt.Sprintf(
		"client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s&access_type=offline&prompt=select_account",
		g.ClientID,
		g.RedirectURL,
		"https://www.googleapis.com/auth/userinfo.email+https://www.googleapis.com/auth/userinfo.profile",
		state,
	)
	return "https://accounts.google.com/o/oauth2/v2/auth?" + params
}

// ValidateState checks the CSRF state token.
func (g *Google) ValidateState(state string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	s, ok := g.states[state]
	if !ok {
		return false
	}
	delete(g.states, state)
	return time.Now().Before(s.expiresAt)
}

// ExchangeCode trades the authorization code for user info.
func (g *Google) ExchangeCode(ctx context.Context, code string) (*GoogleUserInfo, error) {
	// 1. Exchange code for access token
	tokenURL := "https://oauth2.googleapis.com/token"
	body := fmt.Sprintf(
		"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code, g.ClientID, g.ClientSecret, g.RedirectURL,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL,
		newStringReader(body))
	if err != nil {
		return nil, fmt.Errorf("building token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exchanging code: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("google token error: %s", tokenResp.Error)
	}

	// 2. Fetch user info
	userReq, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	userReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		return nil, fmt.Errorf("fetching user info: %w", err)
	}
	defer userResp.Body.Close()

	rawBody, _ := io.ReadAll(userResp.Body)
	var userInfo GoogleUserInfo
	if err := json.Unmarshal(rawBody, &userInfo); err != nil {
		return nil, fmt.Errorf("decoding user info: %w", err)
	}
	if userInfo.Email == "" {
		return nil, fmt.Errorf("google did not return an email")
	}
	return &userInfo, nil
}

// newStringReader wraps a string body so we avoid importing strings repeatedly.
func newStringReader(s string) io.Reader {
	return &stringReader{s: s}
}

type stringReader struct {
	s   string
	pos int
}

func (r *stringReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.s) {
		return 0, io.EOF
	}
	n = copy(p, r.s[r.pos:])
	r.pos += n
	return n, nil
}
