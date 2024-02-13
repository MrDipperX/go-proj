package shared

const UserRegistrationTaskQueue = "USER_REGISTRATION_TASK_QUEUE"

type RegistrationBody struct{
	Name string `json:"name" binding:"required"`
	Surname string `json:"surname" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Login string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
	Confirmed bool `json:"-"`
}

type LoginBody struct{
	Login string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
	Subject  string
	Body     string
}

type RedisUserData struct {
	Username    string `json:"username"`
	ConfirmCode string `json:"confirm-code"`
	TokenExpireAt    int64  `json:"expire-at"`
	Token       string `json:"token"`
}

type UserConfirmation struct {
	Username    string `json:"username"`
	ConfirmCode string `json:"confirm-code"`
}

type ConfirmationSignal struct {
	WorkflowID string `json:"workflow_id"`
	ConfirmationCode string `json:"confirmation_code"`
}