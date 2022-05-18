package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"user-service/internal/core/domain"
)

type userRepository struct {
	Connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) (*userRepository, error) {
	err := db.AutoMigrate(&domain.User{})

	if err != nil {
		return nil, err
	}

	database := userRepository{
		Connection: db,
	}

	return &database, nil
}

func (repository *userRepository) Get(ctx context.Context, id string) (domain.User, error) {
	var user domain.User

	repository.Connection.WithContext(ctx).Preload(clause.Associations).First(&user, "id = ?", id)

	if (user == domain.User{}) {
		return user, errors.New("user not found")
	}

	return user, nil
}

func (repository *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User

	repository.Connection.WithContext(ctx).Find(&users)

	return users, nil
}

func (repository *userRepository) Save(ctx context.Context, user domain.User) (domain.User, error) {
	result := repository.Connection.WithContext(ctx).Create(&user)

	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}

func (repository *userRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	result := repository.Connection.WithContext(ctx).Model(&user).Updates(user)

	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}
