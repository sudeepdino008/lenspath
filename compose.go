package lenspath

// this package provides ways to compose lenspath

func (lp *Lenspath) Compose(lens []Lens) *Lenspath {
	return &Lenspath{
		lens:      append(lp.lens, lens...),
		assumeNil: lp.assumeNil,
	}
}
