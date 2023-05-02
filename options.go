package lenspath

type LenspathOptions func(*Lenspath) error

func (lp *Lenspath) WithOptions(opts ...LenspathOptions) error {
	for _, opt := range opts {
		if err := opt(lp); err != nil {
			return err
		}
	}

	return nil
}

// WithAssumeNil sets the Lenspath.assumeNil field to the given value. This would be used when
// the Lenspath is used for "get" operations, and the user wants to assume nil when the Lenspath
// cannot be resolved.
func WithAssumeNil(assumeNil bool) LenspathOptions {
	return func(lp *Lenspath) error {
		lp.assumeNil = assumeNil
		return nil
	}
}
