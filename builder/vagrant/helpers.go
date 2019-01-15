package vagrant

import (
	"bytes"
	"io"
	"log"
)

// WriteVagrantfile takes a path to a Vagrantfile file and contents in the form of a
// map and writes it out.
func WriteVagrantfile(path string, data map[string]string) (err error) {
	log.Printf("Writing Vagrantfile to: %s", path)
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	var buf bytes.Buffer
	buf.WriteString(EncodeVagrantfile(data))
	if _, err = io.Copy(f, &buf); err != nil {
		return
	}

	return
}

// ReadVagrantfile takes a path to a Vagrantfile file and reads it into a k/v mapping.
func ReadVagrantfile(path string) (map[string]string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseVagrantfile(string(data)), nil
}
