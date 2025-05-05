package model

type User struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Credential string `json:"credential"`
	IsLogin    string `json:"isLogin"`
	Cif        string `json:"cif"`
}
