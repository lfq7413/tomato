package files

type diskAdapter struct {
}

func (d *diskAdapter) createFile(filename string, data []byte, contentType string) error {
	return nil
}

func (d *diskAdapter) deleteFile(filename string) error {
	return nil
}
func (d *diskAdapter) getFileData(filename string) []byte {
	return nil
}
func (d *diskAdapter) getFileLocation(filename string) string {
	return ""
}
