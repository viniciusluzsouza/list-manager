package user_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/user"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	defaultExpectedInsertQuery = "INSERT INTO `user`"
	defaultExpectedUpdateQuery = "UPDATE `user` SET"
	defaultExpectedGetQuery    = "SELECT (.+) FROM `user`"
	defaultExpectedExistsQuery = "SELECT count(*) > 0 FROM `user`"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	t          *testing.T
	repository user.Repository
	sqlMock    sqlmock.Sqlmock
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	var db *sql.DB
	var err error

	db, mock, err := sqlmock.New()
	assert.Nil(s.t, err)

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	assert.Nil(s.t, err)

	s.repository = user.NewRepository(gdb)
	s.t = s.T()
	s.sqlMock = mock
}

type UserRepositorySaveTestSuite struct {
	UserRepositoryTestSuite
}

func TestUserRepositorySaveTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySaveTestSuite))
}

func (s UserRepositorySaveTestSuite) TestUserSaveSuccess() {
	user := getUserToSaveTest()

	s.sqlMock.ExpectBegin()
	s.sqlMock.ExpectExec(defaultExpectedInsertQuery).
		WithArgs(user.Name, user.Email, user.Login, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.sqlMock.ExpectCommit()

	err := s.repository.Save(&user)

	assert.Nil(s.t, err)
	assert.Equal(s.t, uint64(1), user.ID)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func (s UserRepositorySaveTestSuite) TestUserSaveError() {
	user := getUserToSaveTest()
	expectedError := errors.New("error")

	s.sqlMock.ExpectBegin()
	s.sqlMock.ExpectExec(defaultExpectedInsertQuery).
		WithArgs(user.Name, user.Email, user.Login, user.Password).
		WillReturnError(expectedError)
	s.sqlMock.ExpectRollback()

	err := s.repository.Save(&user)

	assert.NotNil(s.t, err)
	assert.Equal(s.t, expectedError, err)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

type UserRepositoryGetTestSuite struct {
	UserRepositoryTestSuite
}

func TestUserRepositoryGetTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryGetTestSuite))
}

func (s UserRepositoryGetTestSuite) TestUserGetSuccess() {
	user := getUserToTest()

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "login", "password",
	}).AddRow(
		user.ID, user.Name, user.Email, user.Login, user.Password,
	)
	s.sqlMock.ExpectQuery(defaultExpectedGetQuery).WithArgs(user.ID).WillReturnRows(rows)

	dbUser, err := s.repository.Get(user.ID)

	assert.Nil(s.t, err)
	assert.Equal(s.t, &user, dbUser)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func (s UserRepositoryGetTestSuite) TestUserGetNotFoundError() {
	id := uint64(1)

	emptyRows := sqlmock.NewRows([]string{"id", "name", "email", "login", "password"})
	s.sqlMock.ExpectQuery(defaultExpectedGetQuery).WithArgs(id).WillReturnRows(emptyRows)

	dbUser, err := s.repository.Get(id)

	assert.Nil(s.t, err)
	assert.Nil(s.t, dbUser)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func (s UserRepositoryGetTestSuite) TestUserGetError() {
	id := uint64(1)
	expectedError := errors.New("error")

	s.sqlMock.ExpectQuery(defaultExpectedGetQuery).WithArgs(id).WillReturnError(expectedError)

	_, err := s.repository.Get(id)

	assert.NotNil(s.t, err)
	assert.Equal(s.t, expectedError, err)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func (s UserRepositoryGetTestSuite) TestUserGetByLoginSuccess() {
	user := getUserToTest()

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "login", "password",
	}).AddRow(
		user.ID, user.Name, user.Email, user.Login, user.Password,
	)
	s.sqlMock.ExpectQuery(defaultExpectedGetQuery).WithArgs(user.Login).WillReturnRows(rows)

	dbUser, err := s.repository.GetByLogin(user.Login)

	assert.Nil(s.t, err)
	assert.Equal(s.t, &user, dbUser)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func (s UserRepositoryGetTestSuite) TestUserGetByLoginNotFoundError() {
	login := "test"

	emptyRows := sqlmock.NewRows([]string{"id", "name", "email", "login", "password"})
	s.sqlMock.ExpectQuery(defaultExpectedGetQuery).WithArgs(login).WillReturnRows(emptyRows)

	dbUser, err := s.repository.GetByLogin(login)

	assert.Nil(s.t, err)
	assert.Nil(s.t, dbUser)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func (s UserRepositoryGetTestSuite) TestUserGetByLoginError() {
	login := "test"
	expectedError := errors.New("error")

	s.sqlMock.ExpectQuery(defaultExpectedGetQuery).WithArgs(login).WillReturnError(expectedError)

	_, err := s.repository.GetByLogin(login)

	assert.NotNil(s.t, err)
	assert.Equal(s.t, expectedError, err)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

type UserRepositoryUpdateTestSuite struct {
	UserRepositoryTestSuite
}

func TestUserRepositoryUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryUpdateTestSuite))
}

func (s UserRepositoryUpdateTestSuite) TestUserUpdateSuccess() {
	user := getUserToTest()

	s.sqlMock.ExpectBegin()
	s.sqlMock.ExpectExec(defaultExpectedUpdateQuery).
		WithArgs(user.Name, user.Email, user.Login, user.Password, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.sqlMock.ExpectCommit()

	err := s.repository.Update(&user)

	assert.Nil(s.t, err)
	assert.Equal(s.t, uint64(1), user.ID)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func (s UserRepositoryUpdateTestSuite) TestUserUpdateError() {
	user := getUserToTest()
	expectedError := errors.New("error")

	s.sqlMock.ExpectBegin()
	s.sqlMock.ExpectExec(defaultExpectedUpdateQuery).
		WithArgs(user.Name, user.Email, user.Login, user.Password, user.ID).
		WillReturnError(expectedError)
	s.sqlMock.ExpectRollback()

	err := s.repository.Update(&user)

	assert.NotNil(s.t, err)
	assert.Equal(s.t, expectedError, err)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

type UserRepositoryExistsTestSuite struct {
	UserRepositoryTestSuite
}

func TestUserRepositoryExistsTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryExistsTestSuite))
}

func (s UserRepositoryExistsTestSuite) TestUserExistsSuccess() {
	user := getUserToTest()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	s.sqlMock.ExpectQuery(regexp.QuoteMeta(defaultExpectedExistsQuery)).WithArgs(user.ID).WillReturnRows(rows)

	exists, err := s.repository.Exists(user.ID)

	assert.Nil(s.t, err)
	assert.Equal(s.t, true, exists)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func (s UserRepositoryExistsTestSuite) TestUserExistsNotExistsSuccess() {
	user := getUserToTest()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	s.sqlMock.ExpectQuery(regexp.QuoteMeta(defaultExpectedExistsQuery)).WithArgs(user.ID).WillReturnRows(rows)

	exists, err := s.repository.Exists(user.ID)

	assert.Nil(s.t, err)
	assert.Equal(s.t, false, exists)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func (s UserRepositoryExistsTestSuite) TestUserExistsError() {
	user := getUserToTest()

	expectedError := errors.New("error")
	s.sqlMock.ExpectQuery(regexp.QuoteMeta(defaultExpectedExistsQuery)).WithArgs(user.ID).WillReturnError(expectedError)

	exists, err := s.repository.Exists(user.ID)

	assert.NotNil(s.t, err)
	assert.Equal(s.t, false, exists)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func (s UserRepositoryExistsTestSuite) TestUserExistsByLoginSuccess() {
	login := "test"

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	s.sqlMock.ExpectQuery(regexp.QuoteMeta(defaultExpectedExistsQuery)).WithArgs(login).WillReturnRows(rows)

	exists, err := s.repository.ExistsByLogin(login)

	assert.Nil(s.t, err)
	assert.Equal(s.t, true, exists)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func (s UserRepositoryExistsTestSuite) TestUserExistsByLoginNotExistsSuccess() {
	login := "test"

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	s.sqlMock.ExpectQuery(regexp.QuoteMeta(defaultExpectedExistsQuery)).WithArgs(login).WillReturnRows(rows)

	exists, err := s.repository.ExistsByLogin(login)

	assert.Nil(s.t, err)
	assert.Equal(s.t, false, exists)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func (s UserRepositoryExistsTestSuite) TestUserExistsByLoginError() {
	login := "test"

	expectedError := errors.New("error")
	s.sqlMock.ExpectQuery(regexp.QuoteMeta(defaultExpectedExistsQuery)).WithArgs(login).WillReturnError(expectedError)

	exists, err := s.repository.ExistsByLogin(login)

	assert.NotNil(s.t, err)
	assert.Equal(s.t, false, exists)
	assert.Nil(s.t, s.sqlMock.ExpectationsWereMet())
}

func getUserToTest() models.User {
	user := models.User{
		ID:       1,
		Name:     "test",
		Email:    "test",
		Login:    "test",
		Password: "test",
	}
	return user
}

func getUserToSaveTest() models.User {
	return models.User{
		Name:     "test",
		Email:    "test",
		Login:    "test",
		Password: "test",
	}
}
