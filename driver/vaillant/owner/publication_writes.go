package owner

import (
	"context"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

// The poller's configuration and time-program writers return native
// results; the writer adapters below gate them behind the admitted source
// and map the native results to the GraphQL and MCP forms with the same
// field names.

// semanticWriterSourceNotAdmittedError is the refusal of a write while no
// source selection is active.
const semanticWriterSourceNotAdmittedError = "source selection not active"

// admittedSourceProvider reports the admitted source and whether one is
// admitted.
type admittedSourceProvider func() (byte, bool)

func admittedSourceActive(provider admittedSourceProvider) bool {
	if provider == nil {
		return false
	}
	source, admitted := provider()
	return admitted && source != 0
}

// admittedSemanticWriter gates the owner's configuration and time-program
// writers behind the admitted source. The configuration results keep the
// native outcome; the time-program results take their output form.
type admittedSemanticWriter struct {
	boiler   VaillantBoilerConfigWriter
	system   VaillantSystemConfigWriter
	schedule VaillantScheduleWriter
	admitted admittedSourceProvider
}

func (writer admittedSemanticWriter) SetBoilerConfig(ctx context.Context, fieldName string, rawValue string) vaillantConfigWriteResult {
	if !admittedSourceActive(writer.admitted) {
		return vaillantConfigWriteResult{Success: false, Error: semanticWriterSourceNotAdmittedError}
	}
	if writer.boiler == nil {
		return vaillantConfigWriteResult{Success: false, Error: "boiler config writer unavailable"}
	}
	return writer.boiler.SetBoilerConfig(ctx, fieldName, rawValue)
}

func (writer admittedSemanticWriter) SetSystemConfig(ctx context.Context, fieldName string, rawValue string) vaillantConfigWriteResult {
	if !admittedSourceActive(writer.admitted) {
		return vaillantConfigWriteResult{Success: false, Error: semanticWriterSourceNotAdmittedError}
	}
	if writer.system == nil {
		return vaillantConfigWriteResult{Success: false, Error: "system config writer unavailable"}
	}
	return writer.system.SetSystemConfig(ctx, fieldName, rawValue)
}

func (writer admittedSemanticWriter) SetZoneTimeProgram(ctx context.Context, zone int, weekday int, slots []output.MCPTimeProgramSlot) (*output.MCPTimeProgramWriteResult, error) {
	if !admittedSourceActive(writer.admitted) {
		return &output.MCPTimeProgramWriteResult{Success: false, Error: semanticWriterSourceNotAdmittedError}, nil
	}
	if writer.schedule == nil {
		return &output.MCPTimeProgramWriteResult{Success: false, Error: "schedule writer unavailable"}, nil
	}
	result, err := writer.schedule.SetZoneTimeProgram(ctx, zone, weekday, timeProgramSlotsFromMCP(slots))
	return mcpTimeProgramWriteResult(result, err)
}

func (writer admittedSemanticWriter) SetDhwTimeProgram(ctx context.Context, weekday int, slots []output.MCPTimeProgramSlot) (*output.MCPTimeProgramWriteResult, error) {
	if !admittedSourceActive(writer.admitted) {
		return &output.MCPTimeProgramWriteResult{Success: false, Error: semanticWriterSourceNotAdmittedError}, nil
	}
	if writer.schedule == nil {
		return &output.MCPTimeProgramWriteResult{Success: false, Error: "schedule writer unavailable"}, nil
	}
	result, err := writer.schedule.SetDhwTimeProgram(ctx, weekday, timeProgramSlotsFromMCP(slots))
	return mcpTimeProgramWriteResult(result, err)
}
