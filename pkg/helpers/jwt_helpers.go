package helpers

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Phone_number string
	jwt.StandardClaims
}

type TOKEN_RESP struct {
	Status string
	Phone  string
}

var jwtKey = []byte("my_secret_key")

func GenerateJWT(phone string) string {
	//get body data
	// expirationTime := time.Now().Add(10 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Phone_number:   phone,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			// ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// println("token is",token)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	tokenString = "Bearer " + tokenString
	println("tokenString", tokenString)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		println("error in jwt creation")
	}
	// Finally, we set the client cookie for "token" as the JWT we just generated
	return tokenString
}

func ValidateJWT(token string) string {

	// c := ctx.Request.Header.Get("Authorization")
	split := strings.Split(token, " ")
	println(split[0])
	length := len(reflect.ValueOf(split).Interface().([]string))
	println("length", length)
	resp := TOKEN_RESP{}
	if length == 2 && split[0] == "Bearer" {
		c := split[1]
		// Get the JWT string
		tknStr := c
		// Initialize a new instance of `Claims`
		claims := &Claims{}
		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		println(claims.Phone_number)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				resp.Status = "invalid"
			}
			resp.Status = "invalid"
		}
		if !tkn.Valid {
			resp.Status = "invalid"
		}
		resp.Status = "success"
		resp.Phone = claims.Phone_number
		// Finally, return the welcome message to the user, along with their
		// username given in the token
	} else {
		resp.Status = "invalid"
	}
	converted_resp, _ := json.Marshal(resp)
	return string(converted_resp)
}
