package registry

import (
	"testing"

	"github.com/d3vi1/helianthus-ebusreg/schema"
)

type metadataTestMethod struct {
	name       string
	readOnly   bool
	mutability MethodMutability
	danger     MethodDanger
	routable   bool
	template   FrameTemplate
}

func (m metadataTestMethod) Name() string {
	return m.name
}

func (m metadataTestMethod) ReadOnly() bool {
	return m.readOnly
}

func (m metadataTestMethod) Template() FrameTemplate {
	return m.template
}

func (m metadataTestMethod) ResponseSchema() schema.SchemaSelector {
	return schema.SchemaSelector{}
}

func (m metadataTestMethod) Mutability() MethodMutability {
	return m.mutability
}

func (m metadataTestMethod) Danger() MethodDanger {
	return m.danger
}

func (m metadataTestMethod) Routable() bool {
	return m.routable
}

type noMetadataMethod struct {
	name     string
	readOnly bool
	template FrameTemplate
}

func (m noMetadataMethod) Name() string {
	return m.name
}

func (m noMetadataMethod) ReadOnly() bool {
	return m.readOnly
}

func (m noMetadataMethod) Template() FrameTemplate {
	return m.template
}

func (m noMetadataMethod) ResponseSchema() schema.SchemaSelector {
	return schema.SchemaSelector{}
}

type mutabilityOnlyMethod struct {
	noMetadataMethod
	mutability MethodMutability
}

func (m mutabilityOnlyMethod) Mutability() MethodMutability {
	return m.mutability
}

type dangerOnlyMethod struct {
	noMetadataMethod
	danger MethodDanger
}

func (m dangerOnlyMethod) Danger() MethodDanger {
	return m.danger
}

type routableOnlyMethod struct {
	noMetadataMethod
	routable bool
}

func (m routableOnlyMethod) Routable() bool {
	return m.routable
}

func TestResolveMethodMetadata_DefaultsFromReadOnly(t *testing.T) {
	t.Parallel()

	readOnly := noMetadataMethod{
		name:     "get_status",
		readOnly: true,
		template: mockTemplate{primary: 0xB5, secondary: 0x04},
	}
	readOnlyMetadata := ResolveMethodMetadata(readOnly)
	if readOnlyMetadata.Mutability != MethodMutabilityReadOnly {
		t.Fatalf("read-only mutability = %q; want %q", readOnlyMetadata.Mutability, MethodMutabilityReadOnly)
	}
	if readOnlyMetadata.Danger != MethodDangerSafe {
		t.Fatalf("read-only danger = %q; want %q", readOnlyMetadata.Danger, MethodDangerSafe)
	}
	if !readOnlyMetadata.Routable {
		t.Fatalf("read-only routable = false; want true")
	}

	mutating := noMetadataMethod{
		name:     "set_status",
		readOnly: false,
		template: mockTemplate{primary: 0xB5, secondary: 0x05},
	}
	mutatingMetadata := ResolveMethodMetadata(mutating)
	if mutatingMetadata.Mutability != MethodMutabilityMutating {
		t.Fatalf("mutating mutability = %q; want %q", mutatingMetadata.Mutability, MethodMutabilityMutating)
	}
	if mutatingMetadata.Danger != MethodDangerDangerous {
		t.Fatalf("mutating danger = %q; want %q", mutatingMetadata.Danger, MethodDangerDangerous)
	}
	if !mutatingMetadata.Routable {
		t.Fatalf("mutating routable = false; want true")
	}
}

func TestResolveMethodMetadata_ExplicitOverrides(t *testing.T) {
	t.Parallel()

	method := metadataTestMethod{
		name:       "set_register",
		readOnly:   false,
		mutability: MethodMutabilityReadOnly,
		danger:     MethodDangerSafe,
		routable:   false,
		template:   mockTemplate{primary: 0xB5, secondary: 0x09},
	}

	metadata := ResolveMethodMetadata(method)
	if metadata.Mutability != MethodMutabilityReadOnly {
		t.Fatalf("mutability = %q; want %q", metadata.Mutability, MethodMutabilityReadOnly)
	}
	if metadata.Danger != MethodDangerSafe {
		t.Fatalf("danger = %q; want %q", metadata.Danger, MethodDangerSafe)
	}
	if metadata.Routable {
		t.Fatalf("routable = true; want false")
	}
}

func TestResolveMethodMetadata_ExplicitUnknownMutabilityIsPreserved(t *testing.T) {
	t.Parallel()

	method := metadataTestMethod{
		name:       "set_register",
		readOnly:   true,
		mutability: MethodMutabilityUnknown,
		danger:     MethodDangerDangerous,
		routable:   true,
		template:   mockTemplate{primary: 0xB5, secondary: 0x09},
	}

	metadata := ResolveMethodMetadata(method)
	if metadata.Mutability != MethodMutabilityUnknown {
		t.Fatalf("mutability = %q; want %q", metadata.Mutability, MethodMutabilityUnknown)
	}
	if metadata.Danger != MethodDangerDangerous {
		t.Fatalf("danger = %q; want %q", metadata.Danger, MethodDangerDangerous)
	}
}

func TestResolveMethodMetadata_InvalidOverrideFallsBack(t *testing.T) {
	t.Parallel()

	method := metadataTestMethod{
		name:       "set_register",
		readOnly:   false,
		mutability: MethodMutability("invalid"),
		danger:     MethodDanger("invalid"),
		routable:   true,
		template:   mockTemplate{primary: 0xB5, secondary: 0x09},
	}

	metadata := ResolveMethodMetadata(method)
	if metadata.Mutability != MethodMutabilityMutating {
		t.Fatalf("mutability = %q; want %q", metadata.Mutability, MethodMutabilityMutating)
	}
	if metadata.Danger != MethodDangerDangerous {
		t.Fatalf("danger = %q; want %q", metadata.Danger, MethodDangerDangerous)
	}
	if !metadata.Routable {
		t.Fatalf("routable = false; want true")
	}
}

func TestResolveMethodMetadata_NilMethodDefaults(t *testing.T) {
	t.Parallel()

	metadata := ResolveMethodMetadata(nil)
	if metadata.Mutability != MethodMutabilityUnknown {
		t.Fatalf("mutability = %q; want %q", metadata.Mutability, MethodMutabilityUnknown)
	}
	if metadata.Danger != MethodDangerDangerous {
		t.Fatalf("danger = %q; want %q", metadata.Danger, MethodDangerDangerous)
	}
	if metadata.Routable {
		t.Fatalf("routable = true; want false")
	}
}

func TestResolveMethodMetadata_PartialProviders(t *testing.T) {
	t.Parallel()

	t.Run("mutability provider influences danger fallback", func(t *testing.T) {
		t.Parallel()

		method := mutabilityOnlyMethod{
			noMetadataMethod: noMetadataMethod{
				name:     "custom",
				readOnly: false,
				template: mockTemplate{primary: 0xB5, secondary: 0x04},
			},
			mutability: MethodMutabilityReadOnly,
		}
		metadata := ResolveMethodMetadata(method)
		if metadata.Mutability != MethodMutabilityReadOnly {
			t.Fatalf("mutability = %q; want %q", metadata.Mutability, MethodMutabilityReadOnly)
		}
		if metadata.Danger != MethodDangerSafe {
			t.Fatalf("danger = %q; want %q", metadata.Danger, MethodDangerSafe)
		}
		if !metadata.Routable {
			t.Fatalf("routable = false; want true")
		}
	})

	t.Run("unknown mutability remains dangerous without explicit danger", func(t *testing.T) {
		t.Parallel()

		method := mutabilityOnlyMethod{
			noMetadataMethod: noMetadataMethod{
				name:     "custom_unknown",
				readOnly: true,
				template: mockTemplate{primary: 0xB5, secondary: 0x04},
			},
			mutability: MethodMutabilityUnknown,
		}
		metadata := ResolveMethodMetadata(method)
		if metadata.Mutability != MethodMutabilityUnknown {
			t.Fatalf("mutability = %q; want %q", metadata.Mutability, MethodMutabilityUnknown)
		}
		if metadata.Danger != MethodDangerDangerous {
			t.Fatalf("danger = %q; want %q", metadata.Danger, MethodDangerDangerous)
		}
	})

	t.Run("danger provider override respected", func(t *testing.T) {
		t.Parallel()

		method := dangerOnlyMethod{
			noMetadataMethod: noMetadataMethod{
				name:     "danger_override",
				readOnly: false,
				template: mockTemplate{primary: 0xB5, secondary: 0x05},
			},
			danger: MethodDangerSafe,
		}
		metadata := ResolveMethodMetadata(method)
		if metadata.Mutability != MethodMutabilityMutating {
			t.Fatalf("mutability = %q; want %q", metadata.Mutability, MethodMutabilityMutating)
		}
		if metadata.Danger != MethodDangerSafe {
			t.Fatalf("danger = %q; want %q", metadata.Danger, MethodDangerSafe)
		}
	})

	t.Run("unknown danger falls back to derived value", func(t *testing.T) {
		t.Parallel()

		method := dangerOnlyMethod{
			noMetadataMethod: noMetadataMethod{
				name:     "danger_unknown",
				readOnly: true,
				template: mockTemplate{primary: 0xB5, secondary: 0x04},
			},
			danger: MethodDangerUnknown,
		}
		metadata := ResolveMethodMetadata(method)
		if metadata.Danger != MethodDangerSafe {
			t.Fatalf("danger = %q; want %q", metadata.Danger, MethodDangerSafe)
		}
	})

	t.Run("routable provider only affects routable", func(t *testing.T) {
		t.Parallel()

		method := routableOnlyMethod{
			noMetadataMethod: noMetadataMethod{
				name:     "non_routable",
				readOnly: true,
				template: mockTemplate{primary: 0xB5, secondary: 0x04},
			},
			routable: false,
		}
		metadata := ResolveMethodMetadata(method)
		if metadata.Mutability != MethodMutabilityReadOnly {
			t.Fatalf("mutability = %q; want %q", metadata.Mutability, MethodMutabilityReadOnly)
		}
		if metadata.Danger != MethodDangerSafe {
			t.Fatalf("danger = %q; want %q", metadata.Danger, MethodDangerSafe)
		}
		if metadata.Routable {
			t.Fatalf("routable = true; want false")
		}
	})
}
