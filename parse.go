package main

import (
	"encoding/json"
	"regexp"
	"strings"
)

const teamString = `steelix (
    uses []
    url "0.0.0.0:8081"
    endpoints (
        /jwtkeypub (
            methods [GET]
        )
        /register (
            methods [POST]
        )
        /logout (
            methods [POST]
        )
        /login (
            methods [POST]
        )
        /apitokens (
            methods [GET]
        )
        /accesstokens (
            methods [GET]
        )
    )
    jwtInfo (
        issuerName "steelix"
        audienceName "steelix"
    )
)

klefki (
    uses [steelix]
    url "0.0.0.0:8083"
    endpoints (
        / (
            methods [GET, PATCH, DELETE]
        )
    )
    jwtInfo (
        audienceName "klefki"
    )
)

alakazam (
    uses [steelix]
    url "0.0.0.0:8082"
    endpoints (
        /ws (
            methods [GET]
        )
    )
    jwtInfo (
        audienceName "alakazam"
    )
)`

func ParseTeamfileString(contents string) Team {
	t := make(Team)

	re := regexp.MustCompile(`[\s\n\t;]`)
	contents = re.ReplaceAllString(contents, "")
	services := parseValue(contents)

	uses := make(map[string][]string)

	for servName, info := range services {
		info := info.(map[string]any)
		info["name"] = servName

		n := servName

		t[n] = NewServiceNode(servName)
		uses[n] = info["uses"].([]string)
		delete(info, "uses")

		// TODO this not not properly parse endpoints xausing nil opinter deref error
		b, _ := json.Marshal(info)
		json.Unmarshal(b, t[n])
	}

	for _, service := range t {
		for _, name := range uses[service.Name] {
			service.Uses(t[name])
		}
	}

	return t
}

func parseValue(contents string) map[string]any {
	value := make(map[string]any)
	attrs := splitAttrs(contents)

	for k, val := range attrs {
		var v any

		if val[0] == '[' {
			v = strings.Split(val[1:len(val)-1], ",")
		} else if val[0] == '(' {
			v = parseValue(val[1 : len(val)-1])
		} else if val[0] == '"' {
			v = val[1 : len(val)-1]
		}

		value[k] = v
	}

	return value
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

		if valEnd == len(buf) - 1 {
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

	var startChar byte
	var endChar byte

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

		if semaphore == 0 {
			break
		}
	}

	end = i
	return start, end
}
