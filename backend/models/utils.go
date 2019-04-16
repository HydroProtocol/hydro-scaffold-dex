package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"
	"reflect"
	"strconv"
	"strings"
)

type OrderByDirection uint

const (
	OrderByAsc OrderByDirection = iota
	OrderByDesc
)

const (
	OP_OR = "or"
)

func insert(model interface{}) (int64, error) {
	sqlString, values := getInsertSQLAndValues(model)
	ret, err := DB.Exec(sqlString, values...)

	if err != nil {
		return 0, err
	} else {
		// If you decide to use Postgres as your driver, please use `returning id` expression.
		// Please see more details at https://github.com/jmoiron/sqlx/issues/154
		id, err := ret.LastInsertId()

		if err != nil {
			return 0, err
		}

		return id, nil
	}
}

func update(model interface{}, fields ...string) error {
	sqlString, values := getUpdateSQLAndValues(model, fields...)
	_, err := DB.Exec(sqlString, values...)
	return err
}

// Will find a way to
func insertAndReturnID(model interface{}) (int64, error) {
	sqlString, values := getInsertSQLAndValues(model)
	sqlString = sqlString + ` RETURNING id`

	var id int64
	err := DB.QueryRow(sqlString, values...).Scan(&id)

	return id, err
}

func quote(str string) string {
	//return fmt.Sprintf(`"%s"`, str)
	return str
}

func getStructPrimaryKeyField(model interface{}) *reflect.StructField {
	t := reflect.ValueOf(model).Elem().Type()

	for i := 0; i < t.NumField(); i++ {
		if _, isPrimaryKey := t.Field(i).Tag.Lookup("primaryKey"); isPrimaryKey {
			f := t.Field(i)
			return &f
		}
	}

	return nil
}

func getStructFieldsList(model interface{}, excludeAuto bool) string {
	t := reflect.ValueOf(model).Elem().Type()
	var fields []string

	for i := 0; i < t.NumField(); i++ {
		if excludeAuto {
			if _, isAuto := t.Field(i).Tag.Lookup("autoIncrement"); isAuto {
				continue
			}
		}

		fields = append(fields, quote(t.Field(i).Tag.Get("db")))
	}
	return strings.Join(fields, ", ")
}

func getStructFieldsIndexList(model interface{}, excludeAuto bool) string {
	t := reflect.ValueOf(model).Elem().Type()
	var buff bytes.Buffer

	j := 0
	for i := 0; i < t.NumField(); i++ {
		if excludeAuto {
			if _, ok := t.Field(i).Tag.Lookup("autoIncrement"); ok {
				continue
			}
			j++
		}
	}

	for i := 0; i < j; i++ {
		buff.WriteString(fmt.Sprintf("$%d", i+1))
		if i < j-1 {
			buff.WriteString(", ")
		}
	}
	return buff.String()
}

func getValuesFromStruct(model interface{}, excludeAuto bool) []interface{} {
	e := reflect.ValueOf(model).Elem()
	t := e.Type()
	var fields []interface{}

	for i := 0; i < e.NumField(); i++ {
		if excludeAuto {
			if _, ok := t.Field(i).Tag.Lookup("autoIncrement"); ok {
				continue
			}

			// if t.Field(i).Type == reflect.TypeOf(time.Time{}) {
			// 	println("time found!!!!!", t.Field(i).Name, e.Field(i).Interface())
			// }
			fields = append(fields, e.Field(i).Interface())
		}
	}

	return fields
}

func getSelectSQLfromStruct(model interface{}) string {
	return fmt.Sprintf("SELECT %s FROM %s", getStructFieldsList(model, false), getTableNameFromStruct(model))
}

func getTableNameFromStruct(model interface{}) string {
	v := reflect.ValueOf(model)
	method := v.MethodByName("TableName")

	if method.IsValid() {
		res := method.Call([]reflect.Value{})
		return quote(res[0].String())
	}

	return quote(ToSnake(v.Elem().Type().Name()) + "s")
}

func getInsertSQLAndValues(model interface{}) (string, []interface{}) {
	s := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", getTableNameFromStruct(model), getStructFieldsList(model, true), getStructFieldsIndexList(model, true))
	values := getValuesFromStruct(model, true)
	return s, values
}

func getUpdateSQLAndValues(model interface{}, fields ...string) (string, []interface{}) {
	buffer := new(bytes.Buffer)
	values := []interface{}{}

	primaryKeyField := getStructPrimaryKeyField(model)

	if primaryKeyField == nil {
		panic(fmt.Errorf("can't update without primaryKey, model: %v", model))
	}

	if len(fields) == 0 {
		panic(fmt.Errorf("can't update without fields list, model: %v", model))
	}

	v := reflect.ValueOf(model).Elem()
	t := reflect.TypeOf(model).Elem()

	for index, key := range fields {
		if field, ok := t.FieldByName(key); ok {
			if dbFieldName, found := field.Tag.Lookup("db"); found {
				buffer.WriteString(quote(dbFieldName))
				buffer.WriteString(" = $")
				buffer.WriteString(strconv.Itoa(index + 1))

				field := v.FieldByName(key)
				values = append(values, field.Interface())
			} else {
				panic(fmt.Errorf("No db tag on field %s", key))
			}
		} else {
			panic(fmt.Errorf("No field %s", key))
		}

		if index < len(fields)-1 {
			buffer.WriteString(", ")
		}
	}

	// primaryKey
	values = append(values, v.FieldByName(primaryKeyField.Name).Interface())

	s := fmt.Sprintf("UPDATE %s SET %s WHERE %s = $%d", getTableNameFromStruct(model), buffer.String(), quote(primaryKeyField.Tag.Get("db")), len(values))
	return s, values
}

func getOrderSQL(orderBy map[string]OrderByDirection) string {
	var res []string
	for k, v := range orderBy {
		if v == OrderByAsc {
			res = append(res, fmt.Sprintf("%s ASC", quote(k)))
		} else {
			res = append(res, fmt.Sprintf("%s DESC", quote(k)))
		}
	}
	return strings.Join(res, ", ")
}

func getSelectSQLfromModels(models interface{}) string {
	// eg:
	// models                      Type *[]*Tx
	// models.Elem()               Type []*Tx
	// models.Elem().Elem()        Type *Tx
	// models.Elem().Elem().Elem() Type Tx
	t := reflect.TypeOf(models).Elem().Elem().Elem()

	// reflect.New(t)             Value *Tx
	// reflect.New(t).Interface() *Tx
	v := reflect.New(t).Interface()
	return getSelectSQLfromStruct(v)
}

func getConditions(conditions map[string]interface{}, startIndex int, connector string) (string, []interface{}) {
	keys := []string{}
	values := []interface{}{}

	for k := range conditions {
		keys = append(keys, k)
	}

	conditionStringBuffer := new(bytes.Buffer)

	for index, key := range keys {
		value := conditions[key]

		switch key {
		case OP_OR:
			s, vs := getConditions(value.(map[string]interface{}), startIndex+index, " OR ")
			conditionStringBuffer.WriteString(fmt.Sprintf("(%s)", s))
			values = append(values, vs...)
			startIndex += len(vs) - 1
		default:
			values = append(values, value)
			conditionStringBuffer.WriteString(quote(key))
			conditionStringBuffer.WriteString(" = $")
			conditionStringBuffer.WriteString(strconv.Itoa(startIndex + index))
		}

		if index < len(keys)-1 {
			conditionStringBuffer.WriteString(connector)
		}
	}

	return conditionStringBuffer.String(), values
}

func getWhereSQLfromConditions(conditions Op) (string, []interface{}) {
	return conditions.GetSqlAndValues(1)
}

func findBy(model interface{}, conditions Op, orderBy map[string]OrderByDirection) {
	sqlString := getSelectSQLfromStruct(model)

	whereSQL, values := getWhereSQLfromConditions(conditions)
	sqlString = fmt.Sprintf("%s WHERE %s", sqlString, whereSQL)

	if orderBy != nil {
		sqlString = fmt.Sprintf("%s ORDER BY %s", sqlString, getOrderSQL(orderBy))
	}

	sqlString = fmt.Sprintf("%s LIMIT 1", sqlString)

	// logSQL(sqlString)
	err := DB.QueryRowx(sqlString, values...).StructScan(model)

	if err == sql.ErrNoRows {
		return
	}
	if err != nil {
		log.Errorf("find by error: %v", err)
		panic(err)
	}
}

func findAllBy(models interface{}, conditions Op, orderBy map[string]OrderByDirection, limit, offset int) {
	sqlString := getSelectSQLfromModels(models)
	var values []interface{}
	if conditions != nil {
		whereSQL, vs := getWhereSQLfromConditions(conditions)
		values = vs
		sqlString = fmt.Sprintf("%s WHERE %s", sqlString, whereSQL)
	}

	if orderBy != nil {
		sqlString = fmt.Sprintf("%s ORDER BY %s", sqlString, getOrderSQL(orderBy))
	}

	if limit >= 0 {
		sqlString = fmt.Sprintf("%s LIMIT %d", sqlString, limit)
	}

	if offset >= 0 {
		sqlString = fmt.Sprintf("%s OFFSET %d", sqlString, offset)
	}

	//fmt.Println(sqlString)
	err := DB.Select(models, sqlString, values...)
	if err != nil {
		panic(err)
	}
}

func findCountBy(model interface{}, conditions Op) uint64 {
	sqlString := fmt.Sprintf("SELECT count(1) FROM %s", getTableNameFromStruct(model))

	var values []interface{}

	if conditions != nil {
		whereSQL, vs := getWhereSQLfromConditions(conditions)
		values = vs
		sqlString = fmt.Sprintf("%s WHERE %s", sqlString, whereSQL)
	}

	var count uint64
	err := DB.Get(&count, sqlString, values...)

	if err != nil {
		panic(err)
	}

	return count
}

type Op interface {
	GetSqlAndValues(index int) (string, []interface{})
}

type OpEq struct {
	Field string
	Value interface{}
}

func (op *OpEq) GetSqlAndValues(index int) (string, []interface{}) {
	return fmt.Sprintf("%s = $%d", op.Field, index), []interface{}{op.Value}
}

type OpGt struct {
	Field string
	Value interface{}
}

func (op *OpGt) GetSqlAndValues(index int) (string, []interface{}) {
	return fmt.Sprintf("%s > $%d", op.Field, index), []interface{}{op.Value}
}

type OpLt struct {
	Field string
	Value interface{}
}

func (op *OpLt) GetSqlAndValues(index int) (string, []interface{}) {
	return fmt.Sprintf("%s < $%d", op.Field, index), []interface{}{op.Value}
}

type OpList struct {
	ops       []Op
	connector string
}

func (opList *OpList) GetSqlAndValues(index int) (string, []interface{}) {
	var parts []string
	var resValues []interface{}

	for i := range opList.ops {
		partial, values := opList.ops[i].GetSqlAndValues(index + i)
		parts = append(parts, partial)
		resValues = append(resValues, values...)
		index += len(values) - 1
	}

	return fmt.Sprintf("(%s)", strings.Join(parts, opList.connector)), resValues
}

func whereAnd(ops ...Op) *OpList {
	return &OpList{
		connector: " AND ",
		ops:       ops,
	}
}

func whereOr(ops ...Op) *OpList {
	return &OpList{
		connector: " OR ",
		ops:       ops,
	}
}

func ToSnake(camel string) string {
	var snake []byte
	for i := range camel {
		if i == 0 {
			snake = append(snake, camel[i])
			continue
		}

		if isLowCase(int32(camel[i-1])) && !isLowCase(int32(camel[i])) {
			snake = append(snake, '_')
		}

		snake = append(snake, camel[i])
	}

	return strings.ToLower(string(snake))
}

func isLowCase(b int32) bool {
	if b >= 'a' && b <= 'z' {
		return true
	}

	if b >= '0' && b <= '9' {
		return true
	}

	return false
}
