package filter

const wildcard = "*"

type Mode uint8

const (
	Include Mode = iota // whitelist: match = included, default = excluded
	Exclude             // blacklist: match = excluded, default = included
)

type List struct {
	table map[string][]string
	mode  Mode
}

func New(entries map[string][]string) *List {
	return &List{
		table: entries,
		mode:  Include,
	}
}

func (l *List) SetMode(mode Mode) *List {
	l.mode = mode
	return l
}

func (l List) Mode() Mode {
	return l.mode
}

func (l List) Empty() bool {
	return len(l.table) < 1
}

func (l *List) Add(contract string, action string) *List {
	if l.table == nil {
		l.table = map[string][]string{}
	}
	l.table[contract] = append(l.table[contract], action)
	return l
}

func (l List) matches(contracts ...string) [][]string {
	ret := [][]string{}
	for _, contract := range contracts {
		if v, ok := l.table[contract]; ok {
			ret = append(ret, v)
		}
	}
	return ret
}

func (l List) In(contract string, action string) bool {
	for _, v := range l.matches(contract, wildcard) {
		for _, act := range v {
			if act == action || act == wildcard {
				return true
			}
		}
	}
	return false
}

func (l List) IsIncluded(contract string, action string) bool {
	res := l.In(contract, action)
	if l.mode == Exclude {
		res = !res
	}
	return res
}

func (l List) IsExcluded(contract string, action string) bool {
	return !l.IsIncluded(contract, action)
}
