package repositories

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"user-service/config"
	"user-service/internal/core/domain"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	TestDb   *gorm.DB
	TestRepo *userRepository
	Cfg      *config.Config
	TestData struct {
		User domain.User
	}
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	cfgPath := "../../test/user.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database)
	db, err := gorm.Open(postgres.Open(dsn))
	db.Debug()

	if err != nil {
		panic(errors.WithStack(err))
	}

	repository, err := NewUserRepository(db)

	if err != nil {
		panic(errors.WithStack(err))
	}

	db.Exec("DELETE FROM public.users")

	db.Exec("INSERT INTO public.users (id, name, last_name, email) VALUES ('test-id', 'test-name', 'test-lastname', 'test@email.com')")

	suite.Cfg = cfg
	suite.TestDb = db
	suite.TestRepo = repository
	suite.TestData = struct {
		User domain.User
	}{
		User: domain.User{
			ID:       "test-id",
			Name:     "test-name",
			LastName: "test-lastname",
			Email:    "test@email.com",
		},
	}
}

func (suite *UserRepositoryTestSuite) TestRepository_Get() {
	result, err := suite.TestRepo.Get(context.Background(), suite.TestData.User.ID)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.User, result)
}

func (suite *UserRepositoryTestSuite) TestRepository_Get_NotFound() {
	_, err := suite.TestRepo.Get(context.Background(), "test")

	suite.Error(err)
}

func (suite *UserRepositoryTestSuite) TestRepository_Save() {
	newUser := suite.TestData.User
	newUser.Name = "test-name-3"
	newUser.ID = "test-id-3"

	_, err := suite.TestRepo.Save(context.Background(), newUser)

	suite.NoError(err)

	queryResult := domain.User{}
	suite.TestDb.Raw("SELECT * FROM public.users WHERE id=?",
		newUser.ID).Scan(&queryResult)

	suite.EqualValues(newUser.Name, queryResult.Name)
}

func (suite *UserRepositoryTestSuite) TestRepository_Update() {
	suite.TestDb.Exec("INSERT INTO public.users (id, name, last_name, email) VALUES ('test-id-2', 'test-name', 'test-lastname', 'test@email.com')")

	updated := suite.TestData.User
	updated.ID = "test-id-2"
	updated.Name = "test-name-3"

	_, err := suite.TestRepo.Update(context.Background(), updated)

	suite.NoError(err)

	queryResult := domain.User{}
	suite.TestDb.Raw("SELECT * FROM public.users WHERE id=?",
		updated.ID).Scan(&queryResult)

	suite.EqualValues(updated.Name, queryResult.Name)
}

func TestIntegration_UserRepositoryTestSuite(t *testing.T) {
	testSuite := new(UserRepositoryTestSuite)
	suite.Run(t, testSuite)
}
