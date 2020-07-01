package loudfixture

import (
	"flag"
	"strings"
	"testing"

	fixturetestSDK "github.com/Pylons-tech/pylons_sdk/cmd/fixture_utils"
	inttestSDK "github.com/Pylons-tech/pylons_sdk/cmd/test_utils"
)

var runSerialMode = false
var connectLocalDaemon = false
var useKnownCookbook = false
var useRest = false
var scenarios = ""
var accounts = ""

func init() {
	flag.BoolVar(&runSerialMode, "runserial", false, "true/false value to check if test will be running in parallel")
	flag.BoolVar(&connectLocalDaemon, "locald", false, "true/false value to check if test will be connecting to local daemon")
	flag.BoolVar(&useRest, "userest", false, "use rest endpoint for Tx send")
	flag.BoolVar(&useKnownCookbook, "use-known-cookbook", false, "use existing cookbook or not")
	flag.StringVar(&scenarios, "scenarios", "", "custom scenario file names")
	flag.StringVar(&accounts, "accounts", "", "custom account names")
}

func TestFixturesViaCLI(t *testing.T) {
	flag.Parse()
	if connectLocalDaemon {
		inttestSDK.CLIOpts.CustomNode = "tcp://localhost:26657"
	} else {
		inttestSDK.CLIOpts.CustomNode = "tcp://35.223.7.2:26657"
	}
	fixturetestSDK.FixtureTestOpts.CreateNewCookbook = !useKnownCookbook
	fixturetestSDK.FixtureTestOpts.IsParallel = !runSerialMode
	if useRest {
		inttestSDK.CLIOpts.RestEndpoint = "http://localhost:1317"
	}
	fixturetestSDK.RegisterDefaultActionRunners()
	scenarioFileNames := []string{}
	if len(scenarios) > 0 {
		scenarioFileNames = strings.Split(scenarios, ",")
	}
	fixturetestSDK.FixtureTestOpts.AccountNames = []string{}
	if len(accounts) > 0 {
		fixturetestSDK.FixtureTestOpts.AccountNames = strings.Split(accounts, ",")
	}
	fixturetestSDK.RunTestScenarios("scenarios", scenarioFileNames, t)
}
