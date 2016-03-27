package main
import (
	"github.com/uniseraph/vlan-docker-plugin/vlan"
	"github.com/docker/go-plugins-helpers/network"
	"github.com/codegangsta/cli"
	log "github.com/Sirupsen/logrus"

	"os"

	"github.com/docker/libkv/store/zookeeper"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/consul"
)



const (
	version = "0.1.0"
)

func init(){
	zookeeper.Register()
	etcd.Register()
	consul.Register()
}


func main() {


	var flClusterStore = cli.StringFlag{
		Name:   "cluster-store",
		EnvVar: "ACSNP_CLUSTER_STORE",
		Usage:  "Set the cluster store",
	}
	var flParentEth = cli.StringFlag{
		Name:   "parent-eth",
		EnvVar: "ACSNP_ETH",
		Usage:  "Set the parent eth for vlan device",
	}


	app := cli.NewApp()
	app.Name = "vlan driver"
	app.Usage = "VLAN Docker Networking"
	app.Version = version
	app.Flags = []cli.Flag{
		flClusterStore,
		flParentEth,
	}

	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stderr)
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		log.SetLevel(level)
		return nil
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "start",
			Usage: "start a vlan network plugin",
			Flags: []cli.Flag{
				flClusterStore,
				flParentEth,
			},
			Action: Start,
		},

	}

	if err:=app.Run(os.Args); err!=nil{
		log.Fatal(err)
	}
}

// Run initializes the driver
func Start(ctx *cli.Context) {

	d, err := vlan.NewDriver(version, 	vlan.ParseContext(ctx))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	h := network.NewHandler(d)
	h.ServeUnix("root", "macvlan")
}
