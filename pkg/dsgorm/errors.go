package dsgorm

import "regexp"

var mySQLErrors = map[string]*regexp.Regexp{
	"duplicate entry": regexp.MustCompile("Error 1062:"),
}

var gormErrors = map[string]*regexp.Regexp{
	"user does not exist": regexp.MustCompile("user does not exist"),
}
