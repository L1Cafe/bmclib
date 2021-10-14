package idrac9

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net"
	"time"

	"github.com/bmc-toolbox/bmclib/internal"
	"github.com/bmc-toolbox/bmclib/internal/helper"
)

// CurrentHTTPSCert returns the current x509 certficates configured on the BMC
// The bool value returned indicates if the BMC supports CSR generation.
// CurrentHTTPSCert implements the Configure interface.
func (i *IDrac9) CurrentHTTPSCert() ([]*x509.Certificate, bool, error) {
	dialer := &net.Dialer{
		Timeout: time.Duration(10) * time.Second,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", i.ip+":"+"443", &tls.Config{InsecureSkipVerify: true})

	if err != nil {
		return []*x509.Certificate{{}}, true, err
	}

	defer conn.Close()

	return conn.ConnectionState().PeerCertificates, true, nil
}

// Screenshot grab screen preview.
func (i *IDrac9) Screenshot() (response []byte, extension string, err error) {
	err = i.httpLogin()
	if err != nil {
		return response, extension, err
	}

	extension = "png"
	endpoint1 := "sysmgmt/2015/server/preview"
	statusCode, _, err := i.get(endpoint1, &map[string]string{})
	if err != nil || statusCode != 200 {
		return nil, "", err
	}

	endpoint2 := "capconsole/scapture0.png"
	statusCode, response, err = i.get(endpoint2, &map[string]string{})
	if err != nil || statusCode != 200 {
		return nil, "", err
	}

	return response, extension, err
}

func (i *IDrac9) queryUsers() (users map[int]User, err error) {
	endpoint := "sysmgmt/2012/server/configgroup/iDRAC.Users"

	statusCode, response, err := i.get(endpoint, &map[string]string{})
	if err != nil || statusCode != 200 {
		i.log.V(1).Error(err, "GET request failed.",
			"IP", i.ip,
			"HardwareType", i.HardwareType(),
			"endpoint", endpoint,
			"step", helper.WhosCalling(),
			"Error", internal.ErrStringOrEmpty(err),
		)
		return users, err
	}

	userData := make(idracUsers)
	err = json.Unmarshal(response, &userData)
	if err != nil {
		i.log.V(1).Error(err, "Unable to unmarshal payload.",
			"IP", i.ip,
			"HardwareType", i.HardwareType(),
			"resource", "User",
			"step", "queryUserInfo",
			"Error", internal.ErrStringOrEmpty(err),
		)
		return users, err
	}

	return userData["iDRAC.Users"], err
}

func (i *IDrac9) queryLdapRoleGroups() (ldapRoleGroups LdapRoleGroups, err error) {
	endpoint := "sysmgmt/2012/server/configgroup/iDRAC.LDAPRoleGroup"

	statusCode, response, err := i.get(endpoint, &map[string]string{})
	if err != nil || statusCode != 200 {
		i.log.V(1).Error(err, "GET request failed.",
			"IP", i.ip,
			"HardwareType", i.HardwareType(),
			"endpoint", endpoint,
			"step", helper.WhosCalling(),
			"Error", internal.ErrStringOrEmpty(err),
		)
		return ldapRoleGroups, err
	}

	idracLdapRoleGroups := make(idracLdapRoleGroups)
	err = json.Unmarshal(response, &idracLdapRoleGroups)
	if err != nil {
		i.log.V(1).Error(err, "Unable to unmarshal payload.",
			"IP", i.ip,
			"HardwareType", i.HardwareType(),
			"resource", "User",
			"step", "queryUserInfo",
			"Error", internal.ErrStringOrEmpty(err),
		)
		return ldapRoleGroups, err
	}

	return idracLdapRoleGroups["iDRAC.LDAPRoleGroup"], err
}
