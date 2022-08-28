package parameters

func NewParamSpec(name string, tp Type, opt TypeOptions)*paramSpec{
	return &paramSpec{Name:name, Type:tp, Options:opt}
}


func (s *Signature)ParamsByName() map[string]*paramSpec{
	return s.paramsByName
}
