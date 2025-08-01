package tool

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

func ldapCheck() {

	l, err := ldap.DialURL("ldap://121.41.22.141:389")
	if err != nil {
		fmt.Println("ldap.DialURL失败:", err.Error())
		return
	}
	_, err = l.SimpleBind(&ldap.SimpleBindRequest{
		Username: "cn=xifeng,dc=aganyunke,dc=com",
		Password: "hzm2mnoto4_",
	})

	if err != nil {
		fmt.Println("L.SimpleBind失败:", err.Error())
		return
	}

	fmt.Println("鉴权成功")
}
