package types

type Blacklist map[string][]string

func (bl Blacklist) Empty() bool {
	return len(bl) < 1
}

func (bl Blacklist) Add(contract string, action string) {
	if len(bl[contract]) < 1 {
		bl[contract] = []string{}
	}
	bl[contract] = append(bl[contract], action)
}

func (bl Blacklist) IsAllowed(contract string, action string) bool {
	if v, ok := bl[contract]; ok {
		for _, act := range v {
			if act == action || act == "*" {
				return false
			}
		}
	}
	return true
}

func (bl Blacklist) IsDenied(contract string, action string) bool {
	return bl.IsAllowed(contract, action)
}
