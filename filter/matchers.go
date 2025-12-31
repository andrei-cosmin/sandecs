package filter

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/sandbox"
)

// Filter is a collection of rules used to query entities.
type Filter struct {
	Rules []sandbox.Rule
}

func createTagRules(ruleType sandbox.Type, tags ...component.Tag) Filter {
	rules := make([]sandbox.Rule, len(tags))
	for index, tag := range tags {
		rules[index] = sandbox.NewTagRule(tag, ruleType)
	}
	return Filter{Rules: rules}
}

// MatchTags matches entities with all specified tags.
func MatchTags(tags ...component.Tag) Filter {
	return createTagRules(sandbox.Match, tags...)
}

// ExcludeTags excludes entities with any of the specified tags.
func ExcludeTags(tags ...component.Tag) Filter {
	return createTagRules(sandbox.Exclude, tags...)
}

// UnionTags matches entities with at least one of the specified tags.
func UnionTags(tags ...component.Tag) Filter {
	return createTagRules(sandbox.Union, tags...)
}

// Match matches entities with component T.
func Match[T component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[T](sandbox.Match),
		},
	}
}

// Match2 matches entities with components A and B.
func Match2[A, B component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Match),
			sandbox.NewComponentRule[B](sandbox.Match),
		},
	}
}

// Match3 matches entities with components A, B, and C.
func Match3[A, B, C component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Match),
			sandbox.NewComponentRule[B](sandbox.Match),
			sandbox.NewComponentRule[C](sandbox.Match),
		},
	}
}

// Match4 matches entities with components A, B, C, and D.
func Match4[A, B, C, D component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Match),
			sandbox.NewComponentRule[B](sandbox.Match),
			sandbox.NewComponentRule[C](sandbox.Match),
			sandbox.NewComponentRule[D](sandbox.Match),
		},
	}
}

// Match5 matches entities with components A, B, C, D, and E.
func Match5[A, B, C, D, E component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Match),
			sandbox.NewComponentRule[B](sandbox.Match),
			sandbox.NewComponentRule[C](sandbox.Match),
			sandbox.NewComponentRule[D](sandbox.Match),
			sandbox.NewComponentRule[E](sandbox.Match),
		},
	}
}

// Exclude excludes entities with component T.
func Exclude[T component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[T](sandbox.Exclude),
		},
	}
}

// Exclude2 excludes entities with components A or B.
func Exclude2[A, B component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Exclude),
			sandbox.NewComponentRule[B](sandbox.Exclude),
		},
	}
}

// Exclude3 excludes entities with components A, B, or C.
func Exclude3[A, B, C component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Exclude),
			sandbox.NewComponentRule[B](sandbox.Exclude),
			sandbox.NewComponentRule[C](sandbox.Exclude),
		},
	}
}

// Exclude4 excludes entities with components A, B, C, or D.
func Exclude4[A, B, C, D component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Exclude),
			sandbox.NewComponentRule[B](sandbox.Exclude),
			sandbox.NewComponentRule[C](sandbox.Exclude),
			sandbox.NewComponentRule[D](sandbox.Exclude),
		},
	}
}

// Exclude5 excludes entities with components A, B, C, D, or E.
func Exclude5[A, B, C, D, E component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Exclude),
			sandbox.NewComponentRule[B](sandbox.Exclude),
			sandbox.NewComponentRule[C](sandbox.Exclude),
			sandbox.NewComponentRule[D](sandbox.Exclude),
			sandbox.NewComponentRule[E](sandbox.Exclude),
		},
	}
}

// Union matches entities with component A.
func Union[A component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Union),
		},
	}
}

// Union2 matches entities with at least one of components A or B.
func Union2[A, B component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Union),
			sandbox.NewComponentRule[B](sandbox.Union),
		},
	}
}

// Union3 matches entities with at least one of components A, B, or C.
func Union3[A, B, C component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Union),
			sandbox.NewComponentRule[B](sandbox.Union),
			sandbox.NewComponentRule[C](sandbox.Union),
		},
	}
}

// Union4 matches entities with at least one of components A, B, C, or D.
func Union4[A, B, C, D component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Union),
			sandbox.NewComponentRule[B](sandbox.Union),
			sandbox.NewComponentRule[C](sandbox.Union),
			sandbox.NewComponentRule[D](sandbox.Union),
		},
	}
}

// Union5 matches entities with at least one of components A, B, C, D, or E.
func Union5[A, B, C, D, E component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Union),
			sandbox.NewComponentRule[B](sandbox.Union),
			sandbox.NewComponentRule[C](sandbox.Union),
			sandbox.NewComponentRule[D](sandbox.Union),
			sandbox.NewComponentRule[E](sandbox.Union),
		},
	}
}
