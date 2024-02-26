package filter

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/sandbox"
)

type Filter struct {
	Rules []sandbox.Rule
}

func createTagRules(ruleType sandbox.Type, tags ...component.Tag) Filter {
	rules := make([]sandbox.Rule, len(tags))
	for index, tag := range tags {
		rules[index] = sandbox.NewTagRule(tag, ruleType)
	}
	return Filter{
		Rules: rules,
	}
}

func MatchTags(tags ...component.Tag) Filter {
	return createTagRules(sandbox.Match, tags...)
}

func ExcludeTags(tags ...component.Tag) Filter {
	return createTagRules(sandbox.Exclude, tags...)
}

func UnionTags(tags ...component.Tag) Filter {
	return createTagRules(sandbox.Union, tags...)
}

func Match[T component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[T](sandbox.Match),
		},
	}
}

func Match2[A, B component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Match),
			sandbox.NewComponentRule[B](sandbox.Match),
		},
	}
}

func Match3[A, B, C component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Match),
			sandbox.NewComponentRule[B](sandbox.Match),
			sandbox.NewComponentRule[C](sandbox.Match),
		},
	}
}

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

func Exclude[T component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[T](sandbox.Exclude),
		},
	}
}

func Exclude2[A, B component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Exclude),
			sandbox.NewComponentRule[B](sandbox.Exclude),
		},
	}
}

func Exclude3[A, B, C component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Exclude),
			sandbox.NewComponentRule[B](sandbox.Exclude),
			sandbox.NewComponentRule[C](sandbox.Exclude),
		},
	}
}

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

func Union[A component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Union),
		},
	}
}

func Union2[A, B component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Union),
			sandbox.NewComponentRule[B](sandbox.Union),
		},
	}
}

func Union3[A, B, C component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewComponentRule[A](sandbox.Union),
			sandbox.NewComponentRule[B](sandbox.Union),
			sandbox.NewComponentRule[C](sandbox.Union),
		},
	}
}

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
