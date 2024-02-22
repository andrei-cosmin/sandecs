package filter

import (
	"github.com/andrei-cosmin/hakkt/component"
	"github.com/andrei-cosmin/hakkt/internal/sandbox"
)

type Filter struct {
	Rules []sandbox.Rule
}

func Match[T component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[T](sandbox.Match),
		},
	}
}

func Match2[A, B component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[A](sandbox.Match),
			sandbox.NewRule[B](sandbox.Match),
		},
	}
}

func Match3[A, B, C component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[A](sandbox.Match),
			sandbox.NewRule[B](sandbox.Match),
			sandbox.NewRule[C](sandbox.Match),
		},
	}
}

func Match4[A, B, C, D component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[A](sandbox.Match),
			sandbox.NewRule[B](sandbox.Match),
			sandbox.NewRule[C](sandbox.Match),
			sandbox.NewRule[D](sandbox.Match),
		},
	}
}

func Match5[A, B, C, D, E component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[A](sandbox.Match),
			sandbox.NewRule[B](sandbox.Match),
			sandbox.NewRule[C](sandbox.Match),
			sandbox.NewRule[D](sandbox.Match),
			sandbox.NewRule[E](sandbox.Match),
		},
	}
}

func Exclude[T component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[T](sandbox.Exclude),
		},
	}
}

func Exclude2[A, B component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[A](sandbox.Exclude),
			sandbox.NewRule[B](sandbox.Exclude),
		},
	}
}

func Exclude3[A, B, C component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[A](sandbox.Exclude),
			sandbox.NewRule[B](sandbox.Exclude),
			sandbox.NewRule[C](sandbox.Exclude),
		},
	}
}

func Exclude4[A, B, C, D component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[A](sandbox.Exclude),
			sandbox.NewRule[B](sandbox.Exclude),
			sandbox.NewRule[C](sandbox.Exclude),
			sandbox.NewRule[D](sandbox.Exclude),
		},
	}
}

func Exclude5[A, B, C, D, E component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[A](sandbox.Exclude),
			sandbox.NewRule[B](sandbox.Exclude),
			sandbox.NewRule[C](sandbox.Exclude),
			sandbox.NewRule[D](sandbox.Exclude),
			sandbox.NewRule[E](sandbox.Exclude),
		},
	}
}

func Union2[A, B component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[A](sandbox.Union),
			sandbox.NewRule[B](sandbox.Union),
		},
	}
}

func Union3[A, B, C component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[A](sandbox.Union),
			sandbox.NewRule[B](sandbox.Union),
			sandbox.NewRule[C](sandbox.Union),
		},
	}
}

func Union4[A, B, C, D component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[A](sandbox.Union),
			sandbox.NewRule[B](sandbox.Union),
			sandbox.NewRule[C](sandbox.Union),
			sandbox.NewRule[D](sandbox.Union),
		},
	}
}

func Union5[A, B, C, D, E component.Component]() Filter {
	return Filter{
		Rules: []sandbox.Rule{
			sandbox.NewRule[A](sandbox.Union),
			sandbox.NewRule[B](sandbox.Union),
			sandbox.NewRule[C](sandbox.Union),
			sandbox.NewRule[D](sandbox.Union),
			sandbox.NewRule[E](sandbox.Union),
		},
	}
}
