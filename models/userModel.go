package models

type UserModel struct {
	ID        uint     `json:"id" gorm:"primaryKey"`
	Name      string   `json:"name"`
	Email     string   `json:"email" binding:"required"`
	Password  string   `json:"password" binding:"required,min=5"`
	ImageURL  string   `json:"image_url"`
	Bio       string   `json:"bio"`
}

type CreateUser struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required,min=5"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUser struct { 
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	Bio      string `json:"bio"`
}