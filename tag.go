package main

/* Tag is a basic type describing a tag to categorize items, either as simple statement or as an expression */


/* ================================================================================ Imports */
import (
	"fmt"
	"strings"
)


/* ================================================================================ Public types */
type Tag struct {
	Expression string
}


/* ================================================================================ Public functions */
func ParseTagEditString(tagString string) []Tag {
	if len(tagString) < 1 {
		return nil
	}

	expressions := strings.Split(tagString, ";")
	validCount  := 0

	for i, expression := range expressions {
		expressions[i] = strings.TrimSpace(expression)
		if len(expressions[i]) > 0 {
			validCount++
		}
	}

	tags := make([]Tag, validCount)
	i    := 0
	for _, expression := range expressions {
		if len(expression) > 0 {
			tags[i].Expression = expression
			i++
		}
	}

	return tags
}


func ComposeTagEditString(tags []Tag) string {
	buffer := &strings.Builder{}

	for _, tag := range tags {
		fmt.Fprintf(buffer, "%s; ", tag.Expression)
	}

	return buffer.String()
}


/* ================================================================================ Public methods */
func (t *Tag) DisplayString() string {
	before, after, found := strings.Cut(t.Expression, "=")

	if(found) {
		return fmt.Sprintf("%s: %s", before, after)
	} else {
		return before
	}
}