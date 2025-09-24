package main

import (
	"context"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostGresSQLsqlx struct {
	db *sqlx.DB // pool of zero or more database connections
}

func NewPostgreSQLsqlx() (*PostGresSQLsqlx, error) {
	// ConnectContext to a database and verify with a ping
	db, err := sqlx.ConnectContext(context.Background(), "postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	// if err := pool.Ping(); err != nil {
	// 	return nil, err
	// }
	return &PostGresSQLsqlx{
		db: db,
	}, nil
}

func (p *PostGresSQLsqlx) Close() {
	p.db.Close()
}

func (p *PostGresSQLsqlx) FindNConst(nconst string) (Name, error) {
	query := `SELECT nconst, primary_name, birth_year, death_year FROM "names" WHERE nconst = $1`
	var result struct {
		NConst    string `db:"nconst"`
		Name      string `db:"primary_name"`
		BirthYear string `db:"birth_year"`
		DeathYear string `db:"death_year"`
	}
	//StructScan a single row into destination
	if err := p.db.QueryRowx(query, nconst).StructScan(&result); err != nil {
		return Name{}, err
	}

	return Name{
		NConst:    result.NConst,
		Name:      result.Name,
		BirthYear: result.BirthYear,
		DeathYear: result.DeathYear,
	}, nil
}
