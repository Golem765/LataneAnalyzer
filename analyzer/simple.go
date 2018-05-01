package analyzer

import (
	"univ/course/model"
	"math"
)

func FindPressure(users []model.User, beta, alpha float64) []model.User {
	for i := 0; i < len(users); i = i + 1 {
		var pressure float64
		var mention = 0
		user := users[i]
		for j, u := range users {
			if i == j {
				continue
			}
			mention = mention + user.Mentions[u.ID]
		}

		nMentions := 0.0
		for j, uj := range users {
			for k, uk := range users {
				if k != 0 && j > k {
					ujMentions := user.Mentions[uj.ID]
					ukMentions := user.Mentions[uk.ID]
					if ujMentions*ukMentions != 0 {
						nMentions = nMentions +
							(float64(ujMentions-ukMentions) / math.Pow(distance(uj, uk), alpha))

					}
				}
			}
		}

		pressure = -beta*float64(mention) - float64(nMentions)
		user.Pressure = pressure
		users[i] = user
	}
	return users
}

func distance(user1, user2 model.User) float64 {
	for _, u := range user1.Friends {
		if u == user2.ID {
			return 1.0
		}
	}
	return 2.0
}

func FindInfluence(users []model.User) []model.User {
	influence := make(map[model.ObjectId]int)
	for i := 0; i < len(users); i = i + 1 {
		user := users[i]
		for userId := range user.Mentions {
			influence[userId] = influence[userId] + 1
		}
	}
	for i := 0; i < len(users); i = i + 1 {
		user := users[i]
		user.Influence = float64(influence[user.ID])
		users[i] = user
	}
	return users
}
