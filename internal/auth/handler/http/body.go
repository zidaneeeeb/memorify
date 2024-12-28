package http

import "hbdtoyou/internal/auth"

// loginSocialRequestData is the data from user to perform loginSocial.
type loginSocialRequestData struct {
	TokenEmail string `json:"token_email"`
}

type loginResponseData struct {
	UserID   string `json:"user_id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type userHTTP struct {
	ID       *string `json:"id"`
	Fullname *string `json:"fullname"`
	Email    *string `json:"email"`
	Type     *string `json:"type"`
	Quota    *int    `json:"quota"`
}

func formatUser(u auth.User) userHTTP {
	types := u.Type.String()

	return userHTTP{
		ID:       &u.ID,
		Fullname: &u.Fullname,
		Email:    &u.Email,
		Type:     &types,
		Quota:    &u.Quota,
	}
}

func (u userHTTP) parseUser(out *auth.User) error {
	if u.ID != nil {
		out.ID = *u.ID
	}

	if u.Fullname != nil {
		out.Fullname = *u.Fullname
	}

	if u.Email != nil {
		out.Email = *u.Email
	}

	if u.Quota != nil {
		out.Quota = *u.Quota
	}

	return nil
}
