package querymaker

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mix3/go-sql-querymaker/util"
)

var matchCc = regexp.MustCompile(`@`)
var matchPh = regexp.MustCompile(`\?`)

type Array []interface{}
type quoteCb func(interface{}) string
type asSql func(column interface{}, quoteCb quoteCb) string
type builder func(string) string

type Option struct {
	Column  interface{}
	QuoteCb quoteCb
}

type QueryMaker struct {
	column interface{}
	asSql  asSql
	bind   []interface{}
}

func newQueryMaker(column interface{}, asSql asSql, bind ...interface{}) *QueryMaker {
	for _, v := range bind {
		switch v.(type) {
		case IHash, Array:
			panic(fmt.Sprintf("cannot bind an Array or an IHash"))
		default:
		}
	}
	return &QueryMaker{
		column: column,
		asSql:  asSql,
		bind:   bind,
	}
}

var fnop = map[string]string{
	`SqlAnd`:        `AND`,
	`SqlOr`:         `OR`,
	`SqlIn`:         `IN`,
	`SqlNotIn`:      `NOT IN`,
	`SqlIsNull`:     `IS NULL`,
	`SqlIsNotNull`:  `IS NOT NULL`,
	`SqlEq`:         `= ?`,
	`SqlNe`:         `!= ?`,
	`SqlLt`:         `< ?`,
	`SqlGt`:         `> ?`,
	`SqlLe`:         `<= ?`,
	`SqlGe`:         `>= ?`,
	`SqlLike`:       `LIKE ?`,
	`SqlBetween`:    `BETWEEN ? AND ?`,
	`SqlNotBetween`: `NOT BETWEEN ? AND ?`,
	`SqlNot`:        `NOT @`,
}

func andOrOp(fn string, args ...interface{}) *QueryMaker {
	op := fnop[fn]
	var opArgs Array
	var column, tmpOpArgs interface{}
	tmpOpArgs, args = util.Pop(args)
	column, args = util.Shift(args)
	switch t := tmpOpArgs.(type) {
	default:
		panic(fmt.Sprintf("arguments to `%s` must be contained in an Array or a IHash", op))
	case IHash:
		if column != nil {
			panic(fmt.Sprintf("cannot specify the column name as another argument when the conditions are listed using IHash"))
		}
		conds := Array{}
		for _, col := range t.Keys() {
			val := t.Value(col)
			switch u := val.(type) {
			default:
				val = SqlEq(col, val)
			case *QueryMaker:
				u.BindColumn(col)
			}
			conds = util.Push(conds, val)
		}
		opArgs = conds
	case Array:
		opArgs = t
	}
	return newQueryMaker(column, func(column interface{}, quoteCb quoteCb) string {
		if len(opArgs) == 0 {
			if op == "AND" {
				return "0=1"
			} else {
				return "1=1"
			}
		}
		terms := []string{}
		for _, arg := range opArgs {
			switch t := arg.(type) {
			default:
				if column == nil {
					panic(fmt.Sprintf("no column binding for %s", fn))
				}
				terms = append(terms, "("+quoteCb(column)+" = ?)")
			case *QueryMaker:
				term := t.AsSql(Option{
					Column:  column,
					QuoteCb: quoteCb,
				})
				terms = append(terms, "("+term+")")
			}
		}
		return strings.Join(terms, " "+op+" ")
	}, func() Array {
		bind := Array{}
		for _, arg := range opArgs {
			switch t := arg.(type) {
			default:
				bind = append(bind, t)
			case *QueryMaker:
				bind = append(bind, t.Bind()...)
			}
		}
		return bind
	}()...)
}

func inNotInOp(fn string, args ...interface{}) *QueryMaker {
	op := fnop[fn]
	var opArgs Array
	var column, tmpOpArgs interface{}
	tmpOpArgs, args = util.Pop(args)
	column, args = util.Shift(args)
	switch t := tmpOpArgs.(type) {
	default:
		panic(fmt.Sprintf("arguments to `%s` must be contained in Array", op))
	case Array:
		opArgs = t
	}
	return newQueryMaker(column, func(column interface{}, quoteCb quoteCb) string {
		if column == nil {
			panic(fmt.Sprintf("no column binding for %s", fn))
		}
		if len(opArgs) == 0 {
			if op == "IN" {
				return "0=1"
			} else {
				return "1=1"
			}
		}
		terms := []string{}
		for _, arg := range opArgs {
			switch t := arg.(type) {
			default:
				terms = append(terms, "?")
			case *QueryMaker:
				term := t.AsSql(Option{
					Column:  nil,
					QuoteCb: quoteCb,
				})
				if term == "?" {
					terms = append(terms, term)
				} else {
					terms = append(terms, "("+term+")")
				}
			}
		}
		return quoteCb(column) + " " + op + " (" + strings.Join(terms, ",") + ")"
	}, func() Array {
		bind := Array{}
		for _, arg := range opArgs {
			switch t := arg.(type) {
			default:
				bind = append(bind, t)
			case *QueryMaker:
				bind = append(bind, t.Bind()...)
			}
		}
		return bind
	}()...)
}

func SqlAnd(args ...interface{}) *QueryMaker        { return andOrOp("SqlAnd", args...) }
func SqlOr(args ...interface{}) *QueryMaker         { return andOrOp("SqlOr", args...) }
func SqlIn(args ...interface{}) *QueryMaker         { return inNotInOp("SqlIn", args...) }
func SqlNotIn(args ...interface{}) *QueryMaker      { return inNotInOp("SqlNotIn", args...) }
func SqlIsNull(args ...interface{}) *QueryMaker     { return fnOp("SqlIsNull", args...) }
func SqlIsNotNull(args ...interface{}) *QueryMaker  { return fnOp("SqlIsNotNull", args...) }
func SqlEq(args ...interface{}) *QueryMaker         { return fnOp("SqlEq", args...) }
func SqlNe(args ...interface{}) *QueryMaker         { return fnOp("SqlNe", args...) }
func SqlLt(args ...interface{}) *QueryMaker         { return fnOp("SqlLt", args...) }
func SqlGt(args ...interface{}) *QueryMaker         { return fnOp("SqlGt", args...) }
func SqlLe(args ...interface{}) *QueryMaker         { return fnOp("SqlLe", args...) }
func SqlGe(args ...interface{}) *QueryMaker         { return fnOp("SqlGe", args...) }
func SqlLike(args ...interface{}) *QueryMaker       { return fnOp("SqlLike", args...) }
func SqlBetween(args ...interface{}) *QueryMaker    { return fnOp("SqlBetween", args...) }
func SqlNotBetween(args ...interface{}) *QueryMaker { return fnOp("SqlNotBetween", args...) }
func SqlNot(args ...interface{}) *QueryMaker        { return fnOp("SqlNot", args...) }

func fnOp(fn string, args ...interface{}) *QueryMaker {
	numArgs, builder := compileBuilder(fnop[fn])
	var column interface{}
	if numArgs < len(args) {
		column, args = util.Shift(args)
	}
	if numArgs != len(args) {
		panic(fmt.Sprintf("the operator expects %d parameters, but got %d", numArgs, len(args)))
	}
	return sqlOp(fn, builder, column, args...)
}

func SqlOp(args ...interface{}) *QueryMaker {
	var tmpOpArgs, expr, column interface{}
	tmpOpArgs, args = util.Pop(args)
	expr, args = util.Pop(args)
	opArgs := tmpOpArgs.(Array)
	numArgs, builder := compileBuilder(fmt.Sprintf("%v", expr))
	if numArgs != len(args) {
		panic(fmt.Sprintf("the operator expects %d parameters, but got %d", numArgs, len(args)))
	}
	column, args = util.Shift(args)
	return sqlOp("SqlOp", builder, column, opArgs...)
}

func sqlOp(fn string, builder builder, column interface{}, args ...interface{}) *QueryMaker {
	return newQueryMaker(column, func(column interface{}, quoteCb quoteCb) string {
		if column == nil {
			panic(fmt.Sprintf(`no column binding for %s(args...)`, fn))
		}
		return builder(quoteCb(column))
	}, args...)
}

func SqlRaw(sql string, bind ...interface{}) *QueryMaker {
	return newQueryMaker("", func(column interface{}, quoteCb quoteCb) string {
		return sql
	}, bind...)
}

func compileBuilder(expr string) (int, builder) {
	if !matchCc.MatchString(expr) {
		expr = `@ ` + expr
	}
	numArgs := len(matchPh.FindAllString(expr, -1))
	exprs := strings.Split(expr, `@`)
	var builder builder = func(arg string) string {
		return strings.Join(exprs, arg)
	}
	return numArgs, builder
}

func (qm *QueryMaker) BindColumn(column interface{}) {
	if column != nil && qm.column != nil {
		panic(fmt.Sprintf("cannot rebind column for `%v` to: `%v`", qm.column, column))
	}
	qm.column = column
}

func (qm *QueryMaker) AsSql(option ...Option) string {
	quoteCb := quoteIdentifier
	if 0 < len(option) {
		if option[0].Column != nil {
			qm.BindColumn(option[0].Column)
		}
		if option[0].QuoteCb != nil {
			quoteCb = option[0].QuoteCb
		}
	}
	return qm.asSql(qm.column, quoteCb)
}

func (qm *QueryMaker) Bind() []interface{} {
	return qm.bind
}

var quoteIdentifier = func(label interface{}) string {
	tmp := []string{}
	for _, v := range strings.Split(fmt.Sprintf("%v", label), `\.`) {
		tmp = append(tmp, "`"+v+"`")
	}
	return strings.Join(tmp, `\.`)
}
