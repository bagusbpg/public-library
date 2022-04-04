package helper

func NilHandler(_x float64) interface{} {
	if _x == 0 {
		return nil
	}

	return _x
}
