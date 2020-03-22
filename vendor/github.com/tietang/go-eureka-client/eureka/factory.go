package eureka

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/tietang/go-eureka-client/eureka/config"
	"github.com/tietang/go-utils"
	"github.com/tietang/props/kvs"
)

const (
	DataCenterNameNetflix = "Netflix"
	DataCenterNameAmazon  = "Amazon"
	DataCenterNameMyOwn   = "MyOwn"
	//
	DataCenterNameMyOwnClass = "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo"

	//
	StatusUp           = "UP"             //Ready to receive traffic
	StatusDown         = "DOWN"           // Do not send traffic- healthcheck callback failed
	StatusStarting     = "STARTING"       //Just about starting- initializations to be done - do not send traffic
	StatusOutOfService = "OUT_OF_SERVICE" // Intentionally shutdown for traffic
	StatusUnknown      = "UNKNOWN"
	//
	DEFAULT_LEASE_RENEWAL_INTERVAL = 30
	DEFAULT_LEASE_DURATION         = 90
)

func CreateDataCenterInfo(dataCenterInfo *DataCenterInfo) *DataCenterInfo {
	if dataCenterInfo == nil {
		dataCenterInfo = &DataCenterInfo{}
	}
	if dataCenterInfo.Name == "" {
		dataCenterInfo.Name = DataCenterNameMyOwn
	}

	if dataCenterInfo.Class == "" {
		dataCenterInfo.Class = DataCenterNameMyOwnClass
	}
	return dataCenterInfo
}

func CreateInstanceInfo(config *config.EurekaConfig) *InstanceInfo {
	instanceConfig := config.Eureka.Instance
	appConfig := config.Application
	ins := createInstanceInfo(instanceConfig, appConfig)

	return ins

}

func defaultFileName() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	fileName := filepath.Base(path)
	return fileName
}

func createInstanceInfo(config config.EurekaInstanceConfig, appConfig config.ApplicationConfig) *InstanceInfo {
	dataCenterInfo := CreateDataCenterInfo(&DataCenterInfo{})

	leaseInfo := NewLeaseInfo(config.LeaseRenewalIntervalInSeconds)
	leaseInfo.DurationInSecs = config.LeaseExpirationDurationInSeconds

	ips, _ := utils.GetExternalIPs()
	hostName := ips[0]

	appName := defaultFileName()

	if appConfig.Name != "" {
		appName = appConfig.Name
	}

	instanceInfo := &InstanceInfo{
		HostName:       hostName,
		App:            appName,
		AppName:        appName,
		AppGroupName:   config.AppGroupName,
		IpAddr:         ips[0],
		Status:         StatusStarting,
		DataCenterInfo: dataCenterInfo,
		LeaseInfo:      leaseInfo,
		Metadata:       nil,
	}

	scheme := "http"
	portj := &Port{
		Port:    appConfig.Port,
		Enabled: true,
	}

	if appConfig.Secured || config.SecurePortEnabled {
		instanceInfo.SecurePort = portj
		instanceInfo.SecureVipAddress = instanceInfo.AppName
		scheme = "https"
	} else {
		instanceInfo.VipAddress = instanceInfo.AppName
		instanceInfo.Port = portj
	}

	stringPort := ":" + strconv.Itoa(portj.Port)

	instanceInfo.StatusPageUrl = scheme + "://" + hostName + stringPort + "/info"
	instanceInfo.HealthCheckUrl = scheme + "://" + hostName + stringPort + "/health"
	instanceInfo.HomePageUrl = scheme + "://" + hostName + stringPort + "/"

	instanceInfo.Metadata = &MetaData{
		Map: make(map[string]string),
	}
	kv := config.MetadataMap
	for k, v := range kv {
		instanceInfo.Metadata.Map[k] = v
	}
	//	instanceInfo.Metadata.Map["foo"] = "bar" //add metadata for example
	instanceId := fmt.Sprintf("%s:%s:%d", instanceInfo.HostName, instanceInfo.AppName, instanceInfo.Port.Port)
	instanceInfo.Metadata.Map["instanceId"] = instanceId
	instanceInfo.InstanceId = instanceId

	return instanceInfo
}

func CreateEurekaClient(eurekaConfig config.EurekaClientConfig) *Client {
	zones := eurekaConfig.GetAvailabilityZones(eurekaConfig.Region)
	machines := make([]string, 0)
	for _, zone := range zones {
		machinesForZone := eurekaConfig.GetEurekaServerServiceUrls(zone)
		if len(machinesForZone) > 0 {
			machines = append(machines, machinesForZone...)
		}
	}
	c := Config{
		// default timeout is one second
		DialTimeout: time.Second,
	}
	client := NewClientByConfig(machines, c)
	client.ClientConfig = &eurekaConfig
	return client
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func CreateEurekaClientByYaml(fileName string) *Client {
	file, _ := os.Getwd()
	fmt.Println("current path:", file)
	if fileName == "" {
		fileName = "application.yml"
	}

	cfg := file + "/" + fileName //"/application.yml"

	c := &config.EurekaConfig{
		Eureka: config.Eureka{
			Client:   config.NewEurekaClientConfig(),
			Instance: config.NewEurekaInstanceConfig(),
		},
	}

	if exists(cfg) {

		data, err := ReadFile(cfg)
		err = yaml.Unmarshal([]byte(data), c)
		if err != nil {
			fmt.Println("error: %v", err)
		}
	} else {
		fmt.Println("error: file %s not exists.", cfg)
	}

	client := newClientbyConfig(c)
	return client
}

func newClientbyConfig(c *config.EurekaConfig) *Client {
	ins := CreateInstanceInfo(c)
	client := CreateEurekaClient(c.Eureka.Client)
	client.InstanceInfo = ins
	client.ClientConfig = &c.Eureka.Client
	client.InstanceConfig = &c.Eureka.Instance
	client.InstanceConfig.Appname = c.Application.Name
	return client
}

// NewClient create a basic client that is configured to be used
// with the given machine list.
func NewClient(conf kvs.ConfigSource) *Client {
	c := &config.EurekaConfig{
		Eureka: config.Eureka{
			Client:   config.NewEurekaClientConfig(),
			Instance: config.NewEurekaInstanceConfig(),
		},
	}

	err := conf.Unmarshal(c)
	if err != nil {
		panic(err)
	}

	//machines := []string{}
	//config := Config{
	//    // default timeout is one second
	//    DialTimeout: time.Second,
	//}
	//
	//client := &Client{
	//    Cluster: NewCluster(machines),
	//    Config:  config,
	//}
	//
	//client.initHTTPClient()
	//
	client := newClientbyConfig(c)
	return client
}

func NewClientByConfig(machines []string, config Config) *Client {
	client := &Client{
		Cluster: NewCluster(machines),
		Config:  config,
	}
	logger.Debugf("%+v", client.Cluster)
	client.initHTTPClient()
	return client
}

func NewClientDefault(machines []string) *Client {
	config := Config{
		// default timeout is one second
		DialTimeout: time.Second,
	}
	client := &Client{
		Cluster: NewCluster(machines),
		Config:  config,
	}
	logger.Debugf("%+v", client.Cluster)
	client.initHTTPClient()
	return client
}

// NewTLSClient create a basic client with TLS configuration
func NewTLSClient(machines []string, cert string, key string, caCerts []string) (*Client, error) {
	// overwrite the default machine to use https
	if len(machines) == 0 {
		machines = []string{"https://127.0.0.1:8761"}
	}

	config := Config{
		// default timeout is one second
		DialTimeout: time.Second,
		CertFile:    cert,
		KeyFile:     key,
		CaCertFile:  make([]string, 0),
	}

	client := &Client{
		Cluster: NewCluster(machines),
		Config:  config,
	}

	err := client.initHTTPSClient(cert, key)
	if err != nil {
		return nil, err
	}

	for _, caCert := range caCerts {
		if err := client.AddRootCA(caCert); err != nil {
			return nil, err
		}
	}
	return client, nil
}

// NewClientFromFile creates a client from a given file path.
// The given file is expected to use the JSON format.
func NewClientFromFile(fpath string) (*Client, error) {
	fi, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	return NewClientFromReader(fi)
}

// NewClientFromReader creates a Client configured from a given reader.
// The configuration is expected to use the JSON format.
func NewClientFromReader(reader io.Reader) (*Client, error) {
	c := new(Client)

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	if c.Config.CertFile == "" {
		c.initHTTPClient()
	} else {
		err = c.initHTTPSClient(c.Config.CertFile, c.Config.KeyFile)
	}

	if err != nil {
		return nil, err
	}

	for _, caCert := range c.Config.CaCertFile {
		if err := c.AddRootCA(caCert); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func ReadFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, err := ioutil.ReadAll(f)
	return src, err
}
