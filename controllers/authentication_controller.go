package controllers

import (
	"net/http"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("PHBH13HarapanBangsaH1tz!!")
var tokenName = "token"

type Claims struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserType int    `json:"user_type"`
	jwt.StandardClaims
}

func generateToken(w http.ResponseWriter, id int, name string, userType int) {
	// tokenExpiryTime := time.Now().Add(5 * time.Minute)

	// claims := &Claims{
	// 	ID:       id,
	// 	Name:     name,
	// 	UserType: userType,
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: tokenExpiryTime.Unix(),
	// 	},
	// }

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// signedToken, err := token.SignedString(jwtKey)
	// if err != nil {
	// 	return
	// }

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     tokenName,
	// 	Value:    signedToken,
	// 	Expires:  tokenExpiryTime,
	// 	Secure:   false,
	// 	HttpOnly: true,
	// 	Path:     "/",
	// })
}

func resetUserToken(w http.ResponseWriter) {
	// http.SetCookie(w, &http.Cookie{
	// 	Name:     tokenName,
	// 	Value:    "",
	// 	Expires:  time.Now(),
	// 	Secure:   false,
	// 	HttpOnly: true,
	// 	Path:     "",
	// })
}

func Authenticate(next http.HandlerFunc, accessType int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// isValidToken := validateUserToken(r, accessType)
		// if !isValidToken {
		// 	sendUnAuthorizedResponse(w)
		// } else {
		// 	next.ServeHTTP(w, r)
		// }
	})
}

func validateUserToken(r *http.Request, accessType int) /*bool*/ {
	// isAccessTokenValid, id, email, userType := validateTokenFromCookies(r)
	// fmt.Print(id, email, userType, accessType, isAccessTokenValid)

	// if isAccessTokenValid {
	// 	isUserValid := userType == accessType
	// 	if isUserValid {
	// 		return true
	// 	}
	// }
	// return false
}

func validateTokenFromCookies(r *http.Request) /*(bool, int, string, int)*/ {
	// cookie, err := r.Cookie(tokenName)
	// if err == nil {
	// 	accessToken := cookie.Value
	// 	accessClaims := &Claims{}
	// 	parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
	// 		return jwtKey, nil
	// 	})
	// 	fmt.Println(err)
	// 	fmt.Println(parsedToken)
	// 	if err == nil && parsedToken.Valid {
	// 		return true, accessClaims.ID, accessClaims.Name, accessClaims.UserType
	// 	}
	// }
	// return false, -1, "", -1
}
