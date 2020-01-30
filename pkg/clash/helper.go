package goclash

import (
	"errors"
	"strings"
)

var (
	errBeforeAfterSet  = errors.New("both Before and After have been set")
	errInvalidOptional = errors.New("could not encode optional arguments")
	errInvalidTag      = errors.New("tag was not valid")
)

func validateTag(tag string) error {
	if tag[:1] != "#" {
		return errInvalidTag
	}
	return nil
}

func buildURLPath(str ...string) string {
	var path strings.Builder
	for _, s := range str {
		path.WriteString(s)
	}
	return path.String()
}
