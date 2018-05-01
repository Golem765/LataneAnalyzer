package mentions

import (
	"univ/course/model"
	"univ/course/service/user"
)

func ProcessUserMentions(users userservice.Array, userService userservice.Service) []model.User {
	var mentionedUsers userservice.Array = make([]model.User, 0)
	for i := 0; i < len(users); i = i + 1 {
		u := users[i]
		if u.Checked == model.CHECKED {
			mentions := make(map[model.ObjectId]int)
			if u.Posts != nil {
				for _, p := range u.Posts {
					if p.Type == model.REPOST {
						mentionedCount := mentions[p.AuthorId]
						mentionedCount = mentionedCount + 1
						mentions[p.AuthorId] = mentionedCount

						mentioned, err := userService.GetOne(p.AuthorId)
						if err == nil {
							mentionedUsers.Add(mentioned)
						}
					}
					if p.Mentions != nil {
						for _, m := range p.Mentions {
							if m.Type == model.FB {
								mentionedCount := mentions[m.UserId]
								mentionedCount = mentionedCount + 1
								mentions[m.UserId] = mentionedCount

								mentioned, err := userService.GetOne(m.UserId)
								if err == nil {
									mentionedUsers.Add(mentioned)
								}
							}
						}
					}
				}
			}
			u.Mentions = mentions
			users[i] = u
		}
	}

	usersPointer := &users
	for _, u := range mentionedUsers {
		usersPointer.Add(u)
	}

	return *usersPointer
}
