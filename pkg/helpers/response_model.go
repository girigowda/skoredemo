package helpers

type Card_success_response struct {
	Data Success_response_data
}

type Success_response_data struct {
	IdNumber           string
	Name               string
	BloodType          string
	Religion           string
	Gender             string
	BirthPlaceBirthday string
	Province           string
	City               string
	District           string
	Village            string
	Rtrw               string
	Occupation         string
	ExpiryDate         string
	Nationality        string
	MaritalStatus      string
	Address            string
}

type Success_response struct {
	Data Success_response_code
}
type Success_response_code struct {
	Code    int
	Message string
	Score   string
}
type Error_response struct {
	Errors Error_response_code
}
type Error_response_code struct {
	Code   string
	Title  string
	Detail string
}
