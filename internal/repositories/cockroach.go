package repositories

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"user-service/internal/core/domain"
)

type cockroachdb struct {
	Connection *gorm.DB
}

func NewCockroachDB(connStr string) (*cockroachdb, error) {
	db, err := gorm.Open(postgres.Open(connStr))

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.User{})

	if err != nil {
		return nil, err
	}

	database := cockroachdb{
		Connection: db,
	}

	return &database, nil
}

func (repository *cockroachdb) Get(id string) (domain.User, error) {
	var user domain.User

	repository.Connection.Preload(clause.Associations).First(&user, "id = ?", id)

	if (user == domain.User{}) {
		return user, errors.New("user not found")
	}

	return user, nil
}

func (repository *cockroachdb) GetAll() ([]domain.User, error) {
	var users []domain.User

	repository.Connection.Find(&users)

	return users, nil
}

func (repository *cockroachdb) Save(user domain.User) (domain.User, error) {
	result := repository.Connection.Create(&user)

	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}

func (repository *cockroachdb) Update(user domain.User) (domain.User, error) {
	result := repository.Connection.Model(&user).Updates(user)

	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}
