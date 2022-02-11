package core

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"regexp"
	"text/template"
)

const (
	Postgresql string = "postgresql"
	Mysql             = "mysql"
	SqlServer         = "sqlserver"
	Oracle			  = "oracle"
)

var pgTransactionTemplate, _ = template.New("pgTransaction").Parse(`{{define "pgTransaction"}}BEGIN;

{{.}}
COMMIT;
{{end}}`)

var SqlConnectionUrlRegex = regexp.MustCompile("^([a-z]+?):\\/\\/(.+?):(.+?)@([\\w:\\.]+?)\\/([\\w]+?)([\\?].+?)?$")

func GetDatabaseConnection(url string) (*sql.DB, error) {
	result := SqlConnectionUrlRegex.FindAllStringSubmatch(url, -1)

	if result == nil || len(result[0]) <= 1 {
		return nil, errors.New("bad format")
	}

	fmt.Println("Connecting " + "(" + result[0][1] + ")" + " to: " + url)
	fmt.Println()

	open, err := sql.Open(result[0][1], url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return open, nil
}

func CreateDatabaseVersionTable(url string) error {
	DB, err := GetDatabaseConnection(url)

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			fmt.Println("could not close database connection")
		}
	}(DB)

	_, err = DB.Query("CREATE TABLE walkline_version (version VARCHAR)")

	if err != nil {
		return err
	}

	return nil
}

func GetCurrentDatabaseVersion(url string) (string, error) {
	DB, err := GetDatabaseConnection(url)
	var version string

	if err != nil {
		return "", err
	}

	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			fmt.Println("could not close database connection")
		}
	}(DB)

	row := DB.QueryRow("SELECT version FROM walkline_version")

	err = row.Scan(&version)

	if err != nil {
		return "", err
	}

	return version, nil
}

func GetSetDatabaseVersionQueryString(version *Version) string {
	return "INSERT INTO walkline_version (version) VALUES ('" + version.Prefix + version.Version + "')"
}

func GenerateTransactionString(flavor string, sql string) (string, error) {
	var out bytes.Buffer
	if flavor == Postgresql {
		err := pgTransactionTemplate.ExecuteTemplate(&out, "pgTransaction", sql)
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("invalid flavor")
	}

	return out.String(), nil
}
