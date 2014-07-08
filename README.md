# go-sql-querymaker

## SYNOPSYS

```
query := SqlEq("foo", V)
query.AsSql() // `foo`=?
query.Bind()  // (V)
```

```
query := SqlLg("foo", V)
query.AsSql() // `foo`<?
query.Bind()  // (V)
```

```
query := SqlIn("foo", Array{V1, V2, V3})
query.AsSql() // `foo` IN (?,?,?)
query.Bind()  // (V1,V2,V3)
```

```
query := SqlAnd("foo", Array{
	SqlGe(Min),
	SqlLt(Max),
})
query.AsSql() // `foo`>=? AND `foo`<?
query.Bind()  // (Min,Max)
```

```
query := SqlAnd("foo", Array{
	SqlEq("foo", V1),
	SqlEq("bar", V2),
})
query.AsSql() // `foo`=? AND `bar`=?
query.Bind()  // (V1,V2)
```

```
query := SqlAnd("foo", Hash{
	foo: V1,
	bar: SqlLt(V2),
})
query.AsSql() // `foo`>=? AND `bar`<?
query.Bind()  // (V1,V2)
```

```
query := SqlAnd("foo", NewOrderedHash(Hash{
	foo: V1,
	bar: SqlLt(V2),
}, "foo", "bar")
query.AsSql() // `foo`>=? AND `bar`<?
query.Bind()  // (V1,V2)
```

## SEE ALSO

Perl version is located at https://github.com/tokuhirom/SQL-Maker
