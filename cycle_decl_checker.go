package anyhash

type visitStatus int

const (
	notVisiting = iota
	inVisiting
	visited
)

type cycleDeclChecker struct {
	typEdge   map[string]map[string]struct{}
	typVisits map[string]visitStatus
}

func (c cycleDeclChecker) isCycle(v string, resetVisits bool) bool {
	if resetVisits {
		for k := range c.typVisits {
			c.typVisits[k] = notVisiting
		}
	}

	c.typVisits[v] = inVisiting
	for k := range c.typEdge[v] {
		switch c.typVisits[k] {
		case inVisiting:
			return true
		case notVisiting:
			if c.isCycle(k, false) {
				return true
			}
		}
	}
	c.typVisits[v] = visited
	return false
}
