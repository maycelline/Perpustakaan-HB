package controllers

import (
	"net/http"
	"time"

	"Perpustakaan-HB/model"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("PHBH13HarapanBangsaH1tz!!")
var tokenName = "token"

type Claims struct {
	ID                int       `json:"idUser,omitempty"`
	FullName          string    `json:"fullName,omitempty"`
	UserName          string    `json:"userName,omitempty"`
	BirthDate         time.Time `json:"birthDate,omitempty"`
	PhoneNumber       string    `json:"phone,omitempty"`
	Email             string    `json:"email,omitempty"`
	Address           string    `json:"address,omitempty"`
	AdditionalAddress string    `json:"additionalAddress,omitempty"`
	Password          string    `json:"password,omitempty"`
	UserType          string    `json:"userType,omitempty"`
	Balance           int       `json:"balance,omitempty"`
	jwt.StandardClaims
}

func generateToken(w http.ResponseWriter, member model.Member) {
	tokenExpiryTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		ID:                member.User.ID,
		FullName:          member.User.FullName,
		UserName:          member.User.UserName,
		BirthDate:         member.User.BirthDate,
		PhoneNumber:       member.User.PhoneNumber,
		Email:             member.User.Email,
		Address:           member.User.Address,
		AdditionalAddress: member.User.AdditionalAddress,
		Password:          member.User.Password,
		UserType:          member.User.UserType,
		Balance:           member.Balance,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    signedToken,
		Expires:  tokenExpiryTime,
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
	})
}

func resetUserToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
		Path:     "",
	})
}

func Authenticate(next http.HandlerFunc, accessTypeInt int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var accessType string
		if accessTypeInt == 1 {
			accessType = "MEMBER"
		} else if accessTypeInt == 2 {
			accessType = "ADMIN"
		} else {
			accessType = "OWNER"
		}
		isValidToken := validateUserToken(r, accessType)
		if !isValidToken {
			sendUnauthorizedResponse(w)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func validateUserToken(r *http.Request, accessType string) bool {
	isAccessTokenValid, userType := validateTokenFromCookies(r)

	if isAccessTokenValid {
		isUserValid := userType == accessType
		if isUserValid {
			return true
		}
	}
	return false
}

func validateTokenFromCookies(r *http.Request) (bool, string) {
	cookie, err := r.Cookie(tokenName)
	if err == nil {
		accessToken := cookie.Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return true, accessClaims.UserType
		}
	}
	return false, ""
}

func getIdFromCookies(r *http.Request) int {
	cookie, err := r.Cookie(tokenName)
	if err == nil {
		accessToken := cookie.Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return accessClaims.ID
		}
	}
	return -1
}

func getDataFromCookies(r *http.Request) (int, string, string, time.Time, string, string, string, string, string, int) {
	cookie, err := r.Cookie(tokenName)
	if err == nil {
		accessToken := cookie.Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return accessClaims.ID, accessClaims.FullName, accessClaims.UserName, accessClaims.BirthDate, accessClaims.PhoneNumber, accessClaims.Email, accessClaims.Address, accessClaims.AdditionalAddress, accessClaims.Password, accessClaims.Balance
		}
	}
	return -1, "", "", time.Time{}, "", "", "", "", "", -1
}
