package models

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	JwtSec = "InR5cCIljaldskWRtaW4iLCJ1c2Vy"
)

// TODO: вынести в конфиг
// const tokenExpiredTime = 1440
// const refreshTokenExpiredMinutes = 201600

// Sessions godoc
type Sessions struct {
	tableName           struct{}  `sql:"crm_sessions"`
	ID                  int       `json:"id" sql:",pk"`
	UserID              int       `json:"user_id"`
	RefreshToken        string    `json:"refresh_token" description:"Токен для обновления"`
	SessionEnd          time.Time `json:"session_end" description:"Дата когда токен отозван"`
	RefreshTokenUsed    time.Time `json:"refresh_token_used" description:"Дата использования токена"`
	RefreshTokenExpired time.Time `json:"refrash_expired" description:"Дата протухания токена"`
	CreatedAt           time.Time `sql:"default:now()" json:"created_at" description:"Дата создания"`
}

type (
	// LoginRequest requested data when logging in
	LoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	// LoginRefreshRequest godoc
	LoginRefreshRequest struct {
		RefreshToken string `json:"refresh"`
	}

	// TokenClaim JWT token structure
	TokenClaim struct {
		Role     string `json:"role"`
		UserID   int    `json:"user_id"`
		Login    string `json:"login"`
		UserUUID string `json:"user_uuid"`
		jwt.StandardClaims
	}

	// LoginResponse responsed when requesting token
	LoginResponse struct {
		UserUUID string `json:"user_uuid"`
		Token    string `json:"token"`
	}
)

// GenerateJWT generates new token
func (logResp *LoginResponse) GenerateJWT(user UsersCRM) error {

	mySigningKey := []byte(JwtSec)

	claims := TokenClaim{
		UserID:   user.ID,
		Role:     user.Role,
		UserUUID: user.UUID,
		Login:    user.Login,
	}
	claims.IssuedAt = time.Now().Unix()

	dur := time.Minute * 3600
	claims.ExpiresAt = time.Now().Add(dur).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return err
	}
	logResp.Token = ss

	return nil
}

// AuthenticateUser godoc
func AuthenticateUser(data LoginRequest) (LoginResponse, error) {

	var oper UsersCRM
	var login LoginResponse
	err := db.Model(&oper).
		Where("deleted is not true AND Login = ?", data.Login).
		First()
	if err != nil {
		return login, err
	}

	if data.Password != oper.Password {
		return login, fmt.Errorf("неверный пароль")
	}

	err = login.GenerateJWT(oper)
	if err != nil {
		return login, err
	}
	login.UserUUID = oper.UUID
	return login, nil
}

func expireUserTokens(userID int) error {

	var sessOld Sessions

	_, err := db.Model(&sessOld).
		Set("refresh_token_used = ?", time.Now()).
		Where("user_id = ?", userID).
		Update()

	if err != nil {
		return err
	}
	return nil
}

func expireToken(token string) (UsersCRM, error) {

	var oper UsersCRM
	var sessOld Sessions

	_, err := db.Model(&sessOld).
		Set("refresh_token_used = ?", time.Now()).
		Where("refresh_token = ? AND CURRENT_TIMESTAMP < refresh_token_expired AND refresh_token_used is NULL", token).
		Returning("*").
		Update(&sessOld)

	if err != nil {
		return oper, fmt.Errorf("Refresh token not found, %s", err)
	}

	err = db.Model(&oper).
		Where("deleted is not true AND ID = ?", sessOld.UserID).
		First()
	if err != nil {
		return oper, err
	}
	return oper, nil
}
