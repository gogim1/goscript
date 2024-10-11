package runtime

import "unicode"

func filterLexical(env []envItem) []envItem {
	newEnv := []envItem{}
	for _, item := range env {
		if len(item.name) > 0 && unicode.IsLower([]rune(item.name)[0]) {
			newEnv = append(newEnv, item)
		}
	}
	return newEnv
}

func lookupEnv(name string, env []envItem) int {
	for i := len(env) - 1; i >= 0; i-- {
		if env[i].name == name {
			return env[i].location
		}
	}
	return -1
}

func lookupStack(name string, layers []*layer) int {
	for i := len(layers) - 1; i >= 0; i-- {
		l := layers[i]
		if l.frame {
			env := *l.env
			for j := len(env) - 1; j >= 0; j-- {
				if env[j].name == name {
					return env[j].location
				}
			}
		}
	}
	return -1
}

func deepcopy(dst *[]*layer, src []*layer) {
	*dst = make([]*layer, len(src))
	for i, l := range src {
		(*dst)[i] = &layer{
			frame:  l.frame,
			expr:   l.expr,
			pc:     l.pc,
			callee: l.callee,
		}
		env := make([]envItem, len(*l.env))
		copy(env, *l.env)
		args := make([]Value, len(l.args))
		copy(args, l.args)

		(*dst)[i].env = &env
		(*dst)[i].args = args
	}
}
