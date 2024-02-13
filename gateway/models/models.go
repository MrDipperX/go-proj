package models

type RegistrationBody struct{
	Name string `json:"name" binding:"required"`
	Surname string `json:"surname" binding:"required"`
	Phone string `json:"phone" binding:"required,numeric"`
	Email string `json:"email" binding:"required,email"`
	Login string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginBody struct{
	Login string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}