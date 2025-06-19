package utils

func Ifelse(cond bool, t, f any) any {
	if cond {
		return t
	}
	return f
}
