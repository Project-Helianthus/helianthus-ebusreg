package output

// NativeScheduleTimerSlot represents one scheduled interval.
type NativeScheduleTimerSlot struct {
	StartHour      int
	StartMinute    int
	EndHour        int
	EndMinute      int
	TemperatureC   *float64
	TemperatureRaw *int
}

// NativeScheduleDay represents the native intervals for one day.
type NativeScheduleDay struct {
	Weekday string
	Slots   []NativeScheduleTimerSlot
}

// NativeScheduleConfig represents the native limits and temperature support of a program.
type NativeScheduleConfig struct {
	MaxSlots       int
	TimeResolution int
	MinDuration    int
	HasTemperature bool
	TempSlots      int
	MinTempC       *float64
	MaxTempC       *float64
}

// NativeScheduleProgram represents the program for one zone and circuit kind.
type NativeScheduleProgram struct {
	Zone      int
	HC        string
	Config    *NativeScheduleConfig
	SlotsUsed []int
	Days      []NativeScheduleDay
}

// NativeTimeProgramSlot represents one requested interval of a time-program write.
type NativeTimeProgramSlot struct {
	StartHour    int
	StartMinute  int
	EndHour      int
	EndMinute    int
	TemperatureC *float64
}

// NativeTimeProgramSlotResult represents the write outcome for one interval.
type NativeTimeProgramSlotResult struct {
	SlotIndex int
	Accepted  bool
	ErrorCode int
	ErrorDesc string
}

// NativeTimeProgramWriteResult represents the result of a time-program write.
type NativeTimeProgramWriteResult struct {
	Success     bool
	SlotResults []NativeTimeProgramSlotResult
	Error       string
}
