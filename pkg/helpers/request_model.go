package helpers

type LoginRequestBody struct {
	Phone      string
	Password   string
	Ip_address string
	Location   string
	Country    string
	Udid       string
}

type RegisterRequestBody struct {
	Phone      string
	Ip_address string
	Location   string
	Country    string
	Udid       string
	Privacy    bool
}

type OTPRequestBody struct {
	Otp string
}

type OTPResendBody struct {
	Phone string
}

type SavepinRequestBody struct {
	Phone    string
	Password string
}
