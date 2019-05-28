package streams

type BaseFlow struct {
	Flow
}

func (f *BaseFlow) Map(work func(msg interface{}) interface{}) *Transform {
	result := &Transform{
		work: work,
	}
	f.actors = append(f.actors, result)
	return result
}

func (f *BaseFlow) Foreach(work func(msg interface{})) *Sink {
	result := &Sink{
		work: work,
	}
	f.actors = append(f.actors, result)
	return result
}

func (f *BaseFlow) Filter(work func(msg interface{}) bool) *Filter {
	result := &Filter{
		work: work,
	}
	f.actors = append(f.actors, result)
	return result
}

func (f *BaseFlow) Window(count int) *Window {
	result := &Window{
		count: count,
	}
	f.actors = append(f.actors, result)
	return result
}
