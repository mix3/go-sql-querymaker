package querymaker

import (
	"reflect"
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
		query := Eq("foo", "bar")
		expect(t, query.AsSql(), "`foo` = ?")
		expect(t, query.Bind(), []interface{}{"bar"})
	}()
	func() {
		query := In("foo", []interface{}{"bar", "baz"})
		expect(t, query.AsSql(), "`foo` IN (?,?)")
		expect(t, query.Bind(), []interface{}{"bar", "baz"})
	}()
	func() {
		query := And([]interface{}{Eq("foo", "bar"), Eq("baz", 123)})
		expect(t, query.AsSql(), "(`foo` = ?) AND (`baz` = ?)")
		expect(t, query.Bind(), []interface{}{"bar", 123})
	}()
	func() {
		query := And("foo", []interface{}{Ge(3), Lt(5)})
		expect(t, query.AsSql(), "(`foo` >= ?) AND (`foo` < ?)")
		expect(t, query.Bind(), []interface{}{3, 5})
	}()
	func() {
		query := Or([]interface{}{Eq("foo", "bar"), Eq("baz", 123)})
		expect(t, query.AsSql(), "(`foo` = ?) OR (`baz` = ?)")
		expect(t, query.Bind(), []interface{}{"bar", 123})
	}()
	func() {
		query := Or("foo", []interface{}{"bar", "baz"})
		expect(t, query.AsSql(), "(`foo` = ?) OR (`foo` = ?)")
		expect(t, query.Bind(), []interface{}{"bar", "baz"})
	}()
	func() {
		query := IsNull("foo")
		expect(t, query.AsSql(), "`foo` IS NULL")
		expect(t, query.Bind(), []interface{}{})
	}()
	func() {
		query := IsNotNull("foo")
		expect(t, query.AsSql(), "`foo` IS NOT NULL")
		expect(t, query.Bind(), []interface{}{})
	}()
	func() {
		query := Between("foo", 1, 2)
		expect(t, query.AsSql(), "`foo` BETWEEN ? AND ?")
		expect(t, query.Bind(), []interface{}{1, 2})
	}()
	func() {
		query := Not("foo")
		expect(t, query.AsSql(), "NOT `foo`")
		expect(t, query.Bind(), []interface{}{})
	}()
	func() {
		query := SqlOp("apples", "MATCH (@) AGAINST (?)", []interface{}{"oranges"})
		expect(t, query.AsSql(), "MATCH (`apples`) AGAINST (?)")
		expect(t, query.Bind(), []interface{}{"oranges"})
	}()
	func() {
		query := Raw("SELECT * FROM t WHERE id=?", 123)
		expect(t, query.AsSql(), "SELECT * FROM t WHERE id=?")
		expect(t, query.Bind(), []interface{}{123})
	}()
	func() {
		query := In("foo", []interface{}{123, Raw("SELECT id FROM t WHERE cat=?", 5)})
		expect(t, query.AsSql(), "`foo` IN (?,(SELECT id FROM t WHERE cat=?))")
		expect(t, query.Bind(), []interface{}{123, 5})
	}()
}
