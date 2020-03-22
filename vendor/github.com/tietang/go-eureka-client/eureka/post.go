package eureka

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tietang/go-utils/errs"
)

func (c *Client) RegisterInstance(appName string, instanceInfo *InstanceInfo) error {
	values := []string{"apps", appName}
	path := strings.Join(values, "/")
	instance := &Instance{
		Instance: instanceInfo,
	}
	body, err := json.Marshal(instance)
	if err != nil {
		return err
	}

	res, err := c.Post(path, body)
	if err != nil {
		return err
	}
	if res != nil {
		log.WithFields(log.Fields{
			"path":   path,
			"status": res.StatusCode,
			"body":   string(body),
		}).Info("RegisterInstance ")
	}

	if res.StatusCode != http.StatusNoContent {
		return errs.NewError("unRegisterInstance", res.StatusCode)
	}
	return err
}
