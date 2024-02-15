package shared


const UserAuthTaskQueue = "USER_AUTH_TASK_QUEUE"

const SecretKey = "AppleVisionPro"

type UserData struct {
	Username    string `json:"username"`
}

type UserToken struct {
	Token       string `json:"token"`
	Expire      string  `json:"expire"`
}

