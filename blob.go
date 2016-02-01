package git

import (
	"bytes"
	"io"
	"io/ioutil"
)

// Blob represents a Git object.
type Blob struct {
	repo *Repository
	*TreeEntry
}

// Data gets content of blob all at once and wrap it as io.Reader.
// This can be very slow and memory consuming for huge content.
func (b *Blob) Data() (io.Reader, error) {
	_, body, err := b.repo.repo.OpenObject(sha2oidp(b.ID))
	if err != nil {
		return nil, err
	}

	// FIXME: we can implement something like EOFCloser who autoclose its reader
	// after it have been fully read. So we can avoid the buffering
	// or we can change the API and return ReadCloser and force the caller to
	// close object.

	buf, err := ioutil.ReadAll(body)
	body.Close()
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(buf), nil
}

func (b *Blob) DataPipeline(stdout, stderr io.Writer) error {
	_, body, err := b.repo.repo.OpenObject(sha2oidp(b.ID))
	if err != nil {
		return err
	}
	_, err = io.Copy(stdout, body)
	return err
}
