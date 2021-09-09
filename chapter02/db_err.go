package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	user, err := findUserByMobile("18616991233")
	if err != nil {
		log.Printf("get user error: %v", err)
		os.Exit(1)
	}

	if user != nil {
		fmt.Println("hello" + user.Name)
	}

}

//errors are values
// 我们通过对error 的判断，来跟踪程序里的意外状况
// 在这里， DAO层返回的ErrNoRows 无非是对“返回结果为空”这个情况的一个说明。
// 他通过ErrNoRows 告诉我们， 返回空是因为没有符合条件的数据， SQL语句是执行成功的。
// 一个error 是否是故障， 对当前方法来说是否是故障，这个问题需要在拿到error的现场进行判断和处理。
// 这与其说是一个错误，不如说是DAO层告诉我们没有故障发生， 没有权限问题，没有SQL语法问题， 没有网络问题，一切都好， 只是没有符合要求的数据而已
// 所以在这里，不应当对ErrNoRows 进行任何形式的包装， 而应该直接返回 nil,nil 。 上层可能会认为没查到是一个错误（比如对已经付费的用户，按照业务逻辑肯定有他的账号，但是居然查不到），那上层自己去处理。 这里的ErrNoRows 并不是错误。
// 而对于DAO层报的其他错误， 比如网络、 数据库权限、SQL语法错等等， 就需要明确告知上层。

func findUserByMobile(mobile string) (user DbRecord, err error) {
	oneUser, err := db.QueryRow("mobile='" + mobile + "'")

	//如果错误不是ErrNoRows ， 就包一下抛出去。
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "Got error on query:"+db.GetLastSql())
	}

	return oneUser, nil
}
