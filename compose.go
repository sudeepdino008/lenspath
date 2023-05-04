package lenspath

// this package provides ways to compose lenspath

func (lp *Lenspath) Compose(lens []Lens) *Lenspath {
	newlens := append(lp.lens, lens...)
	copylens := make([]Lens, len(newlens))
	copy(copylens, newlens)

	return &Lenspath{
		lens:      copylens,
		assumeNil: lp.assumeNil,
	}
}
