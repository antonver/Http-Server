package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	headers_ := NewHeaders()
	// Test: Valid single header
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers_.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers_)
	assert.Equal(t, "localhost:42069", headers_["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers_ = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers_.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	t.Run("Velid 2 headers with existing headers", func (t *testing.T){
		headers_ = NewHeaders()
		test_header1 := "Host: localhost:42069\r\n\r\n"
		test_header2 := "Host: localhost:42069\r\n\r\n"
		m, done, err := headers_.Parse([]byte(test_header1))
		n, done2, err2 := headers_.Parse([]byte(test_header2))
		require.NoError(t, err)
		require.NoError(t, err2)
		assert.Equal(t,n,m)
		assert.Equal(t,len(headers_),1)
		assert.False(t, done)
		assert.False(t, done2)
	})

	t.Run("Valide done", func (t *testing.T){
		headers_ = NewHeaders()
		test_header1 := "Host: localhost:42069\r\n\r\n"
		test_header2 := "\r\n"
		_, done, err := headers_.Parse([]byte(test_header1))
		_, done2, err2 := headers_.Parse([]byte(test_header2))
		require.NoError(t, err)
		require.NoError(t, err2)
		assert.False(t, done)
		assert.True(t, done2)
	})
		t.Run("Wrong character", func (t *testing.T){
		headers_ = NewHeaders()
		test_header1 := "H©st: localhost:42069\r\n\r\n"
		_, done, err := headers_.Parse([]byte(test_header1))
		require.Error(t, err)
		assert.False(t, done)
	})
	t.Run("Two identical keys", func (t *testing.T){
		headers_ = NewHeaders()
		test_header1 := "Host: example1\r\n\r\n"
		test_header2 := "host: example2\r\n\r\n"
		_, done, err := headers_.Parse([]byte(test_header1))
		_, done2, err2 := headers_.Parse([]byte(test_header2))
		require.NoError(t, err)
		require.NoError(t, err2)
		assert.False(t, done)
		assert.False(t, done2)
		assert.Equal(t, len(headers_), 1)
		assert.Equal(t, headers_["host"], "example1,example2")
	})
}
