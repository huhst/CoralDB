package userManage

import "fmt"

// ViewLoggedUser 查看登录的用户和对应的地址。
func (umg *UserManage) ViewLoggedUser() {
	ct := 0
	for _, v := range umg.EffectiveUser {
		ct++
		fmt.Printf("%d %s %s\n", ct, v.UserName, v.Address)
	}
	fmt.Println("-end")
}
