package system

// Ab23Event is an immutable, schema-unqualified raw observation. Its request
// carries identity context; the reply is not treated as echoing that identity.
// No family receives qualified interpretation or physical-support status.
type Ab23Event struct {
	ab23Family   string
	ab23Request  string
	ab23Response string
	ab23Failure  string
	ab23Valid    bool
}

// Ab23DecodeEvent retains all input bytes, including failed or malformed replies.
// Request order is opcode, system type, instance, address, raw code. Instance
// and code are unrestricted u8 values; code has no calendar interpretation.
func Ab23DecodeEvent(ab23Family string, ab23Request, ab23Response []byte, ab23Failure error) Ab23Event {
	ab23Result := Ab23Event{
		ab23Family: ab23Family, ab23Request: string(ab23Request),
		ab23Response: string(ab23Response),
	}
	if ab23Failure != nil {
		ab23Result.ab23Failure = ab23Failure.Error()
		return ab23Result
	}
	if len(ab23Request) != 5 || len(ab23Response) != 8 {
		return ab23Result
	}
	if ab23Request[0] != 0x09 && ab23Request[0] != 0x0B {
		return ab23Result
	}
	ab23MaxAddress := byte(0)
	switch ab23Request[1] {
	case 0x00:
		ab23MaxAddress = 3
	case 0x01, 0x03:
		ab23MaxAddress = 2
	default:
		return ab23Result
	}
	ab23Result.ab23Valid = ab23Request[3] >= 1 && ab23Request[3] <= ab23MaxAddress
	return ab23Result
}

func (ab23Event Ab23Event) Ab23Family() string   { return ab23Event.ab23Family }
func (ab23Event Ab23Event) Ab23Request() []byte  { return []byte(ab23Event.ab23Request) }
func (ab23Event Ab23Event) Ab23Response() []byte { return []byte(ab23Event.ab23Response) }
func (ab23Event Ab23Event) Ab23Failure() string  { return ab23Event.ab23Failure }
func (ab23Event Ab23Event) Ab23Valid() bool      { return ab23Event.ab23Valid }

func (ab23Event Ab23Event) Ab23Configuration() (byte, bool) {
	if !ab23Event.ab23Valid {
		return 0, false
	}
	return ab23Event.ab23Response[0], true
}

func (ab23Event Ab23Event) Ab23Raw() ([7]byte, bool) {
	var ab23Raw [7]byte
	if !ab23Event.ab23Valid {
		return ab23Raw, false
	}
	copy(ab23Raw[:], ab23Event.ab23Response[1:])
	return ab23Raw, true
}

// Ab23EventPair retains each member independently, even when incomplete. A zero
// event represents a missing member. Completeness describes only the raw software
// contract, not schema qualification, fact grade, or physical support.
type Ab23EventPair struct {
	ab23First  Ab23Event
	ab23Second Ab23Event
}

func Ab23PairEvents(ab23First, ab23Second Ab23Event) Ab23EventPair {
	return Ab23EventPair{ab23First: ab23First, ab23Second: ab23Second}
}

func (ab23Pair Ab23EventPair) Ab23First() Ab23Event  { return ab23Pair.ab23First }
func (ab23Pair Ab23EventPair) Ab23Second() Ab23Event { return ab23Pair.ab23Second }

func (ab23Pair Ab23EventPair) Ab23Complete() bool {
	ab23First, ab23Second := ab23Pair.ab23First, ab23Pair.ab23Second
	return ab23First.ab23Valid && ab23Second.ab23Valid &&
		ab23First.ab23Family == ab23Second.ab23Family &&
		ab23First.ab23Request[1:] == ab23Second.ab23Request[1:] &&
		ab23First.ab23Request[0] != ab23Second.ab23Request[0]
}
