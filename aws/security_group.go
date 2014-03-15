package aws

import (
	"net"

	"github.com/vito/cloudformer"
	"github.com/vito/cloudformer/aws/models"
)

type SecurityGroup struct {
	name  string
	model *models.SecurityGroup
}

func (securityGroup SecurityGroup) Ingress(
	protocol cloudformer.ProtocolType,
	network *net.IPNet,
	fromPort uint16,
	toPort uint16,
) {
	securityGroup.model.SecurityGroupIngress =
		&models.SecurityGroupIngress{
			CidrIp:     network.String(),
			IpProtocol: protocol,
			FromPort:   fromPort,
			ToPort:     toPort,
		}
}

func (securityGroup SecurityGroup) Egress(
	protocol cloudformer.ProtocolType,
	network *net.IPNet,
	fromPort uint16,
	toPort uint16,
) {
	securityGroup.model.SecurityGroupIngress =
		&models.SecurityGroupEgress{
			CidrIp:     network.String(),
			IpProtocol: protocol,
			FromPort:   fromPort,
			ToPort:     toPort,
		}
}
