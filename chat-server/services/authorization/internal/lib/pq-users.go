package lib

import (
	"database/sql"
	"errors"

	"github.com/gox7/notify/services/authorization/internal/repository"
	"github.com/gox7/notify/services/authorization/models"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthorizathionService struct {
		db *repository.Postgres
	}
)

func NewAuthorizathion(db *repository.Postgres, model *AuthorizathionService) {
	*model = AuthorizathionService{
		db: db,
	}
}

func (auth *AuthorizathionService) CreateUser(login, email, password, client string) (int64, error) {
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	lastId, err := auth.db.RegisterUser(login, email, string(bcryptHash), client)
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (auth *AuthorizathionService) SearchUser(email string, password string) (*models.PostgresUser, error) {
	user, err := auth.db.SearchUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (auth *AuthorizathionService) SearchUserByID(id int64) (*models.PostgresUser, error) {
	user, err := auth.db.SearchUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (auth *AuthorizathionService) RemoveUser(id int64) error {
	return auth.db.RemoveUserByID(id)
}
