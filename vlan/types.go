package vlan
import (
"sync"
"net"
"github.com/codegangsta/cli"
)


type endpoint struct {
	id      string
	mac     net.HardwareAddr
	addr    *net.IPNet
	srcName string
}

type endpointTable map[string]*endpoint

type network struct {
	id        string
	vlanId    int
	endpoints endpointTable
	gateway   string
	ifaceOpt  string
	modeOpt   string
	sync.Mutex
	cidr *net.IPNet
}


type KvPath struct {
	base string
}

type  networkTable   map[string]*network


type Context struct {
	clusterStorage string
	parentEth string
}

