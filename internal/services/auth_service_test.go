package services

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) NamedExec(query string, arg interface{}) (sql.Result, error) {
	args := m.Called(query, arg)
	return nil, args.Error(1)
}

func (m *MockDB) Get(dest interface{}, query string, args ...interface{}) error {
	return m.Called(dest, query, args).Error(0)
}

/*
func Test_Login(t *testing.T) {
	mockDB := new(MockDB)
	authService := &AuthService{DB: mockDB}

	username := "testuser"
	password := "testpass"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{
		ID:       uuid.New(),
		Username: username,
		Password: string(hashedPassword),
	}

	mockDB.On("Get", mock.Anything, "SELECT * FROM users WHERE username=$1", username).Run(func(args mock.Arguments) {
		dest := args.Get(0).(*models.User)
		*dest = *user
	}).Return(nil)

	loggedInUser, err := authService.Login(username, password)

	assert.NoError(t, err)
	assert.Equal(t, user.Username, loggedInUser.Username)
	assert.Equal(t, user.ID, loggedInUser.ID)
	mockDB.AssertExpectations(t)
}*/
