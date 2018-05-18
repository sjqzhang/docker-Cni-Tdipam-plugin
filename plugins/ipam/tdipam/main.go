package main
import (
	"encoding/json"
	"fmt"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
	"os"
	"log"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/types"
)



type EtcdConfig struct {
	Etcdcluster string `json:etcd server地址`
	Nodenetwork string `json:node network`
	Alreadyusedip string `json:container AvailableIp`
	Containernetwork string `json:containernetr`
	Routes string `json:routes`
	Dns string `json:dns`
}

type IpamConfig struct {
	Ipam	  EtcdConfig
	Name          string      `json:"name"`
	CNIVersion    string      `json:"cniVersion"`
}

func (IpamS *IpamConfig) Load(bytes []byte) error {
	err := json.Unmarshal(bytes, IpamS)
	if err != nil {
		fmt.Println("error in translating,", err.Error())

	}

	return nil
	//fmt.Println("type:", reflect.TypeOf(IpamS.Ipam.Etcdcluster))
}


func main(){
	skel.PluginMain(cmdAdd, cmdDel, version.All)
}

func cmdAdd(args *skel.CmdArgs) error {
	var Config IpamConfig
	var Cli EtcdHelper

	//读取配置
	Config = IpamConfig{}
	Config.Load(args.StdinData)
	//连接etcd
	Cli = Config.etcdConn()

	//NodeRang := Cli.getKey(Config.Ipam.Nodenetwork)
	//err := IsKeyExist(NodeRang,Config.Ipam.Nodenetwork)
	//
	//if err != nil{
	//	fmt.Println(err)
	//	os.Exit(-1)
	//}

	ContainerRange := Cli.getKey(Config.Ipam.Containernetwork)
	err := IsKeyExist(ContainerRange,Config.Ipam.Containernetwork)
	if err != nil{
		fmt.Println(err)
		os.Exit(-1)
	}

	//从etcd获取容器IP范围
	var ContainerR *Range
	ContainerR = &Range{}
	ContainerR.Container(ContainerRange,Config.Ipam.Containernetwork)
	//效验IP地址范围是否正确
	ContainerR.Canonicalize()
	//从etcd获取目前已经使用掉的IP
	AlreadUsedIp := Cli.getKey(Config.Ipam.Alreadyusedip)
	//获取目前可用的IP地址
	IpList,err := Hosts((*ContainerRange)[Config.Ipam.Containernetwork+ "SubNet"])
	if err != nil{
		log.Println("IP地址范围错误")
	}
	AvailableIp := ContainerR.RangeSet(AlreadUsedIp,&IpList,Config.Ipam.Alreadyusedip)
	//将获取到的IP提交到ETCD库中(唯一性?)

	err = Cli.setKey(Config.Ipam.Alreadyusedip,AvailableIp.String(),args.ContainerID)
	if err != nil{
		log.Println("无法存入到etcd库中")
	}
	//返回cni相关
	result := &current.Result{}
	//返回cni 版本
	result.CNIVersion = Config.CNIVersion
	//返给cni ip
	IPs := &current.IPConfig{}
	IPs.Gateway = ContainerR.Gateway
	IPs.Version = "4"
	IPs.Address.IP = AvailableIp
	IPs.Address.Mask = ContainerR.Subnet.Mask
	result.IPs = append(result.IPs,IPs)
	//获取dns配置
	dnsEtcdConfig := Cli.getKey(Config.Ipam.Dns)
	result.DNS = GetDns(dnsEtcdConfig,&Config)

	//自定义容器路由规则
	routeEtcdConfig := Cli.getKey(Config.Ipam.Routes)
	GetRoute(routeEtcdConfig,&Config)
	//_, dstmask, err := net.ParseCIDR("0.0.0.0/0")
	//Routes := &types.Route{}
	//Routes.GW = nil
	//Routes.Dst.IP = dstmask.IP
	//Routes.Dst.Mask = dstmask.Mask
	//result.Routes = append(result.Routes,Routes)


	return types.PrintResult(result, Config.CNIVersion)


}

func cmdDel(args *skel.CmdArgs) error {
	var Config IpamConfig
	var Cli EtcdHelper
	//etcd KEY路劲
	//读取配置
	Config = IpamConfig{}
	Config.Load(args.StdinData)
	//连接etcd
	Cli = Config.etcdConn()
	AlreadUsedIp := Cli.getKey(Config.Ipam.Alreadyusedip)
	//根据容器ID查找IP
	Key := ContainerSearch(AlreadUsedIp,args.ContainerID)
	if len(Key) > 0 {
		Cli.delKey(Key)
	}

	return nil
}





