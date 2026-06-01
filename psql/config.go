// Example config:
//
//   servers:
//     auth:
//       host: "example.com"
//       port: 1234
//       dbname: "mydb"
//       user: "myuser"
//       password: "mypassword"
//       connect_timeout: 10
//     main:
//       hostaddr: "2.3.4.5"
//       port: 2345
//       dbname: "somedb"
//       user: "someuser"
//       password: "somepassword"
//       sslmode: "verify-full"
//       sslrootcert: "/etc/ssl/certs/ca-certificates.crt"
//       ssl_min_protocol_version: "TLSv1.2"
//   
//   queries:
//     get-basic-auth:
//       server: "auth"
//       query: "SELECT * FROM public.users;"
//     get-thing:
//       server: "main"
//       query: "SELECT thing FROM public.stuff;"
package psql

import (
	"errors"
	"fmt"
	"github.com/addyreal/goutils/psql/internal/pgsql"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Config struct {
	Servers map[string]ServerObject `yaml:"servers"`
	Queries map[string]QueryObject  `yaml:"queries"`
}

// Unmarshals a config file
func configFromYaml(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Make strict decoder
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)

	// Unmarshal
	var res Config
	err = dec.Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Config) validate() error {
	// Validate servers
	for name, server := range c.Servers {
		if server.Host == "" && server.HostAddr == "" {
			return errors.New("missing host or hostaddr for server " + name)
		}

		if server.Port == 0 {
			return errors.New("missing port for server " + name)
		}

		if server.DbName == "" {
			return errors.New("missing dbname for server " + name)
		}

		if server.User == "" {
			return errors.New("missing user for server " + name)
		}

		if server.Password == "" {
			return errors.New("missing password for server " + name)
		}

		_, ok := pgsql.SslModeMap[server.SslMode]
		if !ok {
			return errors.New("invalid sslmode '" + server.SslMode + "' for server " + name)
		}

		_, ok = pgsql.SslVersionMap[server.SslMinProtocolVersion]
		if !ok {
			return errors.New("invalid ssl_min_protocol_version '" + server.SslMinProtocolVersion + "' for server " + name)
		}
	}

	// Validate queries
	for name, query := range c.Queries {
		_, ok := c.Servers[query.Server]
		if !ok {
			return errors.New("unknown server '" + query.Server + "' for query " + name)
		}
	}

	return nil
}

type ServerObject struct {
	Host                  string `yaml:"host"`
	HostAddr              string `yaml:"hostaddr"`
	Port                  uint64 `yaml:"port"`
	DbName                string `yaml:"dbname"`
	User                  string `yaml:"user"`
	Password              string `yaml:"password"`
	ConnectTimeout        uint64 `yaml:"connect_timeout"`
	SslMode               string `yaml:"sslmode"`
	SslRootCert           string `yaml:"sslrootcert"`
	SslMinProtocolVersion string `yaml:"ssl_min_protocol_version"`
}

func (s *ServerObject) connstr() string {
	var b strings.Builder

	if s.Host != "" {
		fmt.Fprintf(
			&b, "host=%s ",
			s.Host,
		)
	}

	if s.HostAddr != "" {
		fmt.Fprintf(
			&b, "hostaddr=%s ",
			s.HostAddr,
		)
	}

	fmt.Fprintf(
		&b, "port=%d dbname=%s user=%s password=%s connect_timeout=%d ",
		s.Port, s.DbName, s.User, s.Password, s.ConnectTimeout,
	)

	if s.SslMode != "" {
		fmt.Fprintf(
			&b, "sslmode=%s ",
			s.SslMode,
		)
	}

	if s.SslRootCert != "" {
		fmt.Fprintf(
			&b, "sslrootcert=%s ",
			s.SslRootCert,
		)
	}

	if s.SslMinProtocolVersion != "" {
		fmt.Fprintf(
			&b, "ssl_min_protocol_version=%s ",
			s.SslMinProtocolVersion,
		)
	}

	return strings.TrimSpace(b.String())
}

type QueryObject struct {
	Server string `yaml:"server"`
	Query  string `yaml:"query"`
}
