package querymaker

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if !reflect.DeepEqual(a, b) {
				t.Errorf("Expected %#v (type %v) - Got %#v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
			}
		}
	}()
	if a != b {
		t.Errorf("Expected %#v (type %v) - Got %#v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestCheatSheet(t *testing.T) {
	func() {
		query := SqlEq("foo", "bar")
		expect(t, query.AsSql(), "`foo` = ?")
		expect(t, query.Bind(), []interface{}{"bar"})
	}()
	func() {
		query := SqlIn("foo", Array{"bar", "baz"})
		expect(t, query.AsSql(), "`foo` IN (?,?)")
		expect(t, query.Bind(), []interface{}{"bar", "baz"})
	}()
	func() {
		query := SqlAnd(Array{SqlEq("foo", "bar"), SqlEq("baz", 123)})
		expect(t, query.AsSql(), "(`foo` = ?) AND (`baz` = ?)")
		expect(t, query.Bind(), []interface{}{"bar", 123})
	}()
	func() {
		query := SqlAnd("foo", Array{SqlGe(3), SqlLt(5)})
		expect(t, query.AsSql(), "(`foo` >= ?) AND (`foo` < ?)")
		expect(t, query.Bind(), []interface{}{3, 5})
	}()
	func() {
		query := SqlOr(Array{SqlEq("foo", "bar"), SqlEq("baz", 123)})
		expect(t, query.AsSql(), "(`foo` = ?) OR (`baz` = ?)")
		expect(t, query.Bind(), []interface{}{"bar", 123})
	}()
	func() {
		query := SqlOr("foo", Array{"bar", "baz"})
		expect(t, query.AsSql(), "(`foo` = ?) OR (`foo` = ?)")
		expect(t, query.Bind(), []interface{}{"bar", "baz"})
	}()
	func() {
		query := SqlIsNull("foo")
		expect(t, query.AsSql(), "`foo` IS NULL")
		expect(t, query.Bind(), []interface{}{})
	}()
	func() {
		query := SqlIsNotNull("foo")
		expect(t, query.AsSql(), "`foo` IS NOT NULL")
		expect(t, query.Bind(), []interface{}{})
	}()
	func() {
		query := SqlBetween("foo", 1, 2)
		expect(t, query.AsSql(), "`foo` BETWEEN ? AND ?")
		expect(t, query.Bind(), []interface{}{1, 2})
	}()
	func() {
		query := SqlNot("foo")
		expect(t, query.AsSql(), "NOT `foo`")
		expect(t, query.Bind(), []interface{}{})
	}()
	func() {
		query := SqlOp("apples", "MATCH (@) AGAINST (?)", Array{"oranges"})
		expect(t, query.AsSql(), "MATCH (`apples`) AGAINST (?)")
		expect(t, query.Bind(), []interface{}{"oranges"})
	}()
	func() {
		query := SqlRaw("SELECT * FROM t WHERE id=?", 123)
		expect(t, query.AsSql(), "SELECT * FROM t WHERE id=?")
		expect(t, query.Bind(), []interface{}{123})
	}()
	func() {
		query := SqlIn("foo", Array{123, SqlRaw("SELECT id FROM t WHERE cat=?", 5)})
		expect(t, query.AsSql(), "`foo` IN (?,(SELECT id FROM t WHERE cat=?))")
		expect(t, query.Bind(), []interface{}{123, 5})
	}()
}

func TestAndUsingHash(t *testing.T) {
	func() {
		query := SqlAnd(NewOrderedHash(Hash{
			"foo": 1,
			"bar": SqlEq(2),
			"baz": SqlLt(3),
		}, "foo", "bar", "baz"))
		expect(t, query.AsSql(), "(`foo` = ?) AND (`bar` = ?) AND (`baz` < ?)")
		expect(t, query.Bind(), []interface{}{1, 2, 3})
	}()
	func() {
		query := SqlAnd(NewOrderedHash(Hash{
			"foo": 1,
			"bar": SqlEq(2),
			"baz": SqlLt(3),
		}, "bar", "baz", "foo"))
		expect(t, query.AsSql(), "(`bar` = ?) AND (`baz` < ?) AND (`foo` = ?)")
		expect(t, query.Bind(), []interface{}{2, 3, 1})
	}()
}

func checkErr(t *testing.T, f func() *QueryMaker) {
	var query *QueryMaker
	defer func() {
		if err := recover(); err != nil {
			if query != nil {
				t.Errorf("does not return anything")
			}
			if err == nil {
				t.Errorf("error is thrown")
			}
		}
	}()
	query = f()
	t.Errorf("should not reach")
}

func TestArrayInBind(t *testing.T) {
	checkErr(t, func() *QueryMaker {
		return SqlEq("foo", Array{1, 2, 3})
	})
	checkErr(t, func() *QueryMaker {
		return SqlIn("foo", Array{Array{1, 2, 3}, 4})
	})
	checkErr(t, func() *QueryMaker {
		return SqlAnd("a", Array{Array{1, 2}, 3})
	})
}

func TestAsSqlWithQuoteCb(t *testing.T) {
	func() {
		query := SqlEq("foo.bar", "baz")
		expect(t, query.AsSql(Option{QuoteCb: func(label interface{}) string {
			tmp := []string{}
			for _, v := range strings.Split(fmt.Sprintf("%v", label), `.`) {
				tmp = append(tmp, `"`+v+`"`)
			}
			return strings.Join(tmp, `.`)
		}}), `"foo"."bar" = ?`)
		expect(t, query.Bind(), []interface{}{"baz"})
	}()
}
