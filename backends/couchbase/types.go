package couchbase

import dockerapi "github.com/fsouza/go-dockerclient"

type Document struct {
	ContainerId   string
	ContainerType string
	Location      ServicePort
	MetaData      map[string]string
	Created       int64
	Updated       int64
	LBConfig      LoadBalancerConfig
	Enable        string
	//	Location      ContainerLocation
}

//type ContainerLocation struct {
//	HostIp       string
//	PortsMapping map[string]ServicePort
//}

//type CBMetaData struct {
//	Data map[string]string
//}

type LoadBalancerConfig struct {
	Enable  string
	Mode    string
	Balance string
	Param   string
}

type ServicePort struct {
	HostPort          string
	HostIP            string
	ExposedPort       string
	ExposedIP         string
	PortType          string
	ContainerHostname string
	ContainerID       string
	container         *dockerapi.Container
}
