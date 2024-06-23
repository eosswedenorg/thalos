package types

type Blacklist map[string][]string

func (bl Blacklist) Add(contract string, action string) {
	if len(bl[contract]) < 1 {
		bl[contract] = []string{}
	}
	bl[contract] = append(bl[contract], action)
}

func (bl Blacklist) Lookup(contract string, action string) bool {
	if v, ok := bl[contract]; ok {
		for _, act := range v {
			if act == action || act == "*" {
				return true
			}
		}
	}
	return false
}
