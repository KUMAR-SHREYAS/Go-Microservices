package main

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type PostGresSQLsql struct {
	pool *sql.DB // pool of zero or more database connections
}

func NewPostgreSQLsql() (*PostGresSQLsql, error) {
	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(); err != nil {
		return nil, err
	}
	return &PostGresSQLsql{
		pool: pool,
	}, nil
}

func (p *PostGresSQLsql) Close() {
	p.pool.Close()
}

func (p *PostGresSQLsql) FindNConst(nconst string) (Name, error) {
	query := `SELECT nconst, primary_name, birth_year, death_year FROM "names" WHERE nconst = $1`
	var res Name
	if err := p.pool.QueryRowContext(context.Background(), query, nconst).Scan(&res.NConst, &res.Name, &res.BirthYear, &res.DeathYear); err != nil {
		return Name{}, err
	}
	return res, nil
}
