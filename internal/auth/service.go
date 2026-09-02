package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var ErrUsernameExists = errors.New("username already exists")
var ErrEmailExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid email or password")

type PlayerCreator interface {
	CreatePlayerTx(
		ctx context.Context,
		tx pgx.Tx,
		userID uuid.UUID,
	) error
}

type Service struct {
	repository    *Repository
	tokenService  *TokenService
	playerService PlayerCreator
	db            *pgxpool.Pool
}

func NewService(
	db *pgxpool.Pool,
	repository *Repository,
	tokenService *TokenService,
	playerService PlayerCreator,
) *Service {
	return &Service{
		db:            db,
		repository:    repository,
		tokenService:  tokenService,
		playerService: playerService,
	}
}

func (s *Service) RegisterUser(req RegisterRequest) error {
	err := ValidateRegister(req)
	if err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	}

	ctx := context.Background()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	err = s.repository.CreateUserTx(
		ctx,
		tx,
		user,
	)
	if err != nil {
		return err
	}

	err = s.playerService.CreatePlayerTx(
		ctx,
		tx,
		user.ID,
	)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) LoginUser(req LoginRequest) (User, error) {
	user, err := s.repository.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrInvalidCredentials
		}

		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		return User{}, ErrInvalidCredentials
	}

	return user, nil
}

func ValidateRegister(req RegisterRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	if len(req.Password) < 8 {
		return errors.New("password is too short")
	}

	return nil
}
