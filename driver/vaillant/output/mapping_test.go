package output

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

type gatewayTestMapping struct {
	Carried []struct {
		GatewayFile string `json:"gateway_file"`
		GatewayTest string `json:"gateway_test"`
		DriverFile  string `json:"driver_file"`
		DriverTest  string `json:"driver_test"`
	}
	LeftInGateway []struct {
		GatewayFile string `json:"gateway_file"`
		GatewayTest string `json:"gateway_test"`
		Reason      string
	} `json:"left_in_gateway"`
}

func driverTestNames(t *testing.T, file string) map[string]bool {
	t.Helper()
	parsed, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	names := make(map[string]bool)
	for _, declaration := range parsed.Decls {
		if function, ok := declaration.(*ast.FuncDecl); ok && function.Recv == nil && strings.HasPrefix(function.Name.Name, "Test") {
			names[function.Name.Name] = true
		}
	}
	return names
}

func TestVaillantGatewayTestMappingNamesEveryCarriedDriverTest(t *testing.T) {
	mapping := readFixture[gatewayTestMapping](t, "gateway-test-mapping.json")
	mapped := make(map[string]string)
	gateway := make(map[string]bool)
	for _, entry := range mapping.Carried {
		if !driverTestNames(t, entry.DriverFile)[entry.DriverTest] {
			t.Errorf("%s does not declare %s", entry.DriverFile, entry.DriverTest)
		}
		if previous, dup := mapped[entry.DriverTest]; dup {
			t.Errorf("%s is the driver test of %s and of %s", entry.DriverTest, previous, entry.GatewayTest)
		}
		mapped[entry.DriverTest] = entry.GatewayTest
		if gateway[entry.GatewayTest] {
			t.Errorf("%s is listed twice", entry.GatewayTest)
		}
		gateway[entry.GatewayTest] = true
	}
	for _, entry := range mapping.LeftInGateway {
		if gateway[entry.GatewayTest] {
			t.Errorf("%s is both carried and left in the gateway", entry.GatewayTest)
		}
		gateway[entry.GatewayTest] = true
		if entry.Reason == "" {
			t.Errorf("%s left in the gateway has no reason", entry.GatewayTest)
		}
	}
	for _, file := range []string{"views_test.go", "cache_test.go"} {
		for name := range driverTestNames(t, file) {
			if _, ok := mapped[name]; !ok {
				t.Errorf("%s of %s is not in the mapping", name, file)
			}
		}
	}
}
