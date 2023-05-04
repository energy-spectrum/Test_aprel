package migration

import (
	"app/db"
	"app/internal/util"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

func FillUsers(store db.Store) {
	for i := 0; i < 5; i++ {
		login := fmt.Sprintf("l%d", i)
		hashedPassword := util.HashPassword(fmt.Sprintf("p%d", i))
		logrus.Printf("login %s, hashedPassword %s", login, hashedPassword)
		store.UserRepo.Create(context.Background(), login, hashedPassword)
	}
}
