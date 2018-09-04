package main

import (
	"JsGo/JsMariaDb"
	"fmt"
)

type User struct {
	UID      int
	Role     string
	Username string
	Password string
	Header   string
	Age      int
}

//insert

func main() {
	h := JsMariaDb.NewJsDb("www.junsie.com", "junsie", "junsie", "1qaz2wsx")

	if h != nil {
		fmt.Println("connect success")
		//insert
		db := h.Db

		user := &User{8194188, "Operator", "zhongtie", "123456", "http://www.junsie.com/junsie.png", 26}
		stmt, err := db.Prepare("insert into junsie.User(UID,Role,Username,Password,Header,Age) values(?,?,?,?,?,?)")
		if err == nil {
			_, err := stmt.Exec(user.UID, user.Role, user.Username, user.Password, user.Header, user.Age)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("insert surccess")
			}
			stmt.Close()

		} else {
			fmt.Println(err.Error())
		}

		//update
		stmt, err = db.Prepare("update User set Username=? where UID=?")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		res, err := stmt.Exec("liming", 8194188)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		affect, err := res.RowsAffected()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("affect:", affect)

		//query
		rows, err := db.Query("SELECT * FROM User")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for rows.Next() {

			user := &User{}
			err = rows.Scan(&user.UID, &user.Role, &user.Username, &user.Password, &user.Header, &user.Age)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(user.UID)
			fmt.Println(user.Role)
			fmt.Println(user.Password)
			fmt.Println(user.Age)
		}

		//delete
		stmt, err = db.Prepare("delete from User where uid=?")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		res, err = stmt.Exec(8194188)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		affect, err = res.RowsAffected()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("affect:", affect)
	}
}
