package router

var (
	table = new(TransTable)
)

func init() {
	table.init()
}

func AddItem(to uint64, next uint64, delay100us uint16) {
	if to > 0 && next > 0 {
		table.add(&transItem{to, next, delay100us})
	}
}

func DelItem(to uint64) {
	if to > 0 {
		table.del(to)
	}
}

func NextHop(to uint64) uint64 {
	if to > 0 {
		return table.nextHop(to).next
	}
	return 0
}

func NearMe() []uint64 {
	return table.near()
}
