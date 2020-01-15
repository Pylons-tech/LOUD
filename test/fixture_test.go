package loudFixture

import (
	"testing"

	pylonsFixture "github.com/Pylons-tech/pylons/cmd/fixtures_test"
)

func TestFixturesViaCLI(t *testing.T) {
	pylonsFixture.RunTestScenarios("scenarios", t)
}
