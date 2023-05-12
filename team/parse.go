package team

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"
)

func fromTeamFile(fp string) Team {
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}

	raw := string(b)
	re := regexp.MustCompile(`[\s\n\t;]`)
	raw = re.ReplaceAllString(raw, "")
	
	uses := make(map[string][]string)
	t := make(Team)

	for servName, info := range parseValue(raw) {
		info := info.(map[string]any)
		info["name"] = servName

		n := servName

		t[n] = &Service{}
		uses[n] = info["uses"].([]string)
		delete(info, "uses")

		b, _ := json.Marshal(info)
		json.Unmarshal(b, t[n])
	}

	for _, service := range t {
		for _, name := range uses[service.Name] {
			service.uses(t[name])
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
