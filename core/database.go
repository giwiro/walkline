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
	Postgresql string = "postgres"
	Mysql             = "mysql"
	SqlServer         = "sqlserver"
	Oracle            = "oracle"
)

var pgTransactionTemplate, _ = template.New("pgTransaction").Parse(`{{define "pgTransaction"}}BEGIN;

{{.}}
COMMIT;
{{end}}`)

var SqlConnectionUrlRegex = regexp.MustCompile("^([a-z]+?):\\/\\/(.+?):(.+?)@([\\w:\\.]+?)\\/([\\w_]+?)([\\?].+?)?$")

func GetDatabaseConnection(url string) (*sql.DB, string, error) {
	result := SqlConnectionUrlRegex.FindAllStringSubmatch(url, -1)

	if result == nil || len(result[0]) <= 1 {
		return nil, "", errors.New("connection url bad format")
	}

	var flavor = result[0][1]

	fmt.Println("Connecting " + "(" + flavor + ")" + " to: " + url)

	open, err := sql.Open(flavor, url)
	if err != nil {
		fmt.Println(err)
		return nil, flavor, err
	}

	return open, flavor, nil
}

func CreateDatabaseVersionTable(url string) error {
	DB, _, err := GetDatabaseConnection(url)

	if err != nil {
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

func GetCurrentDatabaseVersion(url string) (*VersionShort, string, error) {
	DB, flavor, err := GetDatabaseConnection(url)
	var version string

	if err != nil {
		return nil, flavor, err
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
		return nil, flavor, err
	}

	versionShort, err := ParseVersionShort(version)

	if err != nil {
		return nil, flavor, err
	}

	return versionShort, flavor, nil
}

func GetInsertVersionQueryString(currentVersion *VersionShort, version *VersionShort) string {
	if currentVersion == nil {
		return "INSERT INTO walkline_version (version) VALUES ('" + version.Prefix + version.Version + "');"
	}
	return "UPDATE walkline_version SET version='" + version.Prefix + version.Version + "' WHERE version='" + currentVersion.Prefix + currentVersion.Version + "';"
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

func ExecuteMigrationString(url string, sqlString string) error {
	DB, _, err := GetDatabaseConnection(url)

	if err != nil {
		return err
	}

	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			fmt.Println("could not close database connection")
		}
	}(DB)

	/*	ctx := context.Background()

		tx, err := DB.BeginTx(ctx, nil)

		if err != nil {
			return err
		}*/

	_, err = DB.Exec(sqlString)

	if err != nil {
		return err
	}

	return nil

	/*_, txErr := tx.ExecContext(ctx, sqlString)

	if txErr != nil {
		err := tx.Rollback()

		if err != nil {
			return err
		}

		return txErr
	}

	err = tx.Commit()

	if err != nil {
		return err
	}*/

	return nil
}
