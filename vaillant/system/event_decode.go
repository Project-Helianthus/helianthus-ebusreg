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

func (eventDecodeEvent Event) EventDecodeFamily() string { return eventDecodeEvent.eventDecodeFamily }
func (eventDecodeEvent Event) EventDecodeRequest() []byte {
	return []byte(eventDecodeEvent.eventDecodeRequest)
}
func (eventDecodeEvent Event) EventDecodeResponse() []byte {
	return []byte(eventDecodeEvent.eventDecodeResponse)
}
func (eventDecodeEvent Event) EventDecodeFailure() string { return eventDecodeEvent.eventDecodeFailure }
func (eventDecodeEvent Event) EventDecodeValid() bool     { return eventDecodeEvent.eventDecodeValid }

func (eventDecodeEvent Event) EventDecodeConfiguration() (byte, bool) {
	if !eventDecodeEvent.eventDecodeValid {
		return 0, false
	}
	return eventDecodeEvent.eventDecodeResponse[0], true
}

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

func PairEvents(eventDecodeFirst, eventDecodeSecond Event) EventPair {
	return EventPair{eventDecodeFirst: eventDecodeFirst, eventDecodeSecond: eventDecodeSecond}
}

func (eventDecodePair EventPair) EventDecodeFirst() Event  { return eventDecodePair.eventDecodeFirst }
func (eventDecodePair EventPair) EventDecodeSecond() Event { return eventDecodePair.eventDecodeSecond }

func (eventDecodePair EventPair) EventDecodeComplete() bool {
	eventDecodeFirst, eventDecodeSecond := eventDecodePair.eventDecodeFirst, eventDecodePair.eventDecodeSecond
	return eventDecodeFirst.eventDecodeValid && eventDecodeSecond.eventDecodeValid &&
		eventDecodeFirst.eventDecodeFamily == eventDecodeSecond.eventDecodeFamily &&
		eventDecodeFirst.eventDecodeRequest[1:] == eventDecodeSecond.eventDecodeRequest[1:] &&
		eventDecodeFirst.eventDecodeRequest[0] != eventDecodeSecond.eventDecodeRequest[0]
}
