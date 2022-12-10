package model

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"regexp"
	"strings"
)

type queryCondition interface{}

const (
	// below name should be same as the struct "types.DynamicQueryRequest" field's tag name
	orderFieldName      = "orderby"
	limitFieldName      = "limit"
	offsetFieldName     = "offset"
	sortFieldName       = "asc"
	fuzzyQueryFieldName = "fuzzyQuery"

	// will be used in "update/delete" operation
	primaryKey = "id"
)

// see base_test.go
func buildFuzzyQueryCondition(query queryCondition) []string {
	v := reflect.ValueOf(query)
	var queryStatement []string
	condition := []string{""}

	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			queryStatement = append(queryStatement, v.Type().Field(i).Name)
			if field.Kind() == reflect.String {
				condition = append(condition, fmt.Sprintf("%%%s%%", field.String()))
			} else if field.Kind() == reflect.Int {
				condition = append(condition, fmt.Sprintf("%%%d%%", field.Int()))
			} else {
				// should not happen
			}
		}
	} else if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)

			queryStatement = append(queryStatement, key.String())
			underValue := value.Interface()
			if reflect.TypeOf(underValue).Kind() == reflect.String {
				condition = append(condition, "%"+underValue.(string)+"%")
			} else if reflect.TypeOf(underValue).Kind() == reflect.Int {
				condition = append(condition, "%"+fmt.Sprintf("%d", underValue.(int))+"%")
			} else {
				// should not happen
			}
		}
	}

	if len(queryStatement) > 1 {
		condition[0] = strings.Join(queryStatement, " like ? and ") + " like ?"
	} else if len(queryStatement) == 1 {
		condition[0] = queryStatement[0] + " like ?"
	}

	return condition
}

// queryCondition is a map or struct
func buildQueryCondition(query queryCondition) []string {
	v := reflect.ValueOf(query)
	var queryStatement []string
	condition := []string{""}

	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			queryStatement = append(queryStatement, v.Type().Field(i).Name)
			if field.Kind() == reflect.String {
				condition = append(condition, fmt.Sprintf("%s", field.String()))
			} else if field.Kind() == reflect.Int {
				condition = append(condition, fmt.Sprintf("%d", field.Int()))
			} else {
				// should not happen
			}
		}
	} else if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)

			queryStatement = append(queryStatement, key.String())
			underValue := value.Interface()
			if reflect.TypeOf(underValue).Kind() == reflect.String {
				condition = append(condition, underValue.(string))
			} else if reflect.TypeOf(underValue).Kind() == reflect.Int {
				condition = append(condition, fmt.Sprintf("%d", underValue.(int)))
			} else {
				// should not happen
			}
		}
	}

	if len(queryStatement) > 1 {
		condition[0] = strings.Join(queryStatement, " = ? and ") + " = ?"
	} else if len(queryStatement) == 1 {
		condition[0] = queryStatement[0] + " = ?"
	}

	return condition
}

// support query by dynamic condition
// can query by null value, such as "id=&name="
func parseDynamicQueryRequest(query map[string]interface{}, fuzzyQuery bool) ([]string, error) {

	var queryCondition = make(map[string]string)

	v := reflect.ValueOf(query)

	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)
		underValue := value.Interface()
		if key.String() == orderFieldName {
			// do something
		} else if key.String() == limitFieldName {

		} else if key.String() == offsetFieldName {

		} else if key.String() == sortFieldName {

		} else if key.String() == fuzzyQueryFieldName {

		} else {

			if reflect.TypeOf(underValue).Kind() == reflect.String {
				if underValue.(string) != "" {
					queryCondition[key.String()] = underValue.(string)
				}
			} else if reflect.TypeOf(underValue).Kind() == reflect.Int {
				if underValue.(int) != 0 {
					queryCondition[key.String()] = fmt.Sprintf("%d", underValue.(int))
				}
			} else {
				// should not happen
			}
		}
	}

	where := make([]string, 0)
	if fuzzyQuery {
		where = buildFuzzyQueryCondition(queryCondition)
	} else {
		where = buildQueryCondition(queryCondition)
	}
	return where, nil
}

type DynamicModel struct {
	db *gorm.DB
}

func NewDynamicModel(db *gorm.DB) DynamicModel {
	return DynamicModel{
		db: db,
	}
}

func (m *DynamicModel) buildWhereCondition(query map[string]interface{}, fuzzyQuery bool) (*gorm.DB, error) {
	where, err := parseDynamicQueryRequest(query, fuzzyQuery)
	if err != nil {
		return nil, err
	}

	if len(where) > 1 {
		// []string can not assign to []interface{}
		tmp := make([]interface{}, 0)
		for i := 0; i < len(where); i++ {
			tmp = append(tmp, where[i])
		}
		return m.db.Where(tmp[0], tmp[1:]...), nil
	} else if len(where) == 1 {
		return m.db.Where(where[0]), nil
	} else {
		return m.db, nil
	}
}

func buildOrder(order string, asc int) string {
	// check sql injection
	compile, err := regexp.Compile("^[a-zA-Z0-9_]+$")
	if err != nil {
		return ""
	}
	if !compile.Match([]byte(order)) {
		return ""
	}
	if asc == 0 {
		return order + " desc"
	} else {
		return order + " asc"
	}
}

func (m *DynamicModel) Query(tableName string, limit int, offset int, orderby string, asc int,
	fuzzyQuery bool, condition map[string]interface{}) (int64, []map[string]interface{}, error) {
	// support fuzzy query
	// support pagination
	// support order by

	// for safety
	var count int64
	dest := make([]map[string]interface{}, 0)

	db, err := m.buildWhereCondition(condition, fuzzyQuery)
	if err != nil {
		return 0, nil, err
	}

	db.Table(tableName).Count(&count)
	db.Table(tableName).Limit(limit).Offset(offset).Order(buildOrder(orderby, asc)).Find(&dest)

	return count, dest, nil
}

// InsertOne support insert one record
func (m *DynamicModel) InsertOne(tableName string, data map[string]interface{}) error {
	retInsertId := m.db.Table(tableName).Create(data)
	if retInsertId.Error != nil {
		return retInsertId.Error
	}
	return nil
}

// Update support update one record
func (m *DynamicModel) Update(tableName string, data map[string]interface{}) error {
	m.db.Table(tableName).Where(fmt.Sprintf("%s=?", primaryKey), data[primaryKey]).Updates(&data) // update all fields
	return nil
}

// DeleteByQuery return delete count
// delete rows which match the query condition. if query is nil, delete all rows
// query is a map, key is column name, value is column value. if value is nil, do nothing
// fuzzyQuery is true, use like query. if fuzzyQuery is false, use equal query
func (m *DynamicModel) DeleteByQuery(tableName string, fuzzyQuery bool, condition map[string]interface{}) (int64, error) {

	db, err := m.buildWhereCondition(condition, fuzzyQuery)
	if err != nil {
		return 0, err
	}
	deleteRows := make(map[string]interface{})
	db.Table(tableName).Delete(&deleteRows) // deletedRows seems to be meaningsless here
	return db.RowsAffected, nil             // who set rowsAffected? mysql server return affected rows?

	// ROW_COUNT() can return delete count, but it does not work in here. I don't know why
	// https://dev.mysql.com/doc/refman/8.0/en/information-functions.html#function_row-count
	//db.Raw("SELECT ROW_COUNT()").Scan(&deleteRows)
	//return deleteRows["ROW_COUNT()"].(int64), nil
}

// DeleteById return delete count
// delete rows which match the id
func (m *DynamicModel) DeleteById(tableName string, id int) (int64, error) {
	var count int64
	m.db.Table(tableName).Where(fmt.Sprintf("%s=?", primaryKey), id).Count(&count)
	m.db.Table(tableName).Where(fmt.Sprintf("%s=?", primaryKey), id).Delete(nil)
	return count, nil
}

// DeleteByIds return delete count
// delete rows which match the ids
func (m *DynamicModel) DeleteByIds(tableName string, ids []int) (int64, error) {
	var count int64
	m.db.Table(tableName).Where(fmt.Sprintf("%s in (?)", primaryKey), ids).Count(&count)
	m.db.Table(tableName).Where(fmt.Sprintf("%s in (?)", primaryKey), ids).Delete(nil)
	return count, nil
}

// DeleteAll return delete count
// delete all rows
func (m *DynamicModel) DeleteAll(tableName string) (int64, error) {
	var count int64
	m.db.Table(tableName).Count(&count)
	m.db.Table(tableName).Delete(nil)
	return count, nil
}

// QueryById return one record
func (m *DynamicModel) QueryById(tableName string, id int) (map[string]interface{}, error) {
	var dest = make(map[string]interface{})
	m.db.Table(tableName).Where(fmt.Sprintf("%s=?", primaryKey), id).Find(&dest)
	return dest, nil
}

// QueryByIds return multi records
func (m *DynamicModel) QueryByIds(tableName string, ids []int) ([]map[string]interface{}, error) {
	var dest = make([]map[string]interface{}, 0)
	m.db.Table(tableName).Where(fmt.Sprintf("%s in (?)", primaryKey), ids).Find(&dest)
	return dest, nil
}
