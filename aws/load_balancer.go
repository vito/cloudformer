package aws

import (
	"fmt"

	"github.com/vito/cloudformer"
	"github.com/vito/cloudformer/aws/models"
)

type LoadBalancer struct {
	name  string
	model *models.LoadBalancer
}

func (balancer LoadBalancer) Listener(
	protocol cloudformer.ProtocolType,
	port uint16,
	destinationProtocol cloudformer.ProtocolType,
	destinationPort uint16,
) {
	listeners := balancer.model.Listeners.([]interface{})

	listeners = append(listeners, models.LoadBalancerListener{
		Protocol:         string(protocol),
		LoadBalancerPort: fmt.Sprintf("%d", port),
		InstanceProtocol: string(destinationProtocol),
		InstancePort:     fmt.Sprintf("%d", destinationPort),
	})

	balancer.model.Listeners = listeners
}

func (balancer LoadBalancer) HealthCheck(check cloudformer.HealthCheck) {
	balancer.model.HealthCheck = models.LoadBalancerHealthCheck{
		Target:             fmt.Sprintf("%s:%d", check.Protocol, check.Port),
		Interval:           fmt.Sprintf("%d", int(check.Interval.Seconds())),
		Timeout:            fmt.Sprintf("%d", int(check.Timeout.Seconds())),
		HealthyThreshold:   fmt.Sprintf("%d", check.HealthyThreshold),
		UnhealthyThreshold: fmt.Sprintf("%d", check.UnhealthyThreshold),
	}
}

func (balancer LoadBalancer) Subnet(subnet cloudformer.Subnet) {
	subnets := balancer.model.Subnets.([]interface{})

	subnets = append(
		subnets,
		models.Ref(subnet.(Subnet).name+"Subnet"),
	)

	balancer.model.Subnets = subnets
}

func (balancer LoadBalancer) SecurityGroup(group cloudformer.SecurityGroup) {
	securityGroups := balancer.model.SecurityGroups.([]interface{})

	securityGroups = append(
		securityGroups,
		models.Ref(group.(SecurityGroup).name+"SecurityGroup"),
	)

	balancer.model.SecurityGroups = securityGroups
}
