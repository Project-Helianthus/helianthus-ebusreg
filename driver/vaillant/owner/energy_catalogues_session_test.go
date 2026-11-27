package owner

import (
	"errors"
	"time"

	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"
)

type energyCatalogueSession struct {
	calendar semdriver.Calendar
	later    semdriver.Calendar
	calls    int
}

func (s *energyCatalogueSession) Native() semdriver.NativeRegistry             { return nil }
func (s *energyCatalogueSession) SemanticInput() semdriver.SemanticInput       { return nil }
func (s *energyCatalogueSession) ScopedOperations() semdriver.ScopedOperations { return nil }
func (s *energyCatalogueSession) Publisher() semdriver.Publisher               { return nil }
func (s *energyCatalogueSession) RegisterExecutor(semdriver.SettingRoute, semdriver.Executor) (func(), error) {
	return func() {}, nil
}
func (s *energyCatalogueSession) Calendar() semdriver.Calendar {
	s.calls++
	if s.later != nil && s.calls > 2 {
		return s.later
	}
	return s.calendar
}

type energyTestCalendar struct {
	day        int
	generation int
	begins     int
}

func (c *energyTestCalendar) BeginAcquisition() semdriver.CalendarAcquisition {
	c.begins++
	return &energyTestAcquisition{calendar: c, day: c.day, generation: c.generation}
}

type energyTestAcquisition struct {
	calendar   *energyTestCalendar
	day        int
	generation int
}

func (a *energyTestAcquisition) Reference() semdriver.CalendarReference { return energyTestReference{} }
func (a *energyTestAcquisition) Check(calendar semdriver.Calendar) error {
	current, ok := calendar.(*energyTestCalendar)
	if !ok || current != a.calendar || current.day != a.day || current.generation != a.generation {
		return errors.New("calendar guard")
	}
	return nil
}

type energyTestReference struct{}

func (energyTestReference) Instant() time.Time {
	return time.Date(2024, time.March, 15, 12, 0, 0, 0, time.UTC)
}
func (energyTestReference) Date() semdriver.CalendarDate {
	return semdriver.CalendarDate{Year: 2024, Month: time.March, Day: 15}
}
func (energyTestReference) Year() int         { return 2024 }
func (energyTestReference) PreviousYear() int { return 2023 }
func (energyTestReference) Month() semdriver.CalendarMonth {
	return semdriver.CalendarMonth{Year: 2024, Month: time.March}
}
func (energyTestReference) PreviousMonth() semdriver.CalendarMonth {
	return semdriver.CalendarMonth{Year: 2024, Month: time.February}
}
