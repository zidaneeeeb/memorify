package service

import (
	"context"
	"hbdtoyou/internal/auth"

	"github.com/golang-jwt/jwt/v4"
)

// jwtClaimss is the claims encapsulated in JWT-generated token.
type jwtClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Type     int    `json:"type"`
	jwt.RegisteredClaims
}

// parseTokenData parse token data from jwt claims.
func (jwtc jwtClaims) parseTokenData() auth.TokenData {
	return auth.TokenData{
		UserID:   jwtc.UserID,
		Username: jwtc.Username,
		Fullname: jwtc.Fullname,
		Email:    jwtc.Email,
		Type:     auth.Type(jwtc.Type),
	}
}

// formatTokenData format token data into jwt claims.
func formatTokenData(data auth.TokenData) jwtClaims {
	return jwtClaims{
		UserID:   data.UserID,
		Username: data.Username,
		Fullname: data.Fullname,
		Email:    data.Email,
		Type:     data.Type.Value(),
	}
}

func (s *service) ValidateToken(ctx context.Context, token string) (auth.TokenData, error) {
	if token == "" {
		return auth.TokenData{}, auth.ErrInvalidToken
	}

	// get jwt token object
	jwtToken, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.config.TokenSecretKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return auth.TokenData{}, auth.ErrInvalidToken
		}

		return auth.TokenData{}, err
	}

	// parse jwt claims
	claims, ok := jwtToken.Claims.(*jwtClaims)
	if !ok {
		return auth.TokenData{}, auth.ErrInvalidToken
	}

	return claims.parseTokenData(), nil
}

func (s *service) RefreshToken(ctx context.Context, token string) (string, error) {
	// validate token
	data, err := s.ValidateToken(ctx, token)
	if err != nil {
		return "", err
	}

	// generate a new token
	return s.generateToken(data)
}

// generateToken returns a new token that encapsulates the
// given token data with some additional information:
//   - token expiration time
//
// Token is generated using JWT HS256.
func (s *service) generateToken(data auth.TokenData) (string, error) {
	claims := formatTokenData(data)

	// add expirations time
	expiresAt := s.timeNow().Add(s.config.TokenExpiration)
	claims.ExpiresAt = jwt.NewNumericDate(expiresAt)

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign token with secret key
	signedToken, err := token.SignedString([]byte(s.config.TokenSecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
