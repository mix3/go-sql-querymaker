package querymaker

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type QueryMaker struct {
	column string
	asSql  asSql
	bind   []interface{}
}

func andOrOp(fn string, args ...interface{}) *QueryMaker {
	op := fnop[fn]
	opArgs := []interface{}{}
	column := ""
	if 1 <= len(args) {
		var ok bool
		opArgs, ok = args[len(args)-1].([]interface{})
		if !ok {
			log.Fatalf("arguments to `%s` must be []interface{}", fn)
		}
	}
	if 2 <= len(args) {
		column = fmt.Sprintf("%v", args[0])
	}
	return _new(column, func(column string, quoteCb quoteCb) string {
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
				if column == "" {
					log.Fatalf("no column binding for %s", fn)
				}
				terms = append(terms, "("+quoteCb(column)+" = ?)")
			case *QueryMaker:
				term := t.AsSql(Option{
					SuppliedColname: column,
					QuoteCb:         quoteCb,
				})
				terms = append(terms, "("+term+")")
			}
		}
		return strings.Join(terms, " "+op+" ")
	}, func() []interface{} {
		bind := []interface{}{}
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
	opArgs := []interface{}{}
	column := ""
	if 1 <= len(args) {
		var ok bool
		opArgs, ok = args[len(args)-1].([]interface{})
		if !ok {
			log.Fatalf("arguments to `%s` must be []interface{}", fn)
		}
	}
	if 2 <= len(args) {
		column = fmt.Sprintf("%v", args[0])
	}
	return _new(column, func(column string, quoteCb quoteCb) string {
		if column == "" {
			log.Fatalf("no column binding for %s", fn)
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
					SuppliedColname: "",
					QuoteCb:         quoteCb,
				})
				if term == "?" {
					terms = append(terms, term)
				} else {
					terms = append(terms, "("+term+")")
				}
			}
		}
		return quoteCb(column) + " " + op + " (" + strings.Join(terms, ",") + ")"
	}, func() []interface{} {
		bind := []interface{}{}
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
	return nil
}

func And(args ...interface{}) *QueryMaker   { return andOrOp("And", args...) }
func Or(args ...interface{}) *QueryMaker    { return andOrOp("Or", args...) }
func In(args ...interface{}) *QueryMaker    { return inNotInOp("In", args...) }
func NotIn(args ...interface{}) *QueryMaker { return inNotInOp("NotIn", args...) }

var fnop = map[string]string{
	`And`:        `AND`,
	`Or`:         `OR`,
	`In`:         `IN`,
	`NotIn`:      `NOT IN`,
	`IsNull`:     `IS NULL`,
	`IsNotNull`:  `IS NOT NULL`,
	`Eq`:         `= ?`,
	`Ne`:         `!= ?`,
	`Lt`:         `< ?`,
	`Gt`:         `> ?`,
	`Le`:         `<= ?`,
	`Ge`:         `>= ?`,
	`Like`:       `LIKE ?`,
	`Between`:    `BETWEEN ? AND ?`,
	`NotBetween`: `NOT BETWEEN ? AND ?`,
	`Not`:        `NOT @`,
}

func IsNull(args ...interface{}) *QueryMaker     { return fnOp("IsNull", args...) }
func IsNotNull(args ...interface{}) *QueryMaker  { return fnOp("IsNotNull", args...) }
func Eq(args ...interface{}) *QueryMaker         { return fnOp("Eq", args...) }
func Ne(args ...interface{}) *QueryMaker         { return fnOp("Ne", args...) }
func Lt(args ...interface{}) *QueryMaker         { return fnOp("Lt", args...) }
func Gt(args ...interface{}) *QueryMaker         { return fnOp("Gt", args...) }
func Le(args ...interface{}) *QueryMaker         { return fnOp("Le", args...) }
func Ge(args ...interface{}) *QueryMaker         { return fnOp("Ge", args...) }
func Like(args ...interface{}) *QueryMaker       { return fnOp("Like", args...) }
func Between(args ...interface{}) *QueryMaker    { return fnOp("Between", args...) }
func NotBetween(args ...interface{}) *QueryMaker { return fnOp("NotBetween", args...) }
func Not(args ...interface{}) *QueryMaker        { return fnOp("Not", args...) }

func fnOp(fn string, args ...interface{}) *QueryMaker {
	numArgs, builder := compileBuilder(fnop[fn])
	column := ""
	if numArgs < len(args) {
		column = fmt.Sprintf("%v", args[0])
		args = args[1:]
	}
	if numArgs != len(args) {
		log.Fatalf("the operator expects %d parameters, but got %d", numArgs, len(args))
	}
	return sqlOp(fn, builder, column, args...)
}

func SqlOp(args ...interface{}) *QueryMaker {
	opArgs := args[len(args)-1].([]interface{})
	args = args[:len(args)-1]
	expr := fmt.Sprintf("%v", args[len(args)-1])
	args = args[:len(args)-1]
	numArgs, builder := compileBuilder(expr)
	if numArgs != len(opArgs) {
		log.Fatalf("the operator expects %d parameters, but got %d", numArgs, len(args))
	}
	column := ""
	if 0 < len(args) {
		column = fmt.Sprintf("%v", args[0])
	}
	return sqlOp("SqlOp", builder, column, opArgs...)
}

func sqlOp(fn string, builder builder, column string, args ...interface{}) *QueryMaker {
	return _new(column, func(column string, quoteCb quoteCb) string {
		if column == "" {
			log.Fatalf(`no column binding for %s(args...)`, fn)
		}
		return builder(quoteCb(column))
	}, args...)
}

func Raw(sql string, bind ...interface{}) *QueryMaker {
	return _new("", func(column string, quoteCb quoteCb) string {
		return sql
	}, bind...)
}

type quoteCb func(string) string
type asSql func(column string, quoteCb quoteCb) string

func _new(column string, asSql asSql, bind ...interface{}) *QueryMaker {
	return &QueryMaker{
		column: column,
		asSql:  asSql,
		bind:   bind,
	}
}

type builder func(string) string

var matchCc = regexp.MustCompile(`@`)
var matchPh = regexp.MustCompile(`\?`)

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

func (qm *QueryMaker) BindColumn(column string) {
	if column != "" && qm.column != "" {
		log.Fatalf("cannot rebind column for `%s` to: `%s`", qm.column, column)
	}
	qm.column = column
}

type Option struct {
	SuppliedColname string
	QuoteCb         quoteCb
}

func (qm *QueryMaker) AsSql(option ...Option) string {
	quoteCb := quoteIdentifier
	if 0 < len(option) {
		if option[0].SuppliedColname != "" {
			qm.BindColumn(option[0].SuppliedColname)
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

var quoteIdentifier = func(label string) string {
	tmp := []string{}
	for _, v := range strings.Split(label, `\.`) {
		tmp = append(tmp, "`"+v+"`")
	}
	return strings.Join(tmp, `\.`)
}
