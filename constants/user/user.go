package constants

type UserFromJSON struct {
	Email    string `json: "email"`
	Password string `json: "password"`
}

type UserToDB struct {
	Email        string
	PasswordHash string
}
