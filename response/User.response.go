package response

import (
	"github.com/youtube/domain"
)

type GetUserInfoResponse struct {
	ID              string   `json:"_id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	Image           string   `json:"image"`
	Subscribers     uint     `json:"subscribers"`
	SubscribedUsers []string `json:"subscribedUsers"`
}

func (u *GetUserInfoResponse) GetDetails(user *domain.User) *GetUserInfoResponse {
	u.ID = user.ID.Hex()
	u.Name = user.Name
	u.Email = user.Email
	u.Image = user.Image
	u.Subscribers = user.Subscribers
	u.SubscribedUsers = user.SubscribedUsers
	return u
}

func (u *GetUserInfoResponse) GetAllUsersResponse(users []*domain.User) []*GetUserInfoResponse {
	var allUsers []*GetUserInfoResponse
	for _, user := range users {
		allUsers = append(allUsers, u.GetDetails(user))
	}
	return allUsers
}
