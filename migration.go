package pgmigration

import (
	"reflect"
	"log"
	"fmt"
	"unicode"
	"strings"
	"database/sql"
	_ "github.com/lib/pq"
)

type MigServ struct {
	Db *sql.DB
}

func (m *MigServ) Migration(values ...interface{}) error {

	tx, err := m.Db.Begin()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	for _, value := range values {
		log.Printf("value:%+v\n", value)
		esql := NewMigration(value, m.Db)
		log.Println(esql)
		if esql == "" {
			continue
		}
		st, err := tx.Prepare(esql)
		if err != nil {
			log.Printf("%+v\n", err)
			tx.Rollback()
			return err
		}
		_, err = st.Exec()
		if err != nil {
			log.Printf("%+v\n", err)
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func NewMigration(value interface{}, db *sql.DB) string {
	t := reflect.TypeOf(value)
	tableName := TableName(t.Elem().Name())
	if HasTable(tableName, db) {
		return ""
	}

	esql := fmt.Sprintf(`CREATE TABLE "%s" (`, tableName)
	for i := 0; i < t.Elem().NumField(); i++ {
		//字段名.
		filed := t.Elem().Field(i).Tag.Get("db")
		sArr := strings.Split(filed, ", ")
		esql += fmt.Sprintf("%s, ", strings.Join(sArr, " "))
	}
	return strings.TrimSuffix(esql, ", ") + ")"
}

//返回表名.
func TableName(name string) string {
	tableName := ""
	for k, v := range name {
		if unicode.IsUpper(v) {
			if k == 0 {
				tableName += string(v)
			} else {
				tableName += fmt.Sprintf("_%s", string(v))
			}
		} else {
			tableName += string(v)
		}
	}
	return strings.ToLower(tableName)
}

//判断表是否存在.
func HasTable(table string, db *sql.DB) bool{
	esql := fmt.Sprintf("SELECT to_regclass('%s') is not null", table)
	row := db.QueryRow(esql)
	var isTrue bool
	if err := row.Scan(&isTrue); err != nil {
		log.Printf("%+v\n", err, esql)
		return false
	}
	return isTrue
}
