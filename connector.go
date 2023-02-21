package rrors

import "strings"

const FirstConnectorMultiple = "┳"
const StraightLine = "━"
const Wall = "┃"
const MiddleConnector = "┣"
const LastConnector = "┗"

func getLine(pos int, len int, multipleBelow bool) string {
	if pos == 0 {
		if multipleBelow {
			return FirstConnectorMultiple
		}

		if len == 1 {
			return ">"
		}

		return StraightLine
	}

	if pos == len-1 {
		return LastConnector
	}

	return MiddleConnector
}

func getWall(multipleBelow bool) string {
	if multipleBelow {
		return Wall
	}

	return " "
}

func ConnectString(tag string, s []string) string {
	for i := range s {
		s[i] = " " + strings.Replace(s[i], "\n", "\n ", 1000)
	}
	tag += " "
	builder := strings.Builder{}
	spaces := strings.Repeat(" ", len(tag))

	for i, element := range s {
		toPrintFirst := spaces
		if i == 0 {
			toPrintFirst = tag
		}

		l := getLine(i, len(s), len(s) > 1)
		builder.WriteString(toPrintFirst + l)
		values := strings.Split(element, "\n")
		if len(values) > 0 {
			builder.WriteString(values[0] + "\n")
			wall := getWall(i != len(s)-1)
			for _, v := range values[1:] {
				builder.WriteString(spaces + wall + v + "\n")
			}
		}
	}

	return strings.TrimSpace(builder.String())
}
