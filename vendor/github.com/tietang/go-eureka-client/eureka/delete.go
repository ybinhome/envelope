package eureka

import "strings"
import log "github.com/sirupsen/logrus"

func (c *Client) UnregisterInstance(appId, instanceId string) error {
	values := []string{"apps", appId, instanceId}
	path := strings.Join(values, "/")
	_, err := c.Delete(path)
	return err
}

func (c *Client) DeleteStatusOverride(appId, instanceId string, status string) error {
	values := []string{"apps", appId, instanceId, "status"}
	path := strings.Join(values, "/") + "?value=" + status
	res, err := c.Delete(path)
	if res != nil {
		log.WithFields(log.Fields{
			"path":   path,
			"status": res.StatusCode,
		}).Info("DeleteStatusOverride ")
	}
	return err
}
