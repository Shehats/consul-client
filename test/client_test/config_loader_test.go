package client_test

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/consul-client/client"
	"github.com/stretchr/testify/assert"
)

func TestReadServerConfigFromYaml(t *testing.T) {
	test_resource := "resources/test_server.yml"
	_, x, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(x), "..")
	test_resource = fmt.Sprintf("%s/%s", dir, test_resource)
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	srvs, errs := client.ReadServerConfigFromYaml(test_resource)

	assert.Nil(t, errs, "An error has occured %v", errs)
	assert.NotNil(t, srvs, "Servers assertion failed")
	assert.Len(t, srvs, 2, "Number of servers isn't accurate %v", len(srvs))

	// Check from URI
	srvs, errs = client.ReadServerConfigFromYaml("https://raw.githubusercontent.com/Shehats/consul-client/master/test/resources/test_server.yml")
	assert.Nil(t, errs, "An error has occured %v", errs)
	assert.NotNil(t, srvs, "Servers assertion failed")
	assert.Len(t, srvs, 2, "Number of servers isn't accurate %v", len(srvs))
}
