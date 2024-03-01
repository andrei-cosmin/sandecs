package filter

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/sandbox"
)

// Filter struct - collection of rules that can be used to filter entities.
//   - Rules []sandbox.Rule - a collection of rules
type Filter struct {
	Rules []sandbox.Rule
}

// createTagRules method - creates a filter with rules for a specific tag types
func createTagRules(ruleType sandbox.Type, tags ...component.Tag) Filter {
	// Create a rule for each tag
	rules := make([]sandbox.Rule, len(tags))
	for index, tag := range tags {
		// Create a new tag rule with the tag and rule type (match, exclude, union)
		rules[index] = sandbox.NewTagRule(tag, ruleType)
	}
	// Return the filter with the rules
	return Filter{
		Rules: rules,
	}
}

// MatchTags creates a filter that matches entities with the specified tags
func MatchTags(tags ...component.Tag) Filter {
	return createTagRules(sandbox.Match, tags...)
}

// ExcludeTags creates a filter that excludes entities with the specified tags
func ExcludeTags(tags ...component.Tag) Filter {
	return createTagRules(sandbox.Exclude, tags...)
}

// UnionTags creates a filter that includes entities with one of the specified tags
func UnionTags(tags ...component.Tag) Filter {
	return createTagRules(sandbox.Union, tags...)
}

// Match creates a filter that matches entities with the specified component
func Match[T component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[T](sandbox.Match),
		},
	}
}

// Match2 creates a filter that matches entities with the specified components
func Match2[A, B component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Match),
			sandbox.NewComponentRule[B](sandbox.Match),
		},
	}
}

// Match3 creates a filter that matches entities with the specified components
func Match3[A, B, C component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Match),
			sandbox.NewComponentRule[B](sandbox.Match),
			sandbox.NewComponentRule[C](sandbox.Match),
		},
	}
}

// Match4 creates a filter that matches entities with the specified components
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

// Match5 creates a filter that matches entities with the specified components
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

// Exclude creates a filter that excludes entities with the specified component
func Exclude[T component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[T](sandbox.Exclude),
		},
	}
}

// Exclude2 creates a filter that excludes entities with the specified components
func Exclude2[A, B component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Exclude),
			sandbox.NewComponentRule[B](sandbox.Exclude),
		},
	}
}

// Exclude3 creates a filter that excludes entities with the specified components
func Exclude3[A, B, C component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Exclude),
			sandbox.NewComponentRule[B](sandbox.Exclude),
			sandbox.NewComponentRule[C](sandbox.Exclude),
		},
	}
}

// Exclude4 creates a filter that excludes entities with the specified components
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

// Exclude5 creates a filter that excludes entities with the specified components
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

// Union creates a filter that includes entities with one the specified components
func Union[A component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Union),
		},
	}
}

// Union2 creates a filter that includes entities with one the specified components
func Union2[A, B component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Union),
			sandbox.NewComponentRule[B](sandbox.Union),
		},
	}
}

// Union3 creates a filter that includes entities with one the specified components
func Union3[A, B, C component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Union),
			sandbox.NewComponentRule[B](sandbox.Union),
			sandbox.NewComponentRule[C](sandbox.Union),
		},
	}
}

// Union4 creates a filter that includes entities with one the specified components
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

// Union5 creates a filter that includes entities with one the specified components
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
