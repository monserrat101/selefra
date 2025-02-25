package module

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestProviderBlock_MarshalYAML(t *testing.T) {
	block := NewProviderBlock()

	block.Name = "foo"
	block.Cache = "1d"
	block.Provider = "aws"
	block.ProvidersConfigYamlString = `    #  Optional, Repeated. Add an accounts block for every account you want to assume-role into and fetch data from.
    accounts:
      #     Optional. User identification
      - account_name: <UNIQUE ACCOUNT IDENTIFIER>
        #    Optional. Named profile in config or credential file from where Selefra should grab credentials
        shared_config_profile: < PROFILE_NAME >
        #    Optional. Location of shared configuration files
        shared_config_files:
          - <FILE_PATH>
        #   Optional. Location of shared credentials files
        shared_credentials_files:
          - <FILE_PATH>
        #    Optional. Role ARN we want to assume when accessing this account
        role_arn: < YOUR_ROLE_ARN >
        #    Optional. Named role session to grab specific operation under the assumed role
        role_session_name: <SESSION_NAME>
        #    Optional. Any outside of the org account id that has additional control
        external_id: <ID>
        #    Optional. Designated region of servers
        default_region: <REGION_CODE>
        #    Optional. by default assumes all regions
        regions:
          - us-east-1
          - us-west-2
    #    The maximum number of times that a request will be retried for failures. Defaults to 10 retry attempts.
    max_attempts: 10
    #    The maximum back off delay between attempts. The backoff delays exponentially with a jitter based on the number of attempts. Defaults to 30 seconds.
    max_backoff: 30`

	out, err := yaml.Marshal(block)
	assert.Nil(t, err)
	t.Log(string(out))
}
