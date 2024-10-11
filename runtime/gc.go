package runtime

func (s *state) new(value Value) int {
	location := len(s.heap)
	s.heap = append(s.heap, value)
	value.SetLocation(location)
	return location
}
