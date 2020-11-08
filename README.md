# ippool

The ippool package provides a basic utility to manage and track a pools of ip addresses and their assignment state.

## usage:

### create a new pool 

InitPrefixPool() will return a pointer to a map of prefixes

```
pool := ippool.InitPrefixPool()
``` 

 

### register a new prefix

ippool expects a net.IPNet struct to register a prefix
you can convert a regular IP Address with prefix using golangs net.ParseCidr function


```
// for IPv4
 _, IPNet, _ := net.ParseCIDR("172.16.10.96/28")
// or for IPv6
 _, IPNet, _ := net.ParseCIDR("FE80::0/96")

``` 

then register the prefix inside the pool

``` 
ippool.RegisterPrefix(&pool, IPnet)
```




### request a new address 

RequestIP() expects a pointer to the pool and the according net.IPNet Prefix.
It will return an ip address as net.IP.

```
ipaddr, err := ippool.RequestIP(&pool, IPNet)
```

### release an assigned address

```
err := ippool.ReleaseIP(&pool, IPnet, ipaddr)
```

### check if an ip is in use

```
bool := IsIPInUse(&pool, IPnet, ipaddr)
```