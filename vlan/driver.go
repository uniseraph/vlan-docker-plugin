package vlan
import (
	"sync"

	"github.com/docker/engine-api/client"
	"github.com/vishvananda/netlink"
	sdk "github.com/docker/go-plugins-helpers/network"

	"net"
	log "github.com/Sirupsen/logrus"

	"errors"
	"github.com/docker/libkv/store"
	"strings"
	"net/url"
	"github.com/docker/docker/vendor/src/github.com/docker/libkv"
	"net/http"
	"encoding/json"
)


const (
	bridgeMode           = "bridge"
	containerIfacePrefix = "eth"
	defaultMTU           = 1500
	minMTU               = 68
)



// Driver is the VLAN Driver
type Driver struct {
	sdk.Driver
	dockerclient *client.Client
	store store.Store
	path  KvPath
	parentDev string

	nameserver string
	sync.Mutex
}


func NewDriver(version string, ctx *Context) (*Driver, error) {


	u , err :=  url.Parse( ctx.clusterStorage )
	if err !=nil {
		return nil , err
	}

	s , err := libkv.NewStore(u.Scheme, strings.Split(u.Host,","),nil)
	if err != nil{
		return nil , err
	}

	dockerclient , err :=
		client.NewClient("unix:///var/run/docker.sock", nil, &http.Client{}, defaultHeaders)

	if(err !=nil){
		return nil , err
	}


	netlink.Link

	return  &Driver{
		store : s,
		dockerclient : dockerclient,
		path : strings.Trim(u.Path,"/"),

	}


}

func (d *Driver) GetCapabilities() (*sdk.CapabilitiesResponse, error){
	return &sdk.CapabilitiesResponse{ Scope: sdk.LocalScope} , nil
}

//
//func (d *Driver) addNetwork(n *network) {
//	d.Lock()
//	d.networks[n.id] = n
//	d.Unlock()
//}

func (d *Driver) CreateNetwork(r *sdk.CreateNetworkRequest) error{

	var netCidr *net.IPNet
	var netGw string
	var err error
	log.Debugf("Network Create Called: [ %+v ]", r)
	for _, v4 := range r.IPv4Data {
		netGw = v4.Gateway
		_, netCidr, err = net.ParseCIDR(v4.Pool)
		if err != nil {
			return err
		}
	}

	n := &network{
		id:        r.NetworkID,
		endpoints: endpointTable{},
		cidr:      netCidr,
		gateway:   netGw,
	}

	// Parse docker network -o opts
	for k, v := range r.Options {
		if k == "com.docker.sdk.generic" {
			if genericOpts, ok := v.(map[string]interface{}); ok {
				for key, val := range genericOpts {
					log.Debugf("Libnetwork Opts Sent: [ %s ] Value: [ %s ]", key, val)
					// Parse -o host_iface from libnetwork generic opts
					if key == "host_iface" {
						n.ifaceOpt = val.(string)
					}
				}
			}
		}
	}

    return d.updateNetwork(n)
}

func (d *Driver) updateNetwork(n network) error {
	data, err := json.Marshal(&n)
	if err != nil {
		return err
	}
	return d.store.Put(d.path.network(n.id), data, nil)
}

func (d *Driver) deleteNetwork(id string) {

}


// DeleteNetwork deletes a network
func (d *Driver) DeleteNetwork(r *sdk.DeleteNetworkRequest) error {
	log.Debugf("Delete network request: %+v", &r)
	d.deleteNetwork(r.NetworkID)
	return nil
}


func (d *Driver) CreateEndpoint(r *sdk.CreateEndpointRequest) (*sdk.CreateEndpointResponse, error){
	var err error

	if r.Interface.Address ==nil {
		return nil , errors.New("Unable to obtain an IP address from libnetwork  ipam")
	}

	rsp := & sdk.CreateEndpointResponse{
		Interface: &sdk.EndpointInterface{
			Address:  r.Interface.Address ,
			MacAddress: makeMac(net.ParseIP(r.Interface.Address)),
		},
	}

	log.Debugf("Create endpoint %s %+v", r.EndpointID, rsp)
	return rsp , err
}


func (d *Driver)  DeleteEndpoint(r *sdk.DeleteEndpointRequest) error{

	var err  error

	d.Lock()



	defer d.Unlock()


	return err

}
