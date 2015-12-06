package aws

import (
	"github.com/vito/cloudformer"
	"github.com/vito/cloudformer/aws/models"
)

type RecordSet struct {
	name  string
	model *models.RecordSet
}

func (recordSet RecordSet) HostedZoneName(hostedZoneName string) {
	recordSet.model.HostedZoneName = hostedZoneName
}

func (recordSet RecordSet) Name(name string) {
	recordSet.model.Name = name
}

func (recordSet RecordSet) PointTo(elasticIP cloudformer.ElasticIP) {
	recordSet.model.RecordSetType = "A"
	recordSet.model.ResourceRecords = []models.Hash{models.Ref(elasticIP.(ElasticIP).name + "EIP")}
}

func (recordSet RecordSet) TTL(ttl int) {
	recordSet.model.TTL = ttl
}
