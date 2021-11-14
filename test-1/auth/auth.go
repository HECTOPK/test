package auth

const token = "acshginmbovxadsvnaf"

func GetToken() string {
	return token
}

func CheckToken(s string) bool {
	return s == token
}
