package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Login    string `json:"login" gorm:"unique"`
	Password string `json:"password"`
}

type List struct {
	ID    uint64
	Title string
	Owner *uint64
}

type Item struct {
	ID          uint64  `json:"id"`
	UserID      *uint64 `json:"user_id"`
	ListID      uint64  `json:"list_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

func NewListFromDTO(listDTO *ListDTO) *List {
	return &List{Title: listDTO.ListParam.Title}
}

func NewItemFromDTO(itemDTO *ItemDTO) *Item {
	return &Item{
		Title:       itemDTO.Title,
		Description: itemDTO.Description,
		UserID:      itemDTO.UserID,
	}
}

func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}

	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
