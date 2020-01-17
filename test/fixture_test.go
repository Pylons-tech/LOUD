package loudFixture

import (
	"testing"

	pylonsFixture "github.com/MikeSofaer/pylons/cmd/fixtures_test"
	pylonSDK "github.com/MikeSofaer/pylons/cmd/test"
)

func TestFixturesViaCLI(t *testing.T) {
	pylonSDK.CLIOpts.CustomNode = "35.223.7.2:26657"
	pylonsFixture.FixtureTestOpts.IsParallel = false
	pylonsFixture.RunTestScenarios("scenarios", t)
}
