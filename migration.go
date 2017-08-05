package pgmigration

import (
	"reflect"
	"log"
	"fmt"
	"unicode"
	"strings"
	"database/sql"
)

func Migration (values... interface{}, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	for _, v := range values {
		esql := NewMigration(&v)
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
		return tx.Commit()
	}
}

func NewMigration(value *interface{}) string {
	t := reflect.TypeOf(&value)
	tableName := TableName(t.Elem().Name())
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
func HasTable(table string,  *sql.DB) bool {

}
