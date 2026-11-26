package owner

import "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"

// BoilerThermalHandle is the boiler thermal handle type of the Vaillant owner.
type BoilerThermalHandle interface {
	Source() semreg.SourceID
	Epoch() semreg.SourceEpochID
	GenerationID() semreg.Uint64
}

// BoilerThermalAssetSpec is the boiler thermal asset spec type of the Vaillant owner.
type BoilerThermalAssetSpec struct {
	AssetID    semreg.AssetID
	SourceID   semreg.SourceID
	Validators []semreg.PackValidator
	OnControl  func() error
}

// BoilerThermalTransaction is the boiler thermal transaction type of the Vaillant owner.
type BoilerThermalTransaction interface {
	Current() (semreg.Snapshot, bool)
	NextSequence() semreg.Uint64
	SourceRetirements() []semreg.SourceEpochID
	Apply(batch semreg.PublicationBatch, monotonic semreg.MonotonicPoint) error
	Commit() error
	Abort()
}

// BoilerThermalHub is the boiler thermal hub type of the Vaillant owner.
type BoilerThermalHub interface {
	Declare(spec BoilerThermalAssetSpec) error
	Begin(handle BoilerThermalHandle, asset semreg.AssetID) (BoilerThermalTransaction, error)
	ReadAdmitted(asset semreg.AssetID) bool
	Snapshot(asset semreg.AssetID) (semreg.Snapshot, error)
}
