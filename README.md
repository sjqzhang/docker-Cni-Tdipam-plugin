[![Linux Build Status](https://travis-ci.org/containernetworking/plugins.svg?branch=master)](https://travis-ci.org/containernetworking/plugins)
[![Windows Build Status](https://ci.appveyor.com/api/projects/status/kcuubx0chr76ev86/branch/master?svg=true)](https://ci.appveyor.com/project/cni-bot/plugins/branch/master)

# 背景
基于cni的提供的接口开发的ipam插件,主要实现容器的IP地址/网关/DNS/路由存储在etcd库中,集中化管理IP地址,如果使用kubernetes macvlan 插件和该插件配合将实现大二层IP地址统一管理的效果.

## Plugins supplied:(官方)
### Main: interface-creating
* `bridge`: Creates a bridge, adds the host and the container to it.
* `ipvlan`: Adds an [ipvlan](https://www.kernel.org/doc/Documentation/networking/ipvlan.txt) interface in the container
* `loopback`: Creates a loopback interface
* `macvlan`: Creates a new MAC address, forwards all traffic to that to the container
* `ptp`: Creates a veth pair.
* `vlan`: Allocates a vlan device.


### IPAM: IP address allocation(官方)
* `dhcp`: Runs a daemon on the host to make DHCP requests on behalf of the container
* `host-local`: maintains a local database of allocated IPs
* `tdipam`: 负责将分配容器的IP地址存储到etcd中 (该代码库)

### Meta: other plugins
* `flannel`: generates an interface corresponding to a flannel config file
* `tuning`: Tweaks sysctl parameters of an existing interface
* `portmap`: An iptables-based portmapping plugin. Maps ports from the host's address space to the container.

### 逻辑导图
![image](https://github.com/panacena/mengqu/blob/master/readme/Screenshot_2016-07-10-22-17- 15_zkk.com.mengqu.png)


###我们该如何使用?
