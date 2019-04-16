package models

import (
	"github.com/HydroProtocol/hydro-sdk-backend/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type utilsTestSuite struct {
	suite.Suite
}

type Foo struct {
	ID    string          `db:"id" primaryKey:"true"`
	Name  string          `db:"name"`
	Age   int             `db:"age"`
	Birth *time.Time      `db:"birth"`
	Iq    decimal.Decimal `db:"iq"`
}

func (s *utilsTestSuite) SetupSuite() {
	test.PreTest()
	InitTestDB()
}

func (s *utilsTestSuite) TearDownSuite() {
}

func (s *utilsTestSuite) TearDownTest() {
}

func (s *utilsTestSuite) TestGetSelectSQLfromStruct() {
	sql := getSelectSQLfromStruct(&Foo{})
	s.Equal(`SELECT id, name, age, birth, iq FROM foos`, sql)
}

func (s *utilsTestSuite) TestGetInsertSQLfromStruct() {
	sql, _ := getInsertSQLAndValues(&Foo{})
	s.Equal(`INSERT INTO foos (id, name, age, birth, iq) VALUES ($1, $2, $3, $4, $5)`, sql)
}

func (s *utilsTestSuite) TestGetUpdateSQLfromStruct() {
	sql, _ := getUpdateSQLAndValues(&Foo{ID: "3"}, "Name")
	s.Equal(`UPDATE foos SET name = $1 WHERE id = $2`, sql)
}

func (s *utilsTestSuite) TestWhere() {

	sql, values := getWhereSQLfromConditions(whereAnd(
		&OpEq{"a", 3},
		&OpEq{"b", 4},
		whereOr(
			&OpEq{"c", 5},
			&OpEq{"d", 4},
			whereAnd(
				&OpEq{"ee", 5},
				&OpEq{"bb", 5},
			),
		),
		&OpEq{"e", 10},
	))

	s.Equal(`(a = $1 AND b = $2 AND (c = $3 OR d = $4 OR (ee = $5 AND bb = $6)) AND e = $7)`, sql)
	s.Equal([]interface{}{3, 4, 5, 4, 5, 5, 10}, values)
}

func (s *utilsTestSuite) TestGetValuesFromStruct() {
	t := time.Now().UTC()
	foo := &Foo{
		"1", "test", 12, &t, decimal.New(1, 0),
	}

	values := getValuesFromStruct(foo, true)

	s.Equal("1", values[0])
	s.Equal("test", values[1])
	s.Equal(12, values[2])
	s.Equal(&t, values[3])
	s.Equal("1", values[4].(decimal.Decimal).String())
}

type Demo struct{}

type Demo2 struct{}

func (f *Demo2) TableName() string { return "Yeah" }

func (s *utilsTestSuite) TestTableName() {
	s.Equal(`demos`, getTableNameFromStruct(&Demo{}))
	s.Equal(`Yeah`, getTableNameFromStruct(&Demo2{}))
}

func (s *utilsTestSuite) TestFindID() {
	// fmt.Printf("%+v\n", OrderDAO.FindByID(DB, "a"))
	// // fmt.Printf("%+v\n", OrderDAO.FindAllByID(DB, "a")[0])
	// fmt.Printf("%+v\n", OrderDAO.FindMarketPending(DB, "LMM-WETH")[0])
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(utilsTestSuite))
}

func TestToSnake(t *testing.T) {
	camel1 := "MyNameIsKEVIN"
	camel2 := ""
	camel3 := "myNameIsKEVIN"
	camel4 := "string2byte"
	camel5 := "String2Byte"

	assert.EqualValues(t, "my_name_is_kevin", ToSnake(camel1))
	assert.EqualValues(t, "", ToSnake(camel2))
	assert.EqualValues(t, "my_name_is_kevin", ToSnake(camel3))
	assert.EqualValues(t, "string2byte", ToSnake(camel4))
	assert.EqualValues(t, "string2_byte", ToSnake(camel5))
}
