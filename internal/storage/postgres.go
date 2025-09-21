package storage

import (
	"context"
	"fmt"
	"simple-golang-crud/pkg/models"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // для регистрации драйвера postgres
)

// Postgres - Структура-обертка для работы с БД
type Postgres struct {
	DB *sqlx.DB
}

// NewPostgres - открытие подключения к БД
func NewPostgres(dsn string) (*Postgres, error) { // Возвращаем кортеж из структуры-обертки и ошибки
	db, err := sqlx.Connect("postgres", dsn) // Подключение и проверка
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Пул
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Postgres{DB: db}, nil // Возвращаем указатель на структуру-обертку и кладем сконфигурированное поле DB и отстутсиве ошибки (nil)
}

// Close - закрытие пула соединений
func (p *Postgres) Close() error {
	return p.DB.Close()
}

// CreateUser - создание пользователя
func (p *Postgres) CreateUser(ctx context.Context, u *models.User) error { // Контекст нужен для того чтобы в случае зависания запроса отменить его и вернуть ошибку
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at`

	return p.DB.QueryRowxContext(ctx, query, u.Name, u.Email).Scan(&u.ID, &u.CreatedAt) // Выполняем запрос и маппим
}

// GetUser - получение пользователя по id
func (p *Postgres) GetUser(ctx context.Context, id int64) (*models.User, error) {
	var u models.User // Пустая структура

	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`

	if err := p.DB.GetContext(ctx, &u, query, id); err != nil { // Выполняем запрос и маппим данные в структуру
		return nil, err
	}

	return &u, nil
}

// GetAllUsers - получение всех пользователей
func (p *Postgres) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	query := `SELECT id, name, email, created_at FROM users ORDER BY id`

	if err := p.DB.SelectContext(ctx, &users, query); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser — обновить имя и email
func (p *Postgres) UpdateUser(ctx context.Context, u *models.User) error {

	query := `UPDATE users SET name=$1, email=$2 WHERE id=$3 RETURNING created_at` // Возвращаем created_at (в БД не меняется)

	return p.DB.QueryRowxContext(ctx, query, u.Name, u.Email, u.ID).Scan(&u.CreatedAt)
}

// DeleteUser — удалить пользователя по id
func (p *Postgres) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id=$1`

	_, err := p.DB.ExecContext(ctx, query, id)

	return err
}
