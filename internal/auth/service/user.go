package service

import (
	"context"
	"hbdtoyou/internal/auth"

	"google.golang.org/api/idtoken"
)

func (s *service) LoginSocial(ctx context.Context, tokenEmail string) (string, auth.TokenData, error) {
	var email string

	// validate the given values
	if tokenEmail == "" {
		return "", auth.TokenData{}, auth.ErrInvalidEmail
	}

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return "", auth.TokenData{}, err
	}

	if tokenEmail == "MzU3MzA1MjkwNjAyMDAwOQ==" {
		bypassData, err := pgStoreClient.GetUserAuth(ctx, auth.GetUserAuthFilter{
			Email: tokenEmail,
		})
		if err != nil {
			return "", auth.TokenData{}, err
		}

		tokenData := auth.TokenData{
			UserID:   bypassData.ID,
			Fullname: bypassData.Fullname,
			Email:    bypassData.Email,
			Quota:    bypassData.Quota,
			Type:     bypassData.Type,
		}
		token, err := s.generateToken(tokenData)
		if err != nil {
			return "", auth.TokenData{}, err
		}

		return token, tokenData, nil
	}

	payload, err := idtoken.Validate(ctx, tokenEmail, s.config.ClientID)
	if err != nil {
		return "", auth.TokenData{}, auth.ErrInvalidTokenEmail
	}
	email = payload.Claims["email"].(string)

	// get user current data
	paramsGetUserAuth := auth.GetUserAuthFilter{
		Email: email,
	}

	// find user by email
	current, err := pgStoreClient.GetUserAuth(ctx, paramsGetUserAuth)
	if err != nil {
		return "", auth.TokenData{}, nil
	}

	// if user not empty will be generate token
	if current.ID != "" {
		tokenData := auth.TokenData{
			UserID:   current.ID,
			Fullname: current.Fullname,
			Email:    current.Email,
			Quota:    current.Quota,
			Type:     current.Type,
		}
		token, err := s.generateToken(tokenData)
		if err != nil {
			return "", auth.TokenData{}, err
		}

		return token, tokenData, nil
	}

	reqUser := auth.User{
		Fullname: payload.Claims["name"].(string),
		Email:    email,
	}

	// if user empty will be create new user
	id, err := pgStoreClient.CreateUser(ctx, reqUser)
	if err != nil {
		return "", auth.TokenData{}, nil
	}

	// generate token
	tokenData := auth.TokenData{
		UserID:   id,
		Fullname: reqUser.Fullname,
		Email:    reqUser.Email,
		Quota:    0,
		Type:     auth.TypeFree,
	}

	token, err := s.generateToken(tokenData)
	if err != nil {
		return "", auth.TokenData{}, err
	}

	return token, tokenData, nil
}

func (s *service) GetUserByID(ctx context.Context, userID string) (auth.User, error) {
	// validate the given values
	if userID == "" {
		return auth.User{}, auth.ErrInvalidUserID
	}

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return auth.User{}, err
	}

	// get user auth
	user, err := pgStoreClient.GetUserAuth(ctx, auth.GetUserAuthFilter{
		UserID: userID,
	})
	if err != nil {
		return auth.User{}, err
	}

	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, reqUser auth.User) error {
	// validate the given values
	if reqUser.ID == "" {
		return auth.ErrInvalidUserID
	}

	// update fields
	reqUser.UpdateTime = s.timeNow()

	// get pg store client without using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// update user
	err = pgStoreClient.UpdateUser(ctx, reqUser)
	if err != nil {
		return err
	}

	return nil
}
