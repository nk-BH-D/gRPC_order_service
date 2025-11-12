package order_db

import (
	"context"
	"database/sql"
	//"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(dsn string) (*Postgres, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	p := &Postgres{DB: db}
	if err := p.ensureSchema(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Postgres) Close() error {
	return p.DB.Close()
}

func (p *Postgres) ensureSchema() error {
	// Создаём простую таблицу orders если её нет
	schema := `
    CREATE TABLE IF NOT EXISTS orders (
        id TEXT PRIMARY KEY,
        item TEXT NOT NULL,
        quantity INTEGER NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
    );`
	_, err := p.DB.Exec(schema)
	return err
}

// InsertOrder возвращает ошибку или nil
func (p *Postgres) InsertOrder(ctx context.Context, id, item string, quantity int32) error {
	_, err := p.DB.ExecContext(ctx, "INSERT INTO orders (id, item, quantity) VALUES ($1,$2,$3)", id, item, quantity)
	return err
}

func (p *Postgres) GetOrder(ctx context.Context, id string) (string, int32, error) {
	var item string
	var quantity int32
	row := p.DB.QueryRowContext(ctx, "SELECT item, quantity FROM orders WHERE id=$1", id)
	if err := row.Scan(&item, &quantity); err != nil {
		return "", 0, err
	}
	return item, quantity, nil
}

func (p *Postgres) UpdateOrder(ctx context.Context, id, item string, quantity int32) error {
	_, err := p.DB.ExecContext(ctx, "UPDATE orders SET item=$1, quantity=$2 WHERE id=$3", item, quantity, id)
	return err
}

func (p *Postgres) DeleteOrder(ctx context.Context, id string) error {
	_, err := p.DB.ExecContext(ctx, "DELETE FROM orders WHERE id=$1", id)
	return err
}

func (p *Postgres) ListOrders(ctx context.Context) ([]map[string]any, error) {
	rows, err := p.DB.QueryContext(ctx, "SELECT id, item, quantity FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []map[string]any
	for rows.Next() {
		var id, item string
		var q int32
		if err := rows.Scan(&id, &item, &q); err != nil {
			return nil, err
		}
		list = append(list, map[string]any{"id": id, "item": item, "quantity": q})
	}
	return list, nil
}
