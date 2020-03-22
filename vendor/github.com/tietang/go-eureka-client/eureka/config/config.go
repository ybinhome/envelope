package config

import (
	"strings"
)

const (
	EurekaAcceptFull    = "full"
	EurekaAcceptCompact = "compact"
	DEFAULT_EUREKA_URL  = "http://localhost:8761/eureka/"
	DEFAULT_EUREKA_ZONE = "defaultZone"
	UP                  = "UP"
	DOWN                = "DOWN"
	STARTING            = "STARTING"
)

type ApplicationConfig struct {
	Name    string `yam:"name" json:"name" xml:"name"`
	Port    int    `yam:"port" json:"port" xml:"port"`
	Secured bool   `yam:"secured" json:"secured" xml:"secured"`
}
type EurekaConfig struct {
	_prefix     string            `prefix:""`
	Eureka      Eureka            `yam:"eureka" json:"eureka" xml:"eureka"`
	Application ApplicationConfig `yam:"application" json:"application" xml:"application"`
}
type Eureka struct {
	//_prefix  string `prefix:"eureka"`
	Client   EurekaClientConfig   `yam:"client" json:"client" xml:"client"`
	Instance EurekaInstanceConfig `yam:"instance" json:"instance" xml:"instance"`
}
type EurekaTransportConfig struct {
	SessionedClientReconnectIntervalSeconds           int     `yam:"sessionedClientReconnectIntervalSeconds" json:"sessionedClientReconnectIntervalSeconds" xml:"sessionedClientReconnectIntervalSeconds"`
	RetryableClientQuarantineRefreshPercentage        float64 `yam:"retryableClientQuarantineRefreshPercentage" json:"retryableClientQuarantineRefreshPercentage" xml:"retryableClientQuarantineRefreshPercentage"`
	BootstrapResolverRefreshIntervalSeconds           int     `yam:"bootstrapResolverRefreshIntervalSeconds" json:"bootstrapResolverRefreshIntervalSeconds" xml:"bootstrapResolverRefreshIntervalSeconds"`
	ApplicationsResolverDataStalenessThresholdSeconds int     `yam:"applicationsResolverDataStalenessThresholdSeconds" json:"applicationsResolverDataStalenessThresholdSeconds" xml:"applicationsResolverDataStalenessThresholdSeconds"`
	AsyncResolverRefreshIntervalMs                    int     `yam:"asyncResolverRefreshIntervalMs" json:"asyncResolverRefreshIntervalMs" xml:"asyncResolverRefreshIntervalMs"`
	AsyncResolverWarmUpTimeoutMs                      int     `yam:"asyncResolverWarmUpTimeoutMs" json:"asyncResolverWarmUpTimeoutMs" xml:"asyncResolverWarmUpTimeoutMs"`
	AsyncExecutorThreadPoolSize                       int     `yam:"asyncExecutorThreadPoolSize" json:"asyncExecutorThreadPoolSize" xml:"asyncExecutorThreadPoolSize"`
	ReadClusterVip                                    string  `yam:"readClusterVip" json:"readClusterVip" xml:"readClusterVip"`
	BootstrapResolverForQuery                         bool    `yam:"bootstrapResolverForQuery" json:"bootstrapResolverForQuery" xml:"bootstrapResolverForQuery"`
}

func NewEurekaTransportConfig() EurekaTransportConfig {
	return EurekaTransportConfig{
		SessionedClientReconnectIntervalSeconds:           1200,
		RetryableClientQuarantineRefreshPercentage:        0.66,
		BootstrapResolverRefreshIntervalSeconds:           300,
		ApplicationsResolverDataStalenessThresholdSeconds: 300,
		AsyncResolverRefreshIntervalMs:                    300000,
		AsyncResolverWarmUpTimeoutMs:                      5000,
		AsyncExecutorThreadPoolSize:                       5,
		BootstrapResolverForQuery:                         true,
	}
}

type EurekaInstanceConfig struct {
	Appname                          string `yam:"appname" json:"appname" xml:"appname"`
	AppGroupName                     string `yam:"appGroupName" json:"appGroupName" xml:"appGroupName"`
	InstanceEnabledOnit              bool   `yam:"instanceEnabledOnit" json:"instanceEnabledOnit" xml:"instanceEnabledOnit"`
	NonSecurePort                    int    `yam:"nonSecurePort" json:"nonSecurePort" xml:"nonSecurePort"`
	SecurePort                       int    `yam:"securePort" json:"securePort" xml:"securePort"`
	NonSecurePortEnabled             bool   `yam:"nonSecurePortEnabled" json:"nonSecurePortEnabled" xml:"nonSecurePortEnabled"`
	SecurePortEnabled                bool   `yam:"securePortEnabled" json:"securePortEnabled" xml:"securePortEnabled"`
	LeaseRenewalIntervalInSeconds    uint   `yam:"leaseRenewalIntervalInSeconds" json:"leaseRenewalIntervalInSeconds" xml:"leaseRenewalIntervalInSeconds"`
	LeaseExpirationDurationInSeconds int    `yam:"leaseExpirationDurationInSeconds" json:"leaseExpirationDurationInSeconds" xml:"leaseExpirationDurationInSeconds"`
	VirtualHostName                  string `yam:"virtualHostName" json:"virtualHostName" xml:"virtualHostName"`
	InstanceId                       string `yam:"instanceId" json:"instanceId" xml:"instanceId"`
	SecureVirtualHostName            string `yam:"secureVirtualHostName" json:"secureVirtualHostName" xml:"secureVirtualHostName"`
	ASGName                          string `yam:"aSGName" json:"aSGName" xml:"aSGName"`
	//	DataCenterInfo                   DataCenterInfoConfig `yam:"dataCenterInfo" json:"dataCenterInfo" xml:"dataCenterInfo"`
	IpAddress            string            `yam:"ipAddress" json:"ipAddress" xml:"ipAddress"`
	StatusPageUrlPath    string            `yam:"statusPageUrlPath" json:"statusPageUrlPath" xml:"statusPageUrlPath"`
	StatusPageUrl        string            `yam:"statusPageUrl" json:"statusPageUrl" xml:"statusPageUrl"`
	HomePageUrlPath      string            `yam:"homePageUrlPath" json:"homePageUrlPath" xml:"homePageUrlPath"`
	HomePageUrl          string            `yam:"homePageUrl" json:"homePageUrl" xml:"homePageUrl"`
	HealthCheckUrlPath   string            `yam:"healthCheckUrlPath" json:"healthCheckUrlPath" xml:"healthCheckUrlPath"`
	HealthCheckUrl       string            `yam:"healthCheckUrl" json:"healthCheckUrl" xml:"healthCheckUrl"`
	SecureHealthCheckUrl string            `yam:"secureHealthCheckUrl" json:"secureHealthCheckUrl" xml:"secureHealthCheckUrl"`
	Namespace            string            `yam:"namespace" json:"namespace" xml:"namespace"`
	Hostname             string            `yam:"hostname" json:"hostname" xml:"hostname"`
	PreferIpAddress      bool              `yam:"preferIpAddress" json:"preferIpAddress" xml:"preferIpAddress"`
	InitialStatus        string            `yam:"initialStatus" json:"initialStatus" xml:"initialStatus"`
	MetadataMap          map[string]string `yam:"metadataMap" json:"metadataMap" xml:"metadataMap"`
}

func NewEurekaInstanceConfig() EurekaInstanceConfig {
	ins := EurekaInstanceConfig{
		Appname:                          "unknow",
		NonSecurePort:                    8080,
		SecurePort:                       443,
		NonSecurePortEnabled:             true,
		LeaseRenewalIntervalInSeconds:    30,
		LeaseExpirationDurationInSeconds: 90,
		StatusPageUrlPath:                "/info",
		HomePageUrlPath:                  "/",
		HealthCheckUrlPath:               "/health",
		Namespace:                        "eureka",
		PreferIpAddress:                  true,
		InitialStatus:                    UP,
	}
	ins.VirtualHostName = ins.Appname

	return ins
}

type EurekaClientConfig struct {
	Transport                                     EurekaTransportConfig `yam:"transport" json:"transport" xml:"transport"`
	RegistryFetchIntervalSeconds                  int                   `yam:"registryFetchIntervalSeconds" json:"registryFetchIntervalSeconds" xml:"registryFetchIntervalSeconds"`
	InstanceInfoReplicationIntervalSeconds        int                   `yam:"instanceInfoReplicationIntervalSeconds" json:"instanceInfoReplicationIntervalSeconds" xml:"instanceInfoReplicationIntervalSeconds"`
	InitialInstanceInfoReplicationIntervalSeconds int                   `yam:"initialInstanceInfoReplicationIntervalSeconds" json:"initialInstanceInfoReplicationIntervalSeconds" xml:"initialInstanceInfoReplicationIntervalSeconds"`
	EurekaServiceUrlPollIntervalSeconds           int                   `yam:"eurekaServiceUrlPollIntervalSeconds" json:"eurekaServiceUrlPollIntervalSeconds" xml:"eurekaServiceUrlPollIntervalSeconds"`
	EurekaServerReadTimeoutSeconds                int                   `yam:"eurekaServerReadTimeoutSeconds" json:"eurekaServerReadTimeoutSeconds" xml:"eurekaServerReadTimeoutSeconds"`
	EurekaServerConnectTimeoutSeconds             int                   `yam:"eurekaServerConnectTimeoutSeconds" json:"eurekaServerConnectTimeoutSeconds" xml:"eurekaServerConnectTimeoutSeconds"`
	BackupRegistryImpl                            string                `yam:"backupRegistryImpl" json:"backupRegistryImpl" xml:"backupRegistryImpl"`
	EurekaServerTotalConnections                  int                   `yam:"eurekaServerTotalConnections" json:"eurekaServerTotalConnections" xml:"eurekaServerTotalConnections"`
	EurekaServerTotalConnectionsPerHost           int                   `yam:"eurekaServerTotalConnectionsPerHost" json:"eurekaServerTotalConnectionsPerHost" xml:"eurekaServerTotalConnectionsPerHost"`
	EurekaServerURLContext                        string                `yam:"eurekaServerURLContext" json:"eurekaServerURLContext" xml:"eurekaServerURLContext"`
	EurekaServerPort                              string                `yam:"eurekaServerPort" json:"eurekaServerPort" xml:"eurekaServerPort"`
	EurekaServerDNSName                           string                `yam:"eurekaServerDNSName" json:"eurekaServerDNSName" xml:"eurekaServerDNSName"`
	Region                                        string                `yam:"region" json:"region" xml:"region"`
	EurekaConnectionIdleTimeoutSeconds            int                   `yam:"eurekaConnectionIdleTimeoutSeconds" json:"eurekaConnectionIdleTimeoutSeconds" xml:"eurekaConnectionIdleTimeoutSeconds"`
	RegistryRefreshSingleVipAddress               string                `yam:"registryRefreshSingleVipAddress" json:"registryRefreshSingleVipAddress" xml:"registryRefreshSingleVipAddress"`
	HeartbeatExecutorThreadPoolSize               int                   `yam:"heartbeatExecutorThreadPoolSize" json:"heartbeatExecutorThreadPoolSize" xml:"heartbeatExecutorThreadPoolSize"`
	HeartbeatExecutorExponentialBackOffBound      int                   `yam:"heartbeatExecutorExponentialBackOffBound" json:"heartbeatExecutorExponentialBackOffBound" xml:"heartbeatExecutorExponentialBackOffBound"`
	CacheRefreshExecutorThreadPoolSize            int                   `yam:"cacheRefreshExecutorThreadPoolSize" json:"cacheRefreshExecutorThreadPoolSize" xml:"cacheRefreshExecutorThreadPoolSize"`
	CacheRefreshExecutorExponentialBackOffBound   int                   `yam:"cacheRefreshExecutorExponentialBackOffBound" json:"cacheRefreshExecutorExponentialBackOffBound" xml:"cacheRefreshExecutorExponentialBackOffBound"`
	GZipContent                                   bool                  `yam:"gZipContent" json:"gZipContent" xml:"gZipContent"`
	UseDnsForFetchingServiceUrls                  bool                  `yam:"useDnsForFetchingServiceUrls" json:"useDnsForFetchingServiceUrls" xml:"useDnsForFetchingServiceUrls"`
	RegisterWithEureka                            bool                  `yam:"registerWithEureka" json:"registerWithEureka" xml:"registerWithEureka"`
	PreferSameZoneEureka                          bool                  `yam:"preferSameZoneEureka" json:"preferSameZoneEureka" xml:"preferSameZoneEureka"`
	LogDeltaDiff                                  bool                  `yam:"logDeltaDiff" json:"logDeltaDiff" xml:"logDeltaDiff"`
	DisableDelta                                  bool                  `yam:"disableDelta" json:"disableDelta" xml:"disableDelta"`
	FetchRemoteRegionsRegistry                    string                `yam:"fetchRemoteRegionsRegistry" json:"fetchRemoteRegionsRegistry" xml:"fetchRemoteRegionsRegistry"`
	FilterOnlyUpInstances                         bool                  `yam:"filterOnlyUpInstances" json:"filterOnlyUpInstances" xml:"filterOnlyUpInstances"`
	FetchRegistry                                 bool                  `yam:"fetchRegistry" json:"fetchRegistry" xml:"fetchRegistry"`
	DollarReplacement                             string                `yam:"dollarReplacement" json:"dollarReplacement" xml:"dollarReplacement"`
	EscapeCharReplacement                         string                `yam:"escapeCharReplacement" json:"escapeCharReplacement" xml:"escapeCharReplacement"`
	AllowRedirects                                bool                  `yam:"allowRedirects" json:"allowRedirects" xml:"allowRedirects"`
	OnDemandUpdateStatusChange                    bool                  `yam:"onDemandUpdateStatusChange" json:"onDemandUpdateStatusChange" xml:"onDemandUpdateStatusChange"`
	EncoderName                                   string                `yam:"encoderName" json:"encoderName" xml:"encoderName"`
	DecoderName                                   string                `yam:"decoderName" json:"decoderName" xml:"decoderName"`
	ClientDataAccept                              string                `yam:"clientDataAccept" json:"clientDataAccept" xml:"clientDataAccept"`
	AvailabilityZones                             map[string]string     `yam:"availabilityZones" json:"availabilityZones" xml:"availabilityZones"`
	ServiceUrl                                    map[string]string     `yam:"serviceUrl" json:"serviceUrl" xml:"serviceUrl"`
}

func NewEurekaClientConfig() EurekaClientConfig {
	return EurekaClientConfig{
		Transport:                                     NewEurekaTransportConfig(),
		RegistryFetchIntervalSeconds:                  30,
		InstanceInfoReplicationIntervalSeconds:        30,
		InitialInstanceInfoReplicationIntervalSeconds: 40,
		EurekaServiceUrlPollIntervalSeconds:           300,
		EurekaServerReadTimeoutSeconds:                8,
		EurekaServerConnectTimeoutSeconds:             5,
		EurekaServerTotalConnections:                  200,
		EurekaServerTotalConnectionsPerHost:           50,
		Region:                                        "cn-east-1",
		EurekaConnectionIdleTimeoutSeconds:            30,
		HeartbeatExecutorThreadPoolSize:               2,
		HeartbeatExecutorExponentialBackOffBound:      10,
		CacheRefreshExecutorThreadPoolSize:            2,
		CacheRefreshExecutorExponentialBackOffBound:   10,
		ServiceUrl: map[string]string{DEFAULT_EUREKA_ZONE: DEFAULT_EUREKA_URL},

		GZipContent:                  true,
		UseDnsForFetchingServiceUrls: false,
		RegisterWithEureka:           true,
		PreferSameZoneEureka:         true,
		AvailabilityZones:            make(map[string]string),
		FilterOnlyUpInstances:        true,
		FetchRegistry:                true,
		DollarReplacement:            "_-",
		EscapeCharReplacement:        "__",
		AllowRedirects:               false,
		OnDemandUpdateStatusChange:   true,
		ClientDataAccept:             EurekaAcceptFull,
	}
}

func (e *EurekaClientConfig) GetAvailabilityZones(region string) []string {
	value, ok := e.AvailabilityZones[region]

	if !ok {
		value = DEFAULT_EUREKA_ZONE
	}

	return strings.Split(value, ",")
}

func (e *EurekaClientConfig) GetEurekaServerServiceUrls(myZone string) []string {
	serviceUrls, ok := e.ServiceUrl[myZone]
	if !ok {
		serviceUrls = e.ServiceUrl[DEFAULT_EUREKA_ZONE]
	}
	if serviceUrls == "" {
		return []string{"https://127.0.0.1:8761"}
	}
	return strings.Split(serviceUrls, ",")
}
