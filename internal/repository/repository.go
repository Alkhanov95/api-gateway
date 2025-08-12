package repository

import (
	"context"

	"github.com/alkhanov95/api-gateway/models"
	"github.com/jackc/pgx/v5" // pgx.ErrNoRows
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors" // errors.Wrap
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) (string, error) {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, name, age) VALUES ($1, $2, $3)`,
		user.ID, user.Name, user.Age,
	)
	if err != nil {
		return "", errors.Wrap(err, "create user") // возвращаем пустой ID и обёрнутую ошибку вставки
	}
	return user.ID, nil // возвращаем созданный ID, ошибки нет
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, name, age FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // не нашли пользователя: nil-указатель и nil-ошибка (это не ошибка, а отсутствие данных)
		}
		return nil, errors.Wrap(err, "get user by id") // ошибка запроса/сканирования: nil-результат и обёрнутая ошибка
	}
	return &user, nil // нашли пользователя: указатель на структуру и nil-ошибка
}

func (r *UserRepo) List(ctx context.Context) ([]models.User, error) {
	const q = `SELECT id, name, age FROM users`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, errors.Wrap(err, "list users: query failed") // не смогли выполнить SELECT: nil-список и обёрнутая ошибка
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			return nil, errors.Wrap(err, "list users: scan row") // ошибка чтения строки: nil-список и обёрнутая ошибка
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "list users: rows err") // ошибка итерации: nil-список и обёрнутая ошибка
	}
	return users, nil // успешный список пользователей и nil-ошибка
}

func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	const q = `UPDATE users SET name = $1, age = $2 WHERE id = $3`

	res, err := r.db.Exec(ctx, q, user.Name, user.Age, user.ID)
	if err != nil {
		return errors.Wrap(err, "update user") // ошибка выполнения UPDATE
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows // ничего не обновили: сигнал наверх "не найдено"
	}
	return nil // успешно обновили: ошибки нет
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM users WHERE id = $1`

	res, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return errors.Wrap(err, "delete user") // ошибка выполнения DELETE
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows // ничего не удалили: сигнал наверх "не найдено"
	}
	return nil // успешно удалили: ошибки нет
}
