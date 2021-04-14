package yaml2json

import "bytes"
import "fmt"
import "strings"
import "regexp"

type parsingState int

const (
	parseDirectives parsingState = iota
	parseContent
)

type lineType int

const (
	lineComment lineType = iota
	lineDirective
	lineDocStart
	lineDocEnd
	lineEmpty
	lineOther
)

type lineRegexps struct {
	comment   *regexp.Regexp
	directive *regexp.Regexp
	docStart  *regexp.Regexp
	docEnd    *regexp.Regexp
	empty     *regexp.Regexp
}

func detectLineType(line string, re *lineRegexps) lineType {
	if re.empty.MatchString(line) {
		return lineEmpty
	}
	if re.comment.MatchString(line) {
		return lineComment
	}
	if re.docStart.MatchString(line) {
		return lineDocStart
	}
	if re.docEnd.MatchString(line) {
		return lineDocEnd
	}
	if re.directive.MatchString(line) {
		return lineDirective
	}
	return lineOther
}

func makeLineRegexps() (*lineRegexps, error) {
	re := new(lineRegexps)

	var err error
	re.comment, err = regexp.Compile("^[\t ]*#")
	if err != nil {
		return nil, fmt.Errorf("Unable to compile comment regexp")
	}
	re.directive, err = regexp.Compile("^[\t ]*%")
	if err != nil {
		return nil, fmt.Errorf("Unable to compile directive regexp")
	}
	re.docStart, err = regexp.Compile("^---$")
	if err != nil {
		return nil, fmt.Errorf("Unable to compile docStart regexp")
	}
	re.docEnd, err = regexp.Compile("^...$")
	if err != nil {
		return nil, fmt.Errorf("Unable to compile docEnd regexp")
	}
	re.empty, err = regexp.Compile("^[\t ]*$")
	if err != nil {
		return nil, fmt.Errorf("Unable to compile empty regexp")
	}

	return re, nil
}

func splitYAML(content []byte) ([][]byte, error) {
	result := make([][]byte, 0)
	buf := bytes.NewBuffer(content)

	lines := strings.SplitAfter(buf.String(), "\n")

	re, err := makeLineRegexps()
	if err != nil {
		return nil, err
	}

	buf.Reset()

	i := 0
	curState := parseDirectives

	for _, line := range lines {
		if len(line) > 0 {
			if line[len(line)-1] == '\n' {
				line = line[0 : len(line)-1]
			}
		}

		lt := detectLineType(line, re)

		moveToNext := func() {
			curState = parseDirectives
			x := buf.Bytes()
			y := make([]byte, len(x))
			copy(y, x)
			result = append(result, y)
			buf.Reset()
			i++
		}

		switch lt {
		case lineComment:
			break

		case lineEmpty:
			break

		case lineDocEnd:
			moveToNext()
			continue

		case lineDocStart:
			if curState == parseDirectives {
				curState = parseContent
				break
			}
			moveToNext()
			break

		case lineDirective:
			break

		case lineOther:
			curState = parseContent
			break
		}

		buf.WriteString(line)
		buf.WriteString("\n")
	}
	if buf.Len() > 0 {
		result = append(result, buf.Bytes())
	}
	return result, nil
}
