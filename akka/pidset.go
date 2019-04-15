package akka

const actorOfSetSliceLen = 16

type ActorOfSet struct {
	s []string
	m map[string]struct{}
}

// NewActorOfSet returns a new ActorOfSet with the given  actorOfs.
func NewActorOfSet(actorOfs ...*ActorRef) *ActorOfSet {
	var s ActorOfSet
	for _, actorOf := range actorOfs {
		s.Add(actorOf)
	}
	return &s
}

func (p *ActorOfSet) indexOf(v *ActorRef) int {
	key := v.key()
	for i, actorOf := range p.s {
		if key == actorOf {
			return i
		}
	}
	return -1
}

func (p *ActorOfSet) migrate() {
	p.m = make(map[string]struct{}, actorOfSetSliceLen)
	for _, v := range p.s {
		p.m[v] = struct{}{}
	}
	p.s = p.s[:0]
}

// Add adds the element v to the set
func (p *ActorOfSet) Add(v *ActorRef) {
	if p.m == nil {
		if p.indexOf(v) > -1 {
			return
		}

		if len(p.s) < actorOfSetSliceLen {
			if p.s == nil {
				p.s = make([]string, 0, actorOfSetSliceLen)
			}
			p.s = append(p.s, v.key())
			return
		}
		p.migrate()
	}
	p.m[v.key()] = struct{}{}
}

// Remove removes v from the set and returns true if them element existed
func (p *ActorOfSet) Remove(v *ActorRef) bool {
	if p.m == nil {
		i := p.indexOf(v)
		if i == -1 {
			return false
		}
		l := len(p.s) - 1
		p.s[i] = p.s[l]
		p.s = p.s[:l]
		return true
	}
	_, ok := p.m[v.key()]
	if !ok {
		return false
	}
	delete(p.m, v.key())
	return true
}

// Contains reports whether v is an element of the set
func (p *ActorOfSet) Contains(v *ActorRef) bool {
	if p.m == nil {
		return p.indexOf(v) != -1
	}
	_, ok := p.m[v.key()]
	return ok
}

// Len returns the number of elements in the set
func (p *ActorOfSet) Len() int {
	if p.m == nil {
		return len(p.s)
	}
	return len(p.m)
}

// Clear removes all the elements in the set
func (p *ActorOfSet) Clear() {
	if p.m == nil {
		p.s = p.s[:0]
	} else {
		p.m = nil
	}
}

// Empty reports whether the set is empty
func (p *ActorOfSet) Empty() bool {
	return p.Len() == 0
}

// Values returns all the elements of the set as a slice
func (p *ActorOfSet) Values() []ActorRef {
	if p.Len() == 0 {
		return nil
	}

	r := make([]ActorRef, p.Len())
	if p.m == nil {
		for i, v := range p.s {
			actorOfFromKey(v, &r[i])
		}
	} else {
		i := 0
		for v := range p.m {
			actorOfFromKey(v, &r[i])
			i++
		}
	}
	return r
}

// ForEach invokes f for every element of the set
func (p *ActorOfSet) ForEach(f func(i int, actorOf ActorRef)) {
	var actorOf ActorRef
	if p.m == nil {
		for i, v := range p.s {
			actorOfFromKey(v, &actorOf)
			f(i, actorOf)
		}
	} else {
		i := 0
		for v := range p.m {
			actorOfFromKey(v, &actorOf)
			f(i, actorOf)
			i++
		}
	}
}

func (p *ActorOfSet) Clone() *ActorOfSet {
	var s ActorOfSet
	if p.s != nil {
		s.s = make([]string, len(p.s))
		for i, v := range p.s {
			s.s[i] = v
		}
	}
	if p.m != nil {
		s.m = make(map[string]struct{}, len(p.m))
		for v := range p.m {
			s.m[v] = struct{}{}
		}
	}
	return &s
}
