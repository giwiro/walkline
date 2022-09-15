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

func GetVersionTableName(schema string) string {
    var tableName = "walkline_version"

    if schema != "" {
        tableName = fmt.Sprintf("%s.walkline_version", schema)
    }

    return tableName
}

func GetDatabaseConnection(url string, verbose bool) (*sql.DB, string, error) {
    result := SqlConnectionUrlRegex.FindAllStringSubmatch(url, -1)

    if result == nil || len(result[0]) <= 1 {
        return nil, "", errors.New("connection url bad format")
    }

    var flavor = result[0][1]

    if verbose == true {
        fmt.Printf("Connecting (%s) to: %s\n", flavor, url)
    }

    open, err := sql.Open(flavor, url)

    if err != nil {
        fmt.Println(err)
        return nil, flavor, err
    }

    return open, flavor, nil
}

func GetCreateVersionTableQueryString(schema string) string {
    var tableName = GetVersionTableName(schema)
    return fmt.Sprintf(
        "CREATE TABLE %s (version VARCHAR);\n"+
            "INSERT INTO %s (version) VALUES ('');\n\n",
        tableName,
        tableName,
    )
}

func CreateVersionTable(url string, verbose bool, schema string) error {
    var crateTableQueryString = GetCreateVersionTableQueryString(schema)

    DB, _, err := GetDatabaseConnection(url, verbose)

    if err != nil {
        return err
    }

    defer func(DB *sql.DB) {
        err := DB.Close()
        if err != nil {
            fmt.Println("could not close database connection")
        }
    }(DB)

    _, err = DB.Query(crateTableQueryString)

    if err != nil {
        return err
    }

    return nil
}

func GetCurrentDatabaseVersion(url string, verbose bool, schema string) (*VersionShort, string, error) {
    var tableName = GetVersionTableName(schema)

    DB, flavor, err := GetDatabaseConnection(url, verbose)
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

    row := DB.QueryRow(fmt.Sprintf("SELECT version FROM %s", tableName))

    err = row.Scan(&version)

    if err != nil && err != sql.ErrNoRows {
        return nil, flavor, err
    }

    if len(version) == 0 {
        return nil, flavor, nil
    }

    versionShort, err := ParseVersionShort(version)

    if err != nil {
        return nil, flavor, err
    }

    return versionShort, flavor, nil
}

func GetUpdateVersionQueryString(init bool, currentVersion *VersionShort, version *VersionShort, schema string) string {
    var tableName = GetVersionTableName(schema)

    const updateFmt = "UPDATE %s SET version='%s' WHERE version='%s';"

    if currentVersion == nil || init {
        return fmt.Sprintf(updateFmt, tableName, version.Prefix+version.Version, "")
    }

    if version == nil {
        return fmt.Sprintf(updateFmt, tableName, "", currentVersion.Prefix+currentVersion.Version)
    }

    return fmt.Sprintf(updateFmt, tableName, version.Prefix+version.Version, currentVersion.Prefix+currentVersion.Version)
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

func ExecuteMigrationString(url string, sqlString string, verbose bool) error {
    DB, _, err := GetDatabaseConnection(url, verbose)

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
}
