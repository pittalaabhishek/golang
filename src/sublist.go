package sublist

// isSublist checks if list a is contained within list b
func isSublist(a, b []int) bool {
	if len(a) == 0 {
		return true
	}
	if len(a) > len(b) {
		return false
	}

	for i := 0; i <= len(b)-len(a); i++ {
		match := true
		for j := 0; j < len(a); j++ {
			if a[j] != b[i+j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// Sublist determines the relationship between two lists
func Sublist(l1, l2 []int) Relation {
	switch {
	case len(l1) == len(l2) && isSublist(l1, l2):
		return RelationEqual
	case len(l1) < len(l2) && isSublist(l1, l2):
		return RelationSublist
	case len(l1) > len(l2) && isSublist(l2, l1):
		return RelationSuperlist
	default:
		return RelationUnequal
	}
}