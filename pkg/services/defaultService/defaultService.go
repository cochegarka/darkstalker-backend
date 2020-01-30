package defaultService

import (
	"context"
	"darkstalker/pkg/services"
	"github.com/SevereCloud/vksdk/api"
	"github.com/pkg/errors"
)

func NewDefaultService(token string) services.Service {
	return &defaultService{api.Init(token)}
}

type defaultService struct {
	vk *api.VK
}

func (d *defaultService) StalkUser(_ context.Context, userId string) (map[string]interface{}, error) {
	dossier := make(map[string]interface{})

	// Get info about user
	users, err := d.vk.UsersGet(api.Params{
		"user_ids": userId,
		"fields":   userFields,
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot get users from api")
	}

	if len(users) == 0 {
		return nil, errors.New("there is no user with given id: " + userId)
	}
	user := users[0]

	dossier["user"] = user

	// Get info about their friends
	rawFriends, err := friendsGet(d.vk, api.Params{
		"user_id": user.ID,
		"fields":  friendsFields,
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot get user's friends from api")
	}

	// Get info about user groups
	friendsLists, _ := d.vk.FriendsGetLists(api.Params{
		"user_id": user.ID,
	})
	//if err != nil {
	//	return nil, errors.Wrap(err, "cannot get user's friends lists from api")
	//}

	// Map VK API friends to slice of our friend struct
	friends := make([]friend, rawFriends.Count)

	for i, v := range rawFriends.Items {
		// Retrieve real group name
		// Speed it up with map[int]string?
		group := ""
		for _, vv := range friendsLists.Items {
			if vv.ID == v.ListId {
				group = vv.Name
				break
			}
		}

		friends[i] = friend{
			Id:      v.Id,
			Name:    v.FirstName + " " + v.LastName,
			Photo:   v.Photo,
			Group:   group,
			GroupID: v.ListId,
		}
	}

	dossier["friends"] = friends

	return dossier, nil
}

type friend struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Photo   string `json:"photo"`
	Group   string `json:"group"`
	GroupID int    `json:"group_id"`
}

type friendsGetResponse struct {
	Count int `json:"count"`
	Items []struct {
		Id        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Photo     string `json:"photo_100"`
		ListId    int    `json:"list_id"`
	} `json:"items"`
}

func friendsGet(vk *api.VK, params api.Params) (response friendsGetResponse, err error) {
	err = vk.RequestUnmarshal("friends.get", params, &response)
	return
}
