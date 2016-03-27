package vlan
import (
	"net"
	"path/filepath"
	"github.com/codegangsta/cli"

	"github.com/docker/docker/vendor/src/github.com/docker/engine-api/client"
	"net/http"
)


// Generate a mac addr
func makeMac(ip net.IP) string {
	hw := make(net.HardwareAddr, 6)
	hw[0] = 0x7a
	hw[1] = 0x42
	copy(hw[2:], ip.To4())
	return hw.String()
}




const (
	pluginPathPrefix         = "vlan/network/v1.0"
	pluginPathNetwork        = pluginPathPrefix + "/network"
	pluginPathEndpoint       = pluginPathPrefix + "/endpoint"
	pluginPathJointEndpoints = pluginPathPrefix + "/endpoints-online"
	defaultHeaders = map[string]string{"User-Agent": "engine-api-cli-1.0"}
	defaultEth = "bond0"

)


func ParseContext(c *cli.Context) (*Context  , error ){

	context :=&Context{}


	if context.clusterStorage = c.String("cluster-storage") ; context.clusterStorage==nil {

		dockerclient, err :=
		client.NewClient("unix:///var/run/docker.sock", nil, &http.Client{}, defaultHeaders)
		if err != nil {
			return nil, err
		}

		Info, err := dockerclient.Info()
		if err != nil {
			return nil, err
		}

		context.clusterStorage = Info.ClusterStore
	}

	if context.parentEth = c.String("parent-eth") ; context.parentEth ==nil{
		context.parentEth = defaultEth
	}

	return context
}


func (p KvPath) relpath(paths []string) string {
	if p.base != "" {
		paths = append([]string{p.base}, paths...)
	}
	return filepath.Join(paths...)
}

func (p KvPath) networks() string {
	return p.relpath([]string{pluginPathNetwork})
}

func (p KvPath) network(nId string) string {
	return p.relpath([]string{pluginPathNetwork, nId})
}

func (p KvPath) endpoints() string {
	return p.relpath([]string{pluginPathEndpoint})
}

func (p KvPath) endpoint(eId string) string {
	return p.relpath([]string{pluginPathEndpoint, eId})
}

func (p KvPath) notifyLink(nid string) string {
	return p.relpath([]string{pluginPathJointEndpoints, nid})
}

func (p KvPath) notifyEndpointLink(nid, eid string) string {
	return p.relpath([]string{pluginPathJointEndpoints, nid, eid})
}
