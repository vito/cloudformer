package aws_test

import (
	"encoding/json"

	"github.com/vito/cloudformer/aws"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AWSCloudFormer", func() {
	It("creates a DNS A record for an elastic IP", func() {
		f := aws.New("elastic IP with A record")

		eip := f.ElasticIP("bosh")

		recordSet := f.RecordSet("bosh")
		recordSet.HostedZoneName("test.com.")
		recordSet.Name("bosh.test.com")
		recordSet.TTL(300)
		recordSet.PointTo(eip)

		Expect(json.Marshal(f.Template)).To(MatchJSON(`
{
  "AWSTemplateFormatVersion": "2010-09-09",
  "Description": "elastic IP with A record",
  "Mappings": {},
  "Resources": {
    "boshEIP": {
      "Type": "AWS::EC2::EIP"
    },
    "boshRecordSet": {
      "Properties": {
        "HostedZoneName": "test.com.",
        "Name": "bosh.test.com",
        "ResourceRecords": [
          {
            "Ref": "boshEIP"
          }
        ],
        "TTL": 300,
        "Type": "A"
      },
      "Type": "AWS::Route53::RecordSet"
    }
  }
}
		`))
	})
})
