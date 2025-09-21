package models

import "time"

// User - сущность пользователя.
// Тег `db` нужен для sqlx, а тег `json` для сериализации

type User struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
