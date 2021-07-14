package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	test "user_test/gen/go/proto"
)

const (
	userPageKey = "user:first_page"
	userTTL     = 3 * time.Minute
)

type UserModel struct {
	ID          int64
	FIO         string
	Email       string
	PhoneNumber string
	CreatedAt   time.Time
}

type Repository struct {
	db    *sql.DB
	cache redis.UniversalClient
}

func NewRepository(db *sql.DB, cache redis.UniversalClient) *Repository {
	return &Repository{db: db, cache: cache}
}

// InsertUser insert record
func (r *Repository) InsertUser(ctx context.Context, user *test.InsertUserRequest) error {
	query := "INSERT INTO public.user_name (fio, email, phone_number, created_at) VALUES ($1, $2, $3, now())"
	err := r.db.QueryRowContext(ctx, query, user.Fio, user.Email, user.Phone).Err()
	if err != nil {
		return err
	}
	_ = r.cache.Del(ctx, userPageKey).Err()
	return nil
}

// DeleteUser Create saves a new user record
func (r *Repository) DeleteUser(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM public.user_name WHERE id = $1", id)
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	_ = r.cache.Del(ctx, userPageKey).Err()
	return nil
}

// ListUsers List retrieves the user records.
func (r *Repository) ListUsers(ctx context.Context, offset, limit int) ([]UserModel, error) {
	var ads []UserModel
	if offset == 0 {
		value, err := r.cache.Get(ctx, userPageKey).Bytes()
		if err != redis.Nil {
			err = json.Unmarshal(value, &ads)
			if err == nil {
				return ads, nil
			}
		}
	}

	rows, err := r.db.QueryContext(ctx, "SELECT id,  fio, email, phone_number, created_at FROM public.user_name LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u UserModel
		err = rows.Scan(&u.ID, &u.FIO, &u.Email, &u.PhoneNumber, &u.CreatedAt)
		if err != nil {
			log.Printf("Unable to scan the row. %v", err)
			return nil, err
		}

		ads = append(ads, u)
	}
	if offset == 0 {
		b, _ := json.Marshal(&ads)
		r.cache.Set(ctx, userPageKey, b, userTTL)
	}

	return ads, nil
}
