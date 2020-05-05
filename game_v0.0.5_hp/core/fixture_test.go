package loudFixture

import (
	"flag"
	"testing"

	pylonsFixture "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test"
	pylonSDK "github.com/Pylons-tech/pylons_sdk/cmd/test"
)

var runSerialMode bool = false
var connectLocalDaemon bool = false

func init() {
	flag.BoolVar(&runSerialMode, "runserial", false, "true/false value to check if test will be running in parallel")
	flag.BoolVar(&connectLocalDaemon, "locald", false, "true/false value to check if test will be connecting to local daemon")
}

func TestFixturesViaCLI(t *testing.T) {
	flag.Parse()
	if connectLocalDaemon {
		pylonSDK.CLIOpts.CustomNode = "localhost:26657"
	} else {
		pylonSDK.CLIOpts.CustomNode = "35.223.7.2:26657"
	}
	pylonsFixture.FixtureTestOpts.IsParallel = !runSerialMode
	pylonsFixture.RunTestScenarios("scenarios", t)
}
