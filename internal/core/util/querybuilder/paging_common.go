package querybuilder

type sortOrder string

func sortOrderFromString(sortOrder string) sortOrder {
	switch sortOrder {
	case "asc":
		return sortOrderASC
	case "desc":
		return sortOrderDESC
	default:
		return sortOrderASC
	}
}

const (
	// sortOrderASC ascending order
	sortOrderASC sortOrder = "asc"
	// sortOrderDESC descending order
	sortOrderDESC sortOrder = "desc"
)
