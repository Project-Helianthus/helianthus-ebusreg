package system

// Event is an immutable, schema-unqualified raw observation. Its request
// carries identity context; the reply is not treated as echoing that identity.
// No family receives qualified interpretation or physical-support status.
type Event struct {
	eventDecodeFamily   string
	eventDecodeRequest  string
	eventDecodeResponse string
	eventDecodeFailure  string
	eventDecodeValid    bool
}

// DecodeEvent retains all input bytes, including failed or malformed replies.
// Request order is opcode, system type, instance, address, raw code. Instance
// and code are unrestricted u8 values; code has no calendar interpretation.
func DecodeEvent(eventDecodeFamily string, eventDecodeRequest, eventDecodeResponse []byte, eventDecodeFailure error) Event {
	eventDecodeResult := Event{
		eventDecodeFamily: eventDecodeFamily, eventDecodeRequest: string(eventDecodeRequest),
		eventDecodeResponse: string(eventDecodeResponse),
	}
	if eventDecodeFailure != nil {
		eventDecodeResult.eventDecodeFailure = eventDecodeFailure.Error()
		return eventDecodeResult
	}
	if len(eventDecodeRequest) != 5 || len(eventDecodeResponse) != 8 {
		return eventDecodeResult
	}
	if eventDecodeRequest[0] != 0x09 && eventDecodeRequest[0] != 0x0B {
		return eventDecodeResult
	}
	eventDecodeMaxAddress := byte(0)
	switch eventDecodeRequest[1] {
	case 0x00:
		eventDecodeMaxAddress = 3
	case 0x01, 0x03:
		eventDecodeMaxAddress = 2
	default:
		return eventDecodeResult
	}
	eventDecodeResult.eventDecodeValid = eventDecodeRequest[3] >= 1 && eventDecodeRequest[3] <= eventDecodeMaxAddress
	return eventDecodeResult
}

// EventDecodeFamily returns the device family the event was recorded for.
func (eventDecodeEvent Event) EventDecodeFamily() string { return eventDecodeEvent.eventDecodeFamily }

// EventDecodeRequest returns a copy of the raw request bytes.
func (eventDecodeEvent Event) EventDecodeRequest() []byte {
	return []byte(eventDecodeEvent.eventDecodeRequest)
}

// EventDecodeResponse returns a copy of the raw reply bytes, even when the
// event is not valid.
func (eventDecodeEvent Event) EventDecodeResponse() []byte {
	return []byte(eventDecodeEvent.eventDecodeResponse)
}

// EventDecodeFailure returns the text of the exchange error, or "" when the
// exchange did not fail.
func (eventDecodeEvent Event) EventDecodeFailure() string { return eventDecodeEvent.eventDecodeFailure }

// EventDecodeValid reports whether the exchange succeeded with a five-byte
// request (read opcode 0x09 or 0x0B, a known system type and an address in
// its range) and an eight-byte reply.
func (eventDecodeEvent Event) EventDecodeValid() bool { return eventDecodeEvent.eventDecodeValid }

// EventDecodeConfiguration returns the first reply byte, uninterpreted, and
// true; it returns false for an event that is not valid.
func (eventDecodeEvent Event) EventDecodeConfiguration() (byte, bool) {
	if !eventDecodeEvent.eventDecodeValid {
		return 0, false
	}
	return eventDecodeEvent.eventDecodeResponse[0], true
}

// EventDecodeRaw returns the seven reply bytes after the first one and true;
// it returns false for an event that is not valid.
func (eventDecodeEvent Event) EventDecodeRaw() ([7]byte, bool) {
	var eventDecodeRaw [7]byte
	if !eventDecodeEvent.eventDecodeValid {
		return eventDecodeRaw, false
	}
	copy(eventDecodeRaw[:], eventDecodeEvent.eventDecodeResponse[1:])
	return eventDecodeRaw, true
}

// EventPair retains each member independently, even when incomplete. A zero
// event represents a missing member. Completeness describes only the raw software
// contract, not schema qualification, fact grade, or physical support.
type EventPair struct {
	eventDecodeFirst  Event
	eventDecodeSecond Event
}

// PairEvents keeps two events as one pair without checking them; use
// EventDecodeComplete to test the pair.
func PairEvents(eventDecodeFirst, eventDecodeSecond Event) EventPair {
	return EventPair{eventDecodeFirst: eventDecodeFirst, eventDecodeSecond: eventDecodeSecond}
}

// EventDecodeFirst returns the first member, a zero Event when it is missing.
func (eventDecodePair EventPair) EventDecodeFirst() Event { return eventDecodePair.eventDecodeFirst }

// EventDecodeSecond returns the second member, a zero Event when it is
// missing.
func (eventDecodePair EventPair) EventDecodeSecond() Event { return eventDecodePair.eventDecodeSecond }

// EventDecodeComplete reports whether both members are valid, share the
// family and every request byte after the opcode, and were read with
// different opcodes.
func (eventDecodePair EventPair) EventDecodeComplete() bool {
	eventDecodeFirst, eventDecodeSecond := eventDecodePair.eventDecodeFirst, eventDecodePair.eventDecodeSecond
	return eventDecodeFirst.eventDecodeValid && eventDecodeSecond.eventDecodeValid &&
		eventDecodeFirst.eventDecodeFamily == eventDecodeSecond.eventDecodeFamily &&
		eventDecodeFirst.eventDecodeRequest[1:] == eventDecodeSecond.eventDecodeRequest[1:] &&
		eventDecodeFirst.eventDecodeRequest[0] != eventDecodeSecond.eventDecodeRequest[0]
}
