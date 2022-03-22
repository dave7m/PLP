package main

import "regexp"

var (
	transitionRE, _  = regexp.Compile("(?m)^>\\s*\\w+\\s*\\([\\s*\\w]+\\)\\s*\\w+\\s*:[^\\n]*$")
	endStateRE, _    = regexp.Compile("(?m)^@\\+[a-zA-Z0-9]*{[^}]*}")
	normalStateRE, _ = regexp.Compile("(?m)^@[a-zA-Z0-9]*{[^}]*}")
	startStateRE, _  = regexp.Compile("(?m)^@\\*[a-zA-Z0-9]*{[^}]*}")
)
