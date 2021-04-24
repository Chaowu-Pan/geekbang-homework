package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Data struct {
	Id      int
	Content string
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname")
	if err != nil {
		panic(err)
	}
}

/**
 * 在数据库操作的时候， dao 层中当遇到一个 sql.ErrNoRows 的时候，不应该 Wrap 这个 error，抛给上层。
 * 因为dao层是一个供不同上层业务代码调用的底层逻辑，应该抛出一个原始的error；
 * 应该由调用方，也就是具体的业务代码来决定是否要处理sql.ErrNoRows，还是Wrap之后继续向上抛出。
 */
func GetUser(id int) (data Data, err error) {
	var data Data
	err = db.QueryRow("select id, content from data where id = $1", id).Scan(&data.Id, &data.Content)
	return
}
