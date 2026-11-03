package system

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func ab23Fixture(ab23Opcode byte) Ab23Event {
	ab23Response := []byte{0xA0, 0x01, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67}
	if ab23Opcode == 0x0B {
		ab23Response = []byte{0xB0, 0x78, 0x89, 0x9A, 0xAB, 0xBC, 0xCD, 0xDE}
	}
	return Ab23DecodeEvent("synthetic-family", []byte{ab23Opcode, 0, 0xFF, 1, 0xFE}, ab23Response, nil)
}

func TestAb23ExactRawAndOwnership(ab23T *testing.T) {
	for _, ab23Opcode := range []byte{0x09, 0x0B} {
		ab23T.Run(fmt.Sprintf("ab23-op-%02x", ab23Opcode), func(ab23T *testing.T) {
			ab23Expected := ab23Fixture(ab23Opcode)
			ab23Request := []byte{ab23Opcode, 0, 0xFF, 1, 0xFE}
			ab23Response := []byte{0xA0, 0x01, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67}
			if ab23Opcode == 0x0B {
				ab23Response = []byte{0xB0, 0x78, 0x89, 0x9A, 0xAB, 0xBC, 0xCD, 0xDE}
			}
			ab23Event := Ab23DecodeEvent("synthetic-family", ab23Request, ab23Response, nil)
			ab23Configuration, ab23OK := ab23Event.Ab23Configuration()
			if !ab23OK || ab23Configuration != ab23Response[0] {
				ab23T.Fatal("ab23 configuration changed")
			}
			ab23Raw, ab23OK := ab23Event.Ab23Raw()
			if !ab23OK || !bytes.Equal(ab23Raw[:], ab23Response[1:]) {
				ab23T.Fatal("ab23 raw positions changed")
			}
			if !ab23Event.Ab23Valid() || ab23Event.Ab23Failure() != "" || ab23Event.Ab23Family() != "synthetic-family" {
				ab23T.Fatal("ab23 observation state")
			}
			for ab23Index := range ab23Response {
				ab23Response[ab23Index] ^= 0xFF
			}
			for ab23Index := range ab23Request {
				ab23Request[ab23Index] ^= 0xFF
			}
			if ab23Event != ab23Expected {
				ab23T.Fatal("ab23 input alias")
			}
			ab23Output := ab23Event.Ab23Response()
			ab23Output[0] = 0
			ab23Output = ab23Event.Ab23Request()
			ab23Output[0] = 0
			ab23Raw[0] = 0
			if ab23Event != ab23Expected {
				ab23T.Fatal("ab23 output alias")
			}
		})
	}
}

func TestAb23PairDimensions(ab23T *testing.T) {
	ab23First, ab23Second := ab23Fixture(0x09), ab23Fixture(0x0B)
	if !Ab23PairEvents(ab23First, ab23Second).Ab23Complete() || !Ab23PairEvents(ab23Second, ab23First).Ab23Complete() {
		ab23T.Fatal("ab23 valid pair")
	}
	for _, ab23Dimension := range []string{"family", "system", "instance", "address", "code", "opcode"} {
		ab23T.Run("ab23-"+ab23Dimension, func(ab23T *testing.T) {
			ab23Family, ab23Request := ab23Second.Ab23Family(), ab23Second.Ab23Request()
			switch ab23Dimension {
			case "family":
				ab23Family = "another-synthetic-family"
			case "system":
				ab23Request[1] = 1
			case "instance":
				ab23Request[2] = 0xFE
			case "address":
				ab23Request[3] = 2
			case "code":
				ab23Request[4] = 0xFD
			case "opcode":
				ab23Request[0] = 0x09
			}
			ab23Changed := Ab23DecodeEvent(ab23Family, ab23Request, ab23Second.Ab23Response(), nil)
			if !ab23Changed.Ab23Valid() {
				ab23T.Fatal("ab23 cross-pair control must be valid")
			}
			ab23Pair := Ab23PairEvents(ab23First, ab23Changed)
			if ab23Pair.Ab23Complete() {
				ab23T.Fatal("ab23 crossed pair dimension")
			}
			if ab23Pair.Ab23First() != ab23First || ab23Pair.Ab23Second() != ab23Changed {
				ab23T.Fatal("ab23 member loss")
			}
			if ab23Dimension != "opcode" {
				ab23Request[0] = 0x09
				ab23Matching := Ab23DecodeEvent(ab23Family, ab23Request, ab23First.Ab23Response(), nil)
				if !Ab23PairEvents(ab23Matching, ab23Changed).Ab23Complete() {
					ab23T.Fatal("ab23 changed identity matching control")
				}
			} else if !Ab23PairEvents(ab23Second, ab23Changed).Ab23Complete() {
				ab23T.Fatal("ab23 changed opcode matching control")
			}
		})
	}
}

func TestAb23IncompleteMembers(ab23T *testing.T) {
	ab23Good := ab23Fixture(0x09)
	ab23Other := ab23Fixture(0x0B)
	for ab23Length := 0; ab23Length <= 9; ab23Length++ {
		if ab23Length == 8 {
			continue
		}
		ab23T.Run(fmt.Sprintf("ab23-length-%d", ab23Length), func(ab23T *testing.T) {
			ab23Response := bytes.Repeat([]byte{0xA5}, ab23Length)
			ab23Bad := Ab23DecodeEvent(ab23Other.Ab23Family(), ab23Other.Ab23Request(), ab23Response, nil)
			ab23CheckIncomplete(ab23T, ab23Good, ab23Bad)
			if !bytes.Equal(ab23Bad.Ab23Response(), ab23Response) {
				ab23T.Fatal("ab23 malformed bytes lost")
			}
		})
	}
	ab23Failed := Ab23DecodeEvent(ab23Other.Ab23Family(), ab23Other.Ab23Request(), ab23Other.Ab23Response(), errors.New("ab23 synthetic failure"))
	ab23CheckIncomplete(ab23T, ab23Good, ab23Failed)
	if ab23Failed.Ab23Failure() != "ab23 synthetic failure" || !bytes.Equal(ab23Failed.Ab23Response(), ab23Other.Ab23Response()) {
		ab23T.Fatal("ab23 failure evidence lost")
	}
	ab23CheckIncomplete(ab23T, ab23Good, Ab23Event{})
	ab23CheckIncomplete(ab23T, Ab23Event{}, Ab23Event{})
}

func ab23CheckIncomplete(ab23T *testing.T, ab23First, ab23Second Ab23Event) {
	ab23T.Helper()
	for _, ab23Pair := range []Ab23EventPair{Ab23PairEvents(ab23First, ab23Second), Ab23PairEvents(ab23Second, ab23First)} {
		if ab23Pair.Ab23Complete() {
			ab23T.Fatal("ab23 incomplete pair admitted")
		}
	}
	ab23Pair := Ab23PairEvents(ab23First, ab23Second)
	if ab23Pair.Ab23First() != ab23First || ab23Pair.Ab23Second() != ab23Second {
		ab23T.Fatal("ab23 independent member lost")
	}
	if !ab23Second.Ab23Valid() {
		if _, ab23OK := ab23Second.Ab23Configuration(); ab23OK {
			ab23T.Fatal("ab23 invalid configuration admitted")
		}
		if _, ab23OK := ab23Second.Ab23Raw(); ab23OK {
			ab23T.Fatal("ab23 invalid raw admitted")
		}
	}
}

func TestAb23ContextAdmissionAndUnknownFamily(ab23T *testing.T) {
	ab23Reply := ab23Fixture(0x09).Ab23Response()
	for _, ab23System := range []byte{0, 1, 3, 2, 0xFF} {
		for ab23Address := byte(0); ab23Address <= 4; ab23Address++ {
			ab23Request := []byte{0x09, ab23System, 0xFF, ab23Address, 0xFF}
			ab23Event := Ab23DecodeEvent("unknown-family", ab23Request, ab23Reply, nil)
			ab23Want := (ab23System == 0 && ab23Address >= 1 && ab23Address <= 3) || ((ab23System == 1 || ab23System == 3) && ab23Address >= 1 && ab23Address <= 2)
			if ab23Event.Ab23Valid() != ab23Want || !bytes.Equal(ab23Event.Ab23Request(), ab23Request) || !bytes.Equal(ab23Event.Ab23Response(), ab23Reply) || ab23Event.Ab23Family() != "unknown-family" {
				ab23T.Fatalf("ab23 context %v", ab23Request)
			}
		}
	}
	for _, ab23Request := range [][]byte{nil, {0x09}, {0x09, 0, 0xFF, 1}, {0x09, 0, 0xFF, 1, 0xFF, 0}, {0x03, 0, 0xFF, 1, 0xFF}} {
		ab23Event := Ab23DecodeEvent("unknown-family", ab23Request, ab23Reply, nil)
		if ab23Event.Ab23Valid() || !bytes.Equal(ab23Event.Ab23Request(), ab23Request) {
			ab23T.Fatal("ab23 malformed request admission")
		}
		ab23CheckIncomplete(ab23T, ab23Fixture(0x0B), ab23Event)
	}
}
