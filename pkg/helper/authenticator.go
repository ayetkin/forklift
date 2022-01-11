package helper

import (
	"fmt"
	"forklift/pkg/config"
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
)

func LoginOSP(configuration *config.Configuration) error {
	var err error
	if err = os.Setenv("OS_AUTH_URL", configuration.Openstack.AuthUrl); err != nil {
		return err
	}
	if err = os.Setenv("OS_USER_DOMAIN_NAME", configuration.Openstack.UserDomainName); err != nil {
		return err
	}
	if err = os.Setenv("OS_USERNAME", configuration.Openstack.Username); err != nil {
		return err
	}
	if err = os.Setenv("OS_PASSWORD", configuration.Openstack.Password); err != nil {
		return err
	}
	if err = os.Setenv("OS_REGION_NAME", configuration.Openstack.RegionName); err != nil {
		return err
	}
	if err = os.Setenv("OS_INTERFACE", configuration.Openstack.Interface); err != nil {
		return err
	}
	if err = os.Setenv("OS_IDENTITY_API_VERSION", configuration.Openstack.IdentityApiVersion); err != nil {
		return err
	}
	log.Info("Openstack environment variables set")
	return nil
}

func VcenterUrl(configuration *config.Configuration) (*url.URL, error) {
	var (
		vCenterPassword string
		err             error
		url             *url.URL
	)

	if vCenterPassword, err = GetVcenterPasswordFromVault(*configuration, "test"); err != nil {
		return nil, err
	}

	vURL := fmt.Sprintf("https://%s:%s@%s/sdk", configuration.Vcenter.Username, vCenterPassword, configuration.Vcenter.Url)
	if url, err = url.Parse(vURL); err != nil {
		return nil, err
	}

	return url, err
}
