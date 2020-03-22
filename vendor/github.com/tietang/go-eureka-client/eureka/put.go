package eureka

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func (c *Client) SendHeartbeat(appId, instanceId string) error {
	values := []string{"apps", appId, instanceId}
	path := strings.Join(values, "/")
	res, err := c.Put(path, nil)
	if res != nil {
		log.WithFields(log.Fields{
			"path":   path,
			"status": res.StatusCode,
		}).Info("SendHeartbeat ")
	}
	return err
}

func (c *Client) UpdateMetadata(appId, instanceId string, metaData map[string]string) error {
	values := []string{"apps", appId, instanceId, "metadata"}
	path := strings.Join(values, "/") + "?"
	for k, v := range metaData {
		path = path + k + "=" + v + "&"
	}
	res, err := c.Put(path, nil)
	if res != nil {
		log.WithFields(log.Fields{
			"path":   path,
			"status": res.StatusCode,
		}).Info("UpdateMetadata ")
	}
	return err
}

//PUT /eureka/v2/apps/appID/instanceID/status?value=OUT_OF_SERVICE
func (c *Client) UpdateStatus(appId, instanceId string, status string) error {
	values := []string{"apps", appId, instanceId, "status"}
	path := strings.Join(values, "/") + "?value=" + status
	res, err := c.Put(path, nil)
	if res != nil {
		log.WithFields(log.Fields{
			"path":   path,
			"status": res.StatusCode,
		}).Info("UpdateStatus ")
	}
	return err
}
