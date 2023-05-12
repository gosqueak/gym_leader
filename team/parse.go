package team

import (
	"encoding/json"
	"os"
	"io"
	"regexp"
	"strings"
)

func fromTeamFile(fp string) Team {
	f, err := os.Open(fp)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	// get string and remove all whitespace and semicolons
	raw := string(b)
	re := regexp.MustCompile(`[\s\n\t;]`)
	raw = re.ReplaceAllString(raw, "")

	// init team
	team := make(Team)

	for name, info := range parseValue(raw) {

		// type assert info and add service name to info
		info := info.(map[string]any)
		info["name"] = name

		team[name] = &Service{}

		// marshall map to JSON then unmarshal into the new Service
		b, _ := json.Marshal(info)
		json.Unmarshal(b, team[name])
	}

	for _, service := range team {
		for _, name := range service.Dependencies {
			other := team[name]
			other.usedBy(service)
		}
	}

	return team
}

func parseValue(contents string) map[string]any {
	value := make(map[string]any)
	attrs := splitAttrs(contents)

	for k, val := range attrs {
		var v any

		if val[0] == '[' {
			v = strings.Split(val[1:len(val)-1], ",")
			if v.([]string)[0] == "" {
				v = []string{}
			}
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
