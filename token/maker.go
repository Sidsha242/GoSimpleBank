package token

import "time"

// general token maker interface [abstracts the logic for creating and verifying tokens]
// This interface allows for different implementations of token makers (e.g., JWT, OAuth, Pasteo) to be used interchangeably in the application.
//***Interfaces are particularly useful when you need to define a contract for behavior while allowing multiple implementations.***

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	
	VerifyToken(token string) (*Payload, error)
	
}