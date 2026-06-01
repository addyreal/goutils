package psql

import (
	"testing"
)

func TestPsql(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		connstr map[string]string
		ok      bool
	}{
		{
			"ok all",
			Config{
				Servers: map[string]ServerObject{
					"main": {
						Host:                  "example.com",
						HostAddr:              "1.2.3.4",
						Port:                  1234,
						DbName:                "mydb",
						User:                  "myuser",
						Password:              "mypassword",
						ConnectTimeout:        10,
						SslMode:               "verify-full",
						SslRootCert:           "/path/to/cert",
						SslMinProtocolVersion: "TLSv1.2",
					},
				},
				Queries: map[string]QueryObject{
					"my-query": {
						Server: "main",
						Query:  "SELECT * FROM table;",
					},
				},
			},
			map[string]string{
				"main": "host=example.com hostaddr=1.2.3.4 port=1234 dbname=mydb user=myuser password=mypassword connect_timeout=10 sslmode=verify-full sslrootcert=/path/to/cert ssl_min_protocol_version=TLSv1.2",
			},
			true,
		},
		{
			"missing host but not hostaddr",
			Config{
				Servers: map[string]ServerObject{
					"main": {
						HostAddr:              "1.2.3.4",
						Port:                  1234,
						DbName:                "mydb",
						User:                  "myuser",
						Password:              "mypassword",
						ConnectTimeout:        10,
						SslMode:               "verify-full",
						SslRootCert:           "/path/to/cert",
						SslMinProtocolVersion: "TLSv1.2",
					},
				},
				Queries: map[string]QueryObject{
					"my-query": {
						Server: "main",
						Query:  "SELECT * FROM table;",
					},
				},
			},
			map[string]string{
				"main": "hostaddr=1.2.3.4 port=1234 dbname=mydb user=myuser password=mypassword connect_timeout=10 sslmode=verify-full sslrootcert=/path/to/cert ssl_min_protocol_version=TLSv1.2",
			},
			true,
		},
		{
			"missing host and hostaddr",
			Config{
				Servers: map[string]ServerObject{
					"main": {
						Port:                  1234,
						DbName:                "mydb",
						User:                  "myuser",
						Password:              "mypassword",
						ConnectTimeout:        10,
						SslMode:               "verify-full",
						SslRootCert:           "/path/to/cert",
						SslMinProtocolVersion: "TLSv1.2",
					},
				},
				Queries: map[string]QueryObject{
					"my-query": {
						Server: "main",
						Query:  "SELECT * FROM table;",
					},
				},
			},
			map[string]string{
				"main": "",
			},
			false,
		},
		{
			"wrong sslmode",
			Config{
				Servers: map[string]ServerObject{
					"main": {
						Host:                  "example.com",
						HostAddr:              "1.2.3.4",
						Port:                  1234,
						DbName:                "mydb",
						User:                  "myuser",
						Password:              "mypassword",
						ConnectTimeout:        10,
						SslMode:               "please",
						SslRootCert:           "/path/to/cert",
						SslMinProtocolVersion: "TLSv1.2",
					},
				},
				Queries: map[string]QueryObject{
					"my-query": {
						Server: "main",
						Query:  "SELECT * FROM table;",
					},
				},
			},
			map[string]string{
				"main": "",
			},
			false,
		},
		{
			"wrong ssl version",
			Config{
				Servers: map[string]ServerObject{
					"main": {
						Host:                  "example.com",
						HostAddr:              "1.2.3.4",
						Port:                  1234,
						DbName:                "mydb",
						User:                  "myuser",
						Password:              "mypassword",
						ConnectTimeout:        10,
						SslMode:               "verify-full",
						SslRootCert:           "/path/to/cert",
						SslMinProtocolVersion: "tlsv1.2",
					},
				},
				Queries: map[string]QueryObject{
					"my-query": {
						Server: "main",
						Query:  "SELECT * FROM table;",
					},
				},
			},
			map[string]string{
				"main": "",
			},
			false,
		},
		{
			"wrong server ref",
			Config{
				Servers: map[string]ServerObject{
					"main": {
						Host:                  "example.com",
						HostAddr:              "1.2.3.4",
						Port:                  1234,
						DbName:                "mydb",
						User:                  "myuser",
						Password:              "mypassword",
						ConnectTimeout:        10,
						SslMode:               "verify-full",
						SslRootCert:           "/path/to/cert",
						SslMinProtocolVersion: "tlsv1.2",
					},
				},
				Queries: map[string]QueryObject{
					"my-query": {
						Server: "nonexistent",
						Query:  "SELECT * FROM table;",
					},
				},
			},
			map[string]string{
				"main": "",
			},
			false,
		},
		{
			"zero port",
			Config{
				Servers: map[string]ServerObject{
					"main": {
						Host:                  "example.com",
						HostAddr:              "1.2.3.4",
						Port:                  0,
						DbName:                "mydb",
						User:                  "myuser",
						Password:              "mypassword",
						ConnectTimeout:        10,
						SslMode:               "verify-full",
						SslRootCert:           "/path/to/cert",
						SslMinProtocolVersion: "TLSv1.2",
					},
				},
				Queries: map[string]QueryObject{
					"my-query": {
						Server: "main",
						Query:  "SELECT * FROM table;",
					},
				},
			},
			map[string]string{
				"main": "",
			},
			false,
		},
		{
			"zero connect_timeout",
			Config{
				Servers: map[string]ServerObject{
					"main": {
						Host:                  "example.com",
						HostAddr:              "1.2.3.4",
						Port:                  1234,
						DbName:                "mydb",
						User:                  "myuser",
						Password:              "mypassword",
						ConnectTimeout:        0,
						SslMode:               "verify-full",
						SslRootCert:           "/path/to/cert",
						SslMinProtocolVersion: "TLSv1.2",
					},
				},
				Queries: map[string]QueryObject{
					"my-query": {
						Server: "main",
						Query:  "SELECT * FROM table;",
					},
				},
			},
			map[string]string{
				"main": "host=example.com hostaddr=1.2.3.4 port=1234 dbname=mydb user=myuser password=mypassword connect_timeout=0 sslmode=verify-full sslrootcert=/path/to/cert ssl_min_protocol_version=TLSv1.2",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if err != nil {
				if tt.ok {
					t.Errorf(`%s: failed. Expected valid`, tt.name)
				}

				t.SkipNow()
			} else {
				if !tt.ok {
					t.Errorf(`%s: failed. Expected invalid`, tt.name)
				}
			}

			for n, s := range tt.config.Servers {
				if tt.connstr[n] != s.connstr() {
					t.Errorf(`%s: failed. Expected connstr "%s", got "%s"`, tt.name, tt.connstr[n], s.connstr())
				}
			}
		})
	}
}
