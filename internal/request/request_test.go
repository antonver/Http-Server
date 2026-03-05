package request

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

func TestRequestFromReader(t *testing.T) {
	t.Run("GET with 3 bytes chunk", func(t *testing.T) {
		reader := &chunkReader{
			data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 3,
		}
		r, err := RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	})

	t.Run("GET with 1 byte chunk", func(t *testing.T) {
		reader := &chunkReader{
			data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 1,
		}
		r, err := RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	})

	t.Run("POST Request", func(t *testing.T) {
		r, err := RequestFromReader(strings.NewReader("POST /cool HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n "))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "POST", r.RequestLine.Method)
		assert.Equal(t, "/cool", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	})

	t.Run("Validation Errors", func(t *testing.T) {
		_, err := RequestFromReader(strings.NewReader("/cool POST HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
		require.Error(t, err)

		_, err = RequestFromReader(strings.NewReader("POST /cool HTTP/2.1\r\nHost: localhost:42069\r\n\r\n"))
		require.Error(t, err)
	})
	reader := &chunkReader{
	data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	numBytesPerRead: 3,
}

	t.Run("Standard Headers", func(t *testing.T) {
reader := &chunkReader{
	data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	numBytesPerRead: 3,
}
r, err := RequestFromReader(reader)
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "localhost:42069", r.Headers["host"])
assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
assert.Equal(t, "*/*", r.Headers["accept"])
	})



		t.Run("Malformed Header", func(t *testing.T) {
reader = &chunkReader{
	data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
	numBytesPerRead: 3,
}
_, err := RequestFromReader(reader)
require.Error(t, err)
	})


		t.Run("Empty Headers", func(t *testing.T) {
reader := &chunkReader{
	data:            "GET / HTTP/1.1\r\n\r\n",
	numBytesPerRead: 3,
}
_, err := RequestFromReader(reader)
require.Error(t, err)

	})

			t.Run("Empty Headers", func(t *testing.T) {
reader := &chunkReader{
	data:            "GET / HTTP/1.1\r\n\r\n",
	numBytesPerRead: 3,
}
_, err := RequestFromReader(reader)
require.Error(t, err)

	})

		t.Run("Duplicated Headers", func(t *testing.T) {
reader := &chunkReader{
	data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nUser-Agent: duplicate\r\n\r\n",
	numBytesPerRead: 3,
}
r, err := RequestFromReader(reader)
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "localhost:42069", r.Headers["host"])
assert.Equal(t, "curl/7.81.0,duplicate", r.Headers["user-agent"])

	})

			t.Run("Case insensitive headers", func(t *testing.T) {
reader := &chunkReader{
	data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nuser-Agent: curl/7.81.0\r\nUser-Agent: duplicate\r\n\r\n",
	numBytesPerRead: 3,
}
r, err := RequestFromReader(reader)
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "localhost:42069", r.Headers["host"])
assert.Equal(t, "curl/7.81.0,duplicate", r.Headers["user-agent"])
assert.Equal(t, len(r.Headers), 2)
	})
				t.Run("Missing End of Headers", func(t *testing.T) {
reader := &chunkReader{
	data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nuser-Agent: curl/7.81.0\r\nUser-Agent: duplicate\r\n",
	numBytesPerRead: 3,
}
_, err := RequestFromReader(reader)
require.Error(t, err)

	})
}


