package cloudformer

import (
	"net"
	"time"
)

type CloudFormer interface {
	InternetGateway() InternetGateway
	VPC() VPC
	ElasticIP(domain string) ElasticIP
	LoadBalancer(name string) LoadBalancer
}

type InternetGateway interface{}

type DHCPOptions struct {
	DomainNameServers []string
}

type VPC interface {
	Network(*net.IPNet)

	AttachInternetGateway(InternetGateway)
	AssociateDHCPOptions(DHCPOptions)

	Subnet(name string) Subnet
	SecurityGroup(name string) SecurityGroup
}

type Subnet interface {
	Network(*net.IPNet)
	AvailabilityZone(string)

	Instance(name string) Instance
	RouteTable() RouteTable
}

type SecurityGroup interface {
	Ingress(ProtocolType, *net.IPNet, uint16, uint16)
}

type RouteTable interface {
	InternetGateway(InternetGateway)
	Instance(Instance)
}

type Instance interface {
	Type(string)
	Image(string)
	PrivateIP(net.IP)
	KeyPair(string)
	SecurityGroup(SecurityGroup)
}

type ElasticIP interface {
	AttachTo(Instance)
}

type LoadBalancer interface {
	Listener(ProtocolType, uint16, ProtocolType, uint16)
	HealthCheck(ProtocolType, uint16, time.Duration, time.Duration, int, int)
	Subnet(Subnet)
	SecurityGroup(SecurityGroup)
}

type ProtocolType string

const TCP = ProtocolType("tcp")
const UDP = ProtocolType("udp")

func CIDR(cidr string) *net.IPNet {
	_, net, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}

	return net
}

func IP(addr string) net.IP {
	ip := net.ParseIP(addr)
	if ip == nil {
		panic("invalid ip: " + addr)
	}

	return ip
}

func Form(f CloudFormer) {
	zone1 := "us-east-1a"

	vpc := f.VPC()
	vpc.Network(CIDR("10.10.0.0/16"))

	vpcGateway := f.InternetGateway()

	vpc.AttachInternetGateway(vpcGateway)

	openSecurityGroup := vpc.SecurityGroup("open")
	boshSecurityGroup := vpc.SecurityGroup("bosh")
	internalSecurityGroup := vpc.SecurityGroup("internal")
	webSecurityGroup := vpc.SecurityGroup("web")

	for _, group := range []SecurityGroup{
		openSecurityGroup,
		boshSecurityGroup,
		internalSecurityGroup,
	} {
		group.Ingress(TCP, CIDR("0.0.0.0/0"), 0, 65535)
		group.Ingress(UDP, CIDR("0.0.0.0/0"), 0, 65535)
	}

	webSecurityGroup.Ingress(TCP, CIDR("0.0.0.0/0"), 80, 80)
	webSecurityGroup.Ingress(TCP, CIDR("0.0.0.0/0"), 8080, 8080)

	boshZ1 := vpc.Subnet("BOSH")
	boshZ1.Network(CIDR("10.10.0.0/24"))
	boshZ1.AvailabilityZone(zone1)
	boshZ1.RouteTable().InternetGateway(vpcGateway)

	droneELBZ1 := vpc.Subnet("DroneELB")
	droneELBZ1.Network(CIDR("10.10.2.0/24"))
	droneELBZ1.AvailabilityZone(zone1)
	droneELBZ1.RouteTable().InternetGateway(vpcGateway)

	boshNAT := boshZ1.Instance("NAT")
	boshNAT.Type("m1.small")
	boshNAT.Image("ami-something")
	boshNAT.PrivateIP(IP("10.10.0.10"))
	boshNAT.KeyPair("bosh")
	boshNAT.SecurityGroup(openSecurityGroup)

	droneZ1 := vpc.Subnet("Drone")
	droneZ1.Network(CIDR("10.10.16.0/20"))
	droneZ1.AvailabilityZone(zone1)
	droneZ1.RouteTable().Instance(boshNAT)

	balancer := f.LoadBalancer("Drone")
	balancer.Listener(TCP, 80, TCP, 80)
	balancer.Listener(TCP, 8080, TCP, 8080)
	balancer.HealthCheck(TCP, 80, 5*time.Second, 30*time.Second, 10, 2)
	balancer.Subnet(droneELBZ1)
	balancer.SecurityGroup(webSecurityGroup)
}
