package userManage

import "fmt"

func (umg *UserManage) Ulist() {
	// ViewLoggedUser 查看登录的用户和对应的地址。
	ct := 0
	for _, v := range umg.EffectiveUser {
		ct++
		fmt.Printf("%d %s %s\n", ct, v.UserName, v.Address)
	}
	fmt.Println("-end")
}
