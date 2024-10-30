package afs

func marshalError(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
