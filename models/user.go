package models

import "time"

type User struct {
	BaseModel
	Email           string     `gorm:"unique;not null" json:"email"`
	Password        string     `gorm:"not null" json:"-"`
	ActivatedAt     *time.Time `gorm:"index" json:"activated_at"`
	ActivationToken string     `gorm:"index" json:"activation_token"`
	RefreshToken    string     `gorm:"index" json:"refresh_token"`
}

func (user *User) Create() error {
	err := DB.Create(&user).Error
	return err
}

func (user *User) Find() error {
	err := DB.Where("email = ?", user.Email).First(&user).Error
	return err
}

func (user *User) FindByActivationToken() error {
	err := DB.Where("activation_token = ?", user.ActivationToken).First(&user).Error
	return err
}

func (user *User) Update() error {
	err := DB.Save(&user).Error
	return err
}
