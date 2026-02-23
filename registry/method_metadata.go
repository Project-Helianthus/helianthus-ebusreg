package registry

// MethodMutability describes whether a method has side effects.
type MethodMutability string

const (
	MethodMutabilityUnknown  MethodMutability = "unknown"
	MethodMutabilityReadOnly MethodMutability = "read_only"
	MethodMutabilityMutating MethodMutability = "mutating"
)

// MethodDanger describes risk level for invoke safety controls.
type MethodDanger string

const (
	MethodDangerUnknown   MethodDanger = "unknown"
	MethodDangerSafe      MethodDanger = "safe"
	MethodDangerDangerous MethodDanger = "dangerous"
)

// MethodMetadata is the normalized method safety/routing contract.
//
// Backward compatibility defaults:
// - Mutability: derived from Method.ReadOnly() when not explicitly provided.
// - Danger: safe for read_only, dangerous for mutating/unknown.
// - Routable: true when not explicitly provided.
type MethodMetadata struct {
	Mutability MethodMutability
	Danger     MethodDanger
	Routable   bool
}

// MethodMutabilityProvider optionally overrides default mutability inference.
type MethodMutabilityProvider interface {
	Mutability() MethodMutability
}

// MethodDangerProvider optionally overrides default danger inference.
type MethodDangerProvider interface {
	Danger() MethodDanger
}

// MethodRoutableProvider optionally overrides default routability.
type MethodRoutableProvider interface {
	Routable() bool
}

// ResolveMethodMetadata returns normalized metadata for a method with
// backward-compatible defaults for legacy Method implementations.
func ResolveMethodMetadata(method Method) MethodMetadata {
	mutability := MethodMutabilityOf(method)
	return MethodMetadata{
		Mutability: mutability,
		Danger:     MethodDangerOf(method),
		Routable:   MethodRoutableOf(method),
	}
}

// MethodMutabilityOf resolves mutability with legacy ReadOnly fallback.
func MethodMutabilityOf(method Method) MethodMutability {
	if method == nil {
		return MethodMutabilityUnknown
	}
	if provider, ok := method.(MethodMutabilityProvider); ok {
		if mutability := normalizeMethodMutability(provider.Mutability()); mutability != MethodMutabilityUnknown {
			return mutability
		}
	}
	if method.ReadOnly() {
		return MethodMutabilityReadOnly
	}
	return MethodMutabilityMutating
}

// MethodDangerOf resolves danger with mutability-based fallback.
func MethodDangerOf(method Method) MethodDanger {
	if method == nil {
		return MethodDangerDangerous
	}
	if provider, ok := method.(MethodDangerProvider); ok {
		if danger := normalizeMethodDanger(provider.Danger()); danger != MethodDangerUnknown {
			return danger
		}
	}
	mutability := MethodMutabilityOf(method)
	if mutability == MethodMutabilityReadOnly {
		return MethodDangerSafe
	}
	return MethodDangerDangerous
}

// MethodRoutableOf resolves routability with legacy default=true behavior.
func MethodRoutableOf(method Method) bool {
	if method == nil {
		return false
	}
	if provider, ok := method.(MethodRoutableProvider); ok {
		return provider.Routable()
	}
	return true
}

func normalizeMethodMutability(mutability MethodMutability) MethodMutability {
	switch mutability {
	case MethodMutabilityReadOnly, MethodMutabilityMutating:
		return mutability
	default:
		return MethodMutabilityUnknown
	}
}

func normalizeMethodDanger(danger MethodDanger) MethodDanger {
	switch danger {
	case MethodDangerSafe, MethodDangerDangerous:
		return danger
	default:
		return MethodDangerUnknown
	}
}
