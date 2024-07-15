package types

type Blacklist struct {
	table       map[string][]string
	isWhitelist bool
}

func NewBlacklist(entries map[string][]string) *Blacklist {
	return &Blacklist{
		table: entries,
	}
}

func (bl *Blacklist) SetWhitelist(value bool) *Blacklist {
	bl.isWhitelist = value
	return bl
}

func (bl Blacklist) Empty() bool {
	return len(bl.table) < 1
}

func (bl *Blacklist) Add(contract string, action string) {
	if bl.table == nil {
		bl.table = map[string][]string{}
	}

	if len(bl.table[contract]) < 1 {
		bl.table[contract] = []string{}
	}
	bl.table[contract] = append(bl.table[contract], action)
}

func (bl Blacklist) IsAllowed(contract string, action string) bool {
	if v, ok := bl.table[contract]; ok {
		for _, act := range v {
			if act == action || act == "*" {
				return bl.isWhitelist == true
			}
		}
	}
	return bl.isWhitelist == false
}

func (bl Blacklist) IsDenied(contract string, action string) bool {
	return bl.IsAllowed(contract, action)
}
