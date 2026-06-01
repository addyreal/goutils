// Example config:
//
//	servers:
//	  auth:
//	    host: "example.com"
//	    port: 1234
//	    dbname: "mydb"
//	    user: "myuser"
//	    password: "mypassword"
//	    connect_timeout: 10
//	  main:
//	    hostaddr: "2.3.4.5"
//	    port: 2345
//	    dbname: "somedb"
//	    user: "someuser"
//	    password: "somepassword"
//	    sslmode: "verify-full"
//	    sslrootcert: "/etc/ssl/certs/ca-certificates.crt"
//	    ssl_min_protocol_version: "TLSv1.2"
//
//	queries:
//	  get-basic-auth:
//	    server: "auth"
//	    query: "SELECT * FROM public.users;"
//	  get-thing:
//	    server: "main"
//	    query: "SELECT thing FROM public.stuff;"
package psql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var (
	ErrNotFound error = errors.New("not found")
)

type Database struct {
	servers map[string]*pgxpool.Pool
	queries map[string]query
}

type query struct {
	query string
	conn  *pgxpool.Pool
}

func FromYaml(path string) (*Database, error) {
	cfg, err := configFromYaml(path)
	if err != nil {
		return nil, err
	}

	d, err := FromConfig(cfg)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func FromConfig(c *Config) (*Database, error) {
	// Validate config
	err := c.validate()
	if err != nil {
		return nil, err
	}

	// Initialize all servers
	servers := make(map[string]*pgxpool.Pool)
	for serverName, serverObject := range c.Servers {
		// Connect
		conn, err := pgxpool.New(context.Background(), serverObject.connstr())
		if err != nil {
			return nil, err
		}

		// Ping
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = conn.Ping(ctx)
		if err != nil {
			return nil, err
		}

		servers[serverName] = conn
	}

	// Make a map to query and connection
	queries := make(map[string]query)
	for queryName, queryObject := range c.Queries {
		// Match to server
		queries[queryName] = query{
			query: queryObject.Query,
			conn:  servers[queryObject.Server],
		}
	}

	return &Database{
		servers: servers,
		queries: queries,
	}, nil
}

func (d *Database) HasQuery(name string) bool {
	_, ok := d.queries[name]
	return ok
}

func (d *Database) QueryRow(ctx context.Context, name string, args ...any) (pgx.Row, error) {
	query, ok := d.queries[name]
	if !ok {
		return nil, ErrNotFound
	}

	return query.conn.QueryRow(ctx, query.query, args...), nil
}

func (d *Database) Query(ctx context.Context, name string, args ...any) (pgx.Rows, error) {
	query, ok := d.queries[name]
	if !ok {
		return nil, ErrNotFound
	}

	return query.conn.Query(ctx, query.query, args...)
}

func (d *Database) CloseAll() {
	for _, conn := range d.servers {
		conn.Close()
	}
}

func (d *Database) Close(name string) error {
	conn, ok := d.servers[name]
	if !ok {
		return ErrNotFound
	}

	conn.Close()
	return nil
}
