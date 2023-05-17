package team

import (
	"fmt"
	"strings"
)

func parseObject(contents string) map[string]any {
	values := make(map[string]any)
	attrs := splitAttrs(contents)

	for k, v := range attrs {
		var val any
		startChar := v[0]
		v = v[1 : len(v)-1]

		if startChar == '[' { // string array
			val = strings.Split(v, ",")
			if val.([]string)[0] == "" { // TODO uhh what?
				val = []string{}
			}
		} else if startChar == '(' { // object
			val = parseObject(v)
		} else if startChar == '"' { // string
			val = v
		}

		values[k] = val
	}

	return values
}

func splitAttrs(contents string) map[string]string {
	attrs := make(map[string]string)
	buf := []byte(contents)

	var nameBuilder strings.Builder
	var valueBuilder strings.Builder

	for {
		valStart, valEnd := boundNextValue(buf)

		nameBuilder.Write(buf[:valStart])
		valueBuilder.Write(buf[valStart : valEnd+1])
		attrs[nameBuilder.String()] = valueBuilder.String()

		nameBuilder.Reset()
		valueBuilder.Reset()

		if valEnd == len(buf)-1 {
			break
		}

		buf = buf[valEnd+1:]
	}

	return attrs
}

func boundNextValue(contentBuf []byte) (start, end int) {
	endChars := map[byte]byte{
		'[': ']',
		'(': ')',
		'"': '"',
	}

	var (
		startChar byte
		endChar   byte
	)

	// find start
	for i := 0; i < len(contentBuf); i++ {
		if strings.ContainsAny(string(contentBuf[i]), `[("`) {
			start = i

			startChar = contentBuf[i]
			endChar = endChars[startChar]

			break
		}
	}

	// find end
	var semaphore int
	var readingString bool
	i := start
	for ; i < len(contentBuf); i++ {
		char := contentBuf[i]

		if !readingString && startChar == '"' && char == '"' {
			readingString = true
			semaphore++
			continue
		}

		if readingString && char == endChar {
			break
		}

		if char == startChar {
			semaphore++
		} else if char == endChar {
			semaphore--
		}

		if semaphore < 0 {
			panic(fmt.Errorf("unbalanced %q or %q in Teamfile", startChar, endChar))
		}

		if semaphore == 0 {
			break
		}
	}

	end = i
	return start, end
}
