package console

func (con *SpecterClient) Printf(format string, args ...any) {
	con.printf(format, args...)
}
