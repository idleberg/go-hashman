package algo

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"

	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
)

// Algorithm defines a hash algorithm with its CLI flag, display name, and constructor.
type Algorithm struct {
	ID         string
	Flag       string
	Display    string
	NewHash    func() hash.Hash
	Deprecated bool
}

var (
	crc32cTable = crc32.MakeTable(crc32.Castagnoli)
	crc64Table  = crc64.MakeTable(crc64.ISO)
)

// Registry contains all supported hash algorithms in display order.
var Registry = []Algorithm{
	{ID: "adler32", Flag: "adler32", Display: "Adler-32", NewHash: func() hash.Hash { return adler32.New() }},
	{ID: "crc32", Flag: "crc32", Display: "CRC32", NewHash: func() hash.Hash { return crc32.NewIEEE() }},
	{ID: "crc32c", Flag: "crc32c", Display: "CRC32C", NewHash: func() hash.Hash { return crc32.New(crc32cTable) }},
	{ID: "crc64", Flag: "crc64", Display: "CRC64", NewHash: func() hash.Hash { return crc64.New(crc64Table) }},
	{ID: "md4", Flag: "md4", Display: "MD4", NewHash: md4.New, Deprecated: true},
	{ID: "md5", Flag: "md5", Display: "MD5", NewHash: md5.New},
	{ID: "ripemd160", Flag: "rmd160", Display: "RIPEMD-160", NewHash: ripemd160.New, Deprecated: true},
	{ID: "sha1", Flag: "sha-1", Display: "SHA-1", NewHash: sha1.New},
	{ID: "sha224", Flag: "sha-224", Display: "SHA-224", NewHash: sha256.New224},
	{ID: "sha256", Flag: "sha-256", Display: "SHA-256", NewHash: sha256.New},
	{ID: "sha384", Flag: "sha-384", Display: "SHA-384", NewHash: sha512.New384},
	{ID: "sha512", Flag: "sha-512", Display: "SHA-512", NewHash: sha512.New},
	{ID: "sha3224", Flag: "sha3-224", Display: "SHA3-224", NewHash: func() hash.Hash { return sha3.New224() }},
	{ID: "sha3256", Flag: "sha3-256", Display: "SHA3-256", NewHash: func() hash.Hash { return sha3.New256() }},
	{ID: "sha3384", Flag: "sha3-384", Display: "SHA3-384", NewHash: func() hash.Hash { return sha3.New384() }},
	{ID: "sha3512", Flag: "sha3-512", Display: "SHA3-512", NewHash: func() hash.Hash { return sha3.New512() }},
}
