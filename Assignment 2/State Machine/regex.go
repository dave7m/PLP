package main

import "regexp"

var (
	transitionRE, _  = regexp.Compile("(?m)^>\\s*\\w+\\s*\\([\\s*\\w]+\\)\\s*\\w+\\s*:[^\\n]*$")
	endStateRE, _    = regexp.Compile("(?m)^@\\s*\\+\\s*[a-zA-Z0-9]*\\s*{[^}]*}")
	normalStateRE, _ = regexp.Compile("(?m)^@\\s*[a-zA-Z0-9]*\\s*{[^}]*}")
	startStateRE, _  = regexp.Compile("(?m)^@\\s*\\*\\s*[a-zA-Z0-9]*\\s*{[^}]*}")
)
