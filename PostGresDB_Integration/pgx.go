package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type PostGresSQLpgx struct {
	pool *pgxpool.Pool // pool of zero or more database connections
}

func NewPostgreSQLpgx() (*PostGresSQLpgx, error) {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	return &PostGresSQLpgx{
		pool: pool,
	}, nil
}

func (p *PostGresSQLpgx) Close() {
	p.pool.Close()
}

func (p *PostGresSQLpgx) FindNConst(nconst string) (Name, error) {
	query := `SELECT nconst, primary_name, birth_year, death_year FROM "names" WHERE nconst = $1`
	var res Name
	if err := p.pool.QueryRow(context.Background(), query, nconst).Scan(&res.NConst, &res.Name, &res.BirthYear, &res.DeathYear); err != nil {
		return Name{}, err
	}
	return res, nil
}
