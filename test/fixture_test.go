package loudFixture

import (
	"testing"

	pylonsFixture "github.com/Pylons-tech/pylons/cmd/fixtures_test"
	pylonSDK "github.com/Pylons-tech/pylons/cmd/test"
)

func TestFixturesViaCLI(t *testing.T) {
	pylonSDK.CLIOpts.CustomNode = "35.223.7.2:26657"
	pylonsFixture.FixtureTestOpts.IsParallel = false
	pylonsFixture.RunTestScenarios("scenarios", t)
}
