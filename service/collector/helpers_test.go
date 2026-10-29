package collector

// Values the tests share, named for goconst.
const (
	testAPIVersion           = "application.giantswarm.io/v1alpha1"
	testAppVersion           = "1.0.0"
	testOperatorVersionOld   = "6.3.0"
	testOperatorVersionNew   = "6.4.0"
	testNamespaceABC01       = "abc01"
	testNamespaceTesting     = "testing"
	testCaseFlawless         = "flawless"
	testCaseFlawlessVersions = "flawless with v* versions"
	testTeam                 = "honeybadger"
	testProvider             = "aws"
	testCustomerApp          = "customer-example-" + testAppVersion
	testTeamAtlas            = "team-atlas"
	testExpectedMetrics3     = "testdata/expected.3"
)
