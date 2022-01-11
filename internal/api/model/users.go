package model

type AuthResponse struct {
	Success  bool `json:"success"`
	AuthUser struct {
		FullName  string `json:"fullName"`
		UserName  string `json:"userName"`
		UserEmail string `json:"userEmail"`
		UserId    string `json:"userId"`
	}
}
