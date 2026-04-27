package dto

import "media-equipment-tracker/internal/domain"

type ShortUserItem struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Nice     int    `json:"nice"`
	IsAdmin  bool   `json:"is_admin"`
}

func ShortFromUser(user *domain.User) ShortUserItem {
	return ShortUserItem{
		ID:       user.ID.String(),
		FullName: user.FullName,
		Email:    user.Email,
		Nice:     user.Nice,
		IsAdmin:  user.IsAdmin,
	}
}
