package files

type diskAdapter struct {
}

func (d *diskAdapter) createFile(filename string, data []byte, contentType string) {

}

func (d *diskAdapter) deleteFile(filename string) {

}
func (d *diskAdapter) getFileData(filename string) []byte {
	return nil
}
func (d *diskAdapter) getFileLocation(filename string) string {
	return ""
}
