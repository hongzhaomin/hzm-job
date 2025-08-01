package tool

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/prop"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
	"strings"
)

func LoginCheck2Ldap(userName, password string) error {
	ldapProperties := ezconfig.Get[*prop.LdapProperties]()
	if !ldapProperties.Enabled() {
		return errors.New("ldap is not enabled")
	}
	conn, err := ldap.DialURL(ldapProperties.Addr)
	if err != nil {
		global.SingletonPool().Log.Error("ldap.DialURL err: %v", err)
		return err
	}
	defer conn.Close()

	var ou string
	if ldapProperties.Group != "" {
		for _, s := range strings.Split(ldapProperties.Group, ".") {
			ou += ",ou=" + s
		}
	}

	var dc string
	for _, s := range strings.Split(ldapProperties.Dc, ".") {
		dc += ",dc=" + s
	}

	username := fmt.Sprintf("%s=%s%s%s", ldapProperties.CnKey, userName, ou, dc)
	_, err = conn.SimpleBind(&ldap.SimpleBindRequest{
		//Username: "cn=xifeng,dc=aganyunke,dc=com",
		Username: username,
		Password: password,
	})

	if err != nil {
		global.SingletonPool().Log.Error("ldap.SimpleBind err: %v", err)
		return errors.New("ldap authentication failed")
	}

	return nil
}
