package usecase

import (
	"context"
	"errors"
	"os"
	"time"

	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo domain.UserRepository
	google   *infra.Google
}

func NewAuthUsecase(userRepo domain.UserRepository, google *infra.Google) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo, google: google}
}

// ─── Input / Output types ────────────────────────────────────────────────────

type RegisterInput struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	// Optional: let the client specify their role (photographer or client).
	// Defaults to photographer so existing dashboard flow still works.
	Role string `json:"role"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        *domain.User `json:"user"`
}

// ─── Register ────────────────────────────────────────────────────────────────

func (u *AuthUsecase) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	existing, _ := u.userRepo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Determine role
	role := domain.RolePhotographer
	if input.Role == string(domain.RoleClient) {
		role = domain.RoleClient
	}

	user := &domain.User{
		ID:               uuid.New(),
		Name:             input.Name,
		Email:            input.Email,
		PasswordHash:     string(hashedPassword),
		Role:             role,
		SubscriptionTier: domain.TierBasic,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := generateJWT(user)
	if err != nil {
		return nil, err
	}

	user.HasPassword = user.PasswordHash != ""

	return &AuthResponse{AccessToken: token, User: user}, nil
}

// ─── Login ───────────────────────────────────────────────────────────────────

func (u *AuthUsecase) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.IsSuspended {
		return nil, errors.New("account has been suspended")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := generateJWT(user)
	if err != nil {
		return nil, err
	}

	user.HasPassword = user.PasswordHash != ""

	return &AuthResponse{AccessToken: token, User: user}, nil
}

// ─── Google OAuth ─────────────────────────────────────────────────────────────

// GoogleLoginURL returns the Google consent-screen URL.
// The client should redirect the browser to this URL.
func (u *AuthUsecase) GoogleLoginURL(state string) string {
	return u.google.AuthURL(state)
}

// GoogleCallback processes the OAuth2 callback code+state.
// It creates or fetches the user and returns a JWT — same as normal login.
func (u *AuthUsecase) GoogleCallback(ctx context.Context, code, state string) (*AuthResponse, error) {
	// 1. Validate CSRF state
	if !u.google.ValidateState(state) {
		return nil, errors.New("invalid or expired OAuth state")
	}

	// 2. Exchange authorization code for user info
	googleUser, err := u.google.ExchangeCode(ctx, code)
	if err != nil {
		return nil, errors.New("failed to authenticate with Google: " + err.Error())
	}

	if !googleUser.VerifiedEmail {
		return nil, errors.New("Google account email is not verified")
	}

	// 3. Find or create user
	user, err := u.userRepo.FindByEmail(ctx, googleUser.Email)
	if err != nil {
		// User does not exist — create one (role=client by default for Google sign-in)
		newUser := &domain.User{
			ID:               uuid.New(),
			Name:             googleUser.Name,
			Email:            googleUser.Email,
			PasswordHash:     "", // no password — Google-only account
			Role:             domain.RoleClient,
			SubscriptionTier: domain.TierBasic,
			IsVerified:       true,
		}
		if err := u.userRepo.Create(ctx, newUser); err != nil {
			return nil, errors.New("failed to create user account")
		}
		user = newUser
	}

	if user.IsSuspended {
		return nil, errors.New("account has been suspended")
	}

	// 4. Issue JWT
	token, err := generateJWT(user)
	if err != nil {
		return nil, err
	}

	user.HasPassword = user.PasswordHash != ""

	return &AuthResponse{AccessToken: token, User: user}, nil
}

// ─── JWT helper ──────────────────────────────────────────────────────────────

func generateJWT(user *domain.User) (string, error) {
	claims := &middleware.JWTClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("APP_SECRET")))
}
