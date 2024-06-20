package services

import (
	"ShortUrl/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DB *sqlx.DB
}

func (s *AuthService) Login(username string, password string) (*models.User, error) {

	var user models.User
	//Sử dụng phương thức Get của sqlx.DB để thực hiện truy vấn cơ sở dữ liệu.
	err := s.DB.Get(&user, "SELECT * FROM users WHERE username=$1", username)
	//Nếu truy vấn không thành công trả về nil và lỗi err.
	if err != nil {
		return nil, err
	}

	//So sánh mật khẩu đã mã hóa lưu trong cơ sở dữ liệu với mật khẩu người dùng nhập vào.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (s *AuthService) Register(username, password string) (*models.User, error) {
	//mã hóa mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       uuid.New(),
		Username: username,
		Password: string(hashedPassword),
	}

	_, err = s.DB.NamedExec(`INSERT INTO users (id, username, password) VALUES (:id, :username, :password)`, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
