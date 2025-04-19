package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTMaker struct {
	secretKey string
}

//This is the JWT implementation of the Maker interface ; [returns a JWTMaker Interface] ; it should have the functions defined in the Maker interface

//Symmetric key for signing and verifying JWT tokens

//Creates a new instance of JWTMaker with the provided secretKey.
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 32 { //length should not be less than 32 characters
		return nil, fmt.Errorf("invalid key size: must be at least 32 characters")
	}

	return &JWTMaker{secretKey}, nil
}

//JWTMaker implements the Maker interface by providing concrete implementations for the CreateToken and VerifyToken

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error){
 	payload, err := NewPayload(username, duration) //create new payload
	if err != nil {
		return "", err
	}


	//create jwt token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(maker.secretKey)) //sign the token with the secret key and return it

	//signing algorithm, claims (payload)

}
	
func (maker *JWTMaker)	VerifyToken(tokenString string) (*Payload, error) {
	   keyFunc := func(token *jwt.Token) (interface{}, error) {
        // Ensure the signing method is HMAC (HS256)
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
        }
        return []byte(maker.secretKey), nil
    }

    // Parse the token with claims
    token, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyFunc)
    if err != nil {
        return nil, err
    }

    // Validate the token and extract the claims
    if claims, ok := token.Claims.(*Payload); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}


//(maker *JWTMaker) : This is a pointer receiver for the JWTMaker struct.
//It allows the function to access the fields of the JWTMaker struct (e.g., secretKey) and modify them if needed.
//Using a pointer receiver avoids copying the JWTMaker struct, which is more efficient.