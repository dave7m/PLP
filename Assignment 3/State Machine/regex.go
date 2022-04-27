package main

import "regexp"

var (
	transitionRE, _  = regexp.Compile("(?m)^>\\s*\\w+\\s*\\(.*\\)\\s*\\w+\\s*:[^\\n]*$")
	endStateRE, _    = regexp.Compile("(?m)^@\\s*((!\\+|\\+!)\\s*[a-zA-Z]*\\s*[0-9]+|\\+\\s*[a-zA-Z]*)\\s*{[^}]*}")
	normalStateRE, _ = regexp.Compile("(?m)^@\\s*(!\\s*[a-zA-Z]*\\s*[0-9]+|[a-zA-Z]*)\\s*{[^}]*}")
	startStateRE, _  = regexp.Compile("(?m)^@\\s*((!\\*|\\*!)\\s*[a-zA-Z]*\\s*[0-9]+|\\*\\s*[a-zA-Z]*)\\s*{[^}]*}")
)
