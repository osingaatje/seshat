package multiplicity

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	multiplicity          *regexp.Regexp = regexp.MustCompile(`^(\d+|\*)(\.\.(\d+|\*))?$`)
	singleMatchCount                     = 2 // first one is always the entire string
	rangeMatchCount                      = 4
	fromMultiplicityIndex                = 1
	toMultiplicityIndex                  = 3
)

func GetMultiplicity(in string) (*Multiplicity, bool) {
	groups := multiplicity.FindStringSubmatch(in)

	if len(groups) < singleMatchCount || groups[fromMultiplicityIndex] == "" {
		return nil, false
	}

	if groups[toMultiplicityIndex] == "" { // single match
		fromMult := groups[fromMultiplicityIndex]

		fromMultInt, err := stringToMult(fromMult)
		if err != nil {
			return nil, false
		}

		return &Multiplicity{Start: fromMultInt, HasEndMult: false}, true // only a start index
	}
	fromMult := groups[fromMultiplicityIndex]
	toMult := groups[toMultiplicityIndex]

	fromMultInt, errFrom := stringToMult(fromMult)
	toMultInt, errTo := stringToMult(toMult)

	if errFrom != nil || errTo != nil { // not integers
		return nil, false
	}

	if fromMultInt > toMultInt { // swap around range if inverted
		temp := fromMultInt
		fromMultInt = toMultInt
		toMultInt = temp
	}

	if fromMultInt == toMultInt { // squash 1..1 for example into just 1
		return &Multiplicity{Start: fromMultInt, HasEndMult: false}, true
	}

	return &Multiplicity{
		Start:      fromMultInt,
		End:        toMultInt,
		HasEndMult: true,
	}, true
}

type Multiplicity struct {
	Start      int // 0, 1, ..., * == -1
	End        int // 1, 2, ..., * == -1
	HasEndMult bool
}

func (m Multiplicity) String() string {
	if !m.HasEndMult {
		return multToString(m.Start)
	}
	return multToString(m.Start) + ".." + multToString(m.End)
}

func stringToMult(in string) (int, error) {
	if strings.TrimSpace(in) == "*" {
		return -1, nil
	}
	intVal, err := strconv.Atoi(in)
	return intVal, err
}

func multToString(mult int) string {
	if mult < 0 {
		return "*"
	}
	return strconv.Itoa(mult)
}
