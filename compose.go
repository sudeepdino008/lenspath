package lenspath

// this package provides ways to compose lenspath

func (lp *Lenspath) Compose(lens []Lens) (*Lenspath, error) {
	newlens := append(lp.lens, lens...)
	copylens := make([]Lens, len(newlens))
	copy(copylens, newlens)

	if lens, err := Create(copylens); err != nil {
		return nil, err
	} else {
		return lens, nil
	}
}
