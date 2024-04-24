// Package keyconvert
// взял отсюда https://github.com/wemeetagain/go-hdwallet
// несколько переписал, чтоб можно было работать с разными криптовалютами
package hdwallet

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
)

// HDWallet defines the components of a hierarchical deterministic hdwallet
type HDWallet struct {
	prvMagic    [4]byte
	pubMagic    [4]byte
	Vbytes      []byte // 4 bytes
	Depth       uint16 // 1 byte
	Fingerprint []byte // 4 bytes
	I           []byte // 4 bytes
	Chaincode   []byte // 32 bytes
	Key         []byte // 33 bytes
}

// Child returns the ith child of hdwallet w. Values of i >= 2^31
// signify private key derivation. Attempting private key derivation
// with a public key will throw an error.
func (w *HDWallet) Child(i uint32) (*HDWallet, error) {
	var fingerprint, I, newkey []byte
	switch {
	case bytes.Equal(w.Vbytes, w.prvMagic[:4]):
		pub := privToPub(w.Key)
		mac := hmac.New(sha512.New, w.Chaincode)
		if i >= uint32(0x80000000) {
			_, writeErr := mac.Write(append(w.Key, uint32ToByte(i)...))
			if writeErr != nil {
				return nil, writeErr
			}
		} else {
			_, writeErr := mac.Write(append(pub, uint32ToByte(i)...))
			if writeErr != nil {
				return nil, writeErr
			}
		}

		I = mac.Sum(nil)
		iL := new(big.Int).SetBytes(I[:32])
		if iL.Cmp(curve.N) >= 0 || iL.Sign() == 0 {
			return &HDWallet{}, errors.New("invalid child")
		}
		newkey = addPrivKeys(I[:32], w.Key)
		raw, err := hash160(privToPub(w.Key))
		if err != nil {
			return nil, err
		}

		fingerprint = raw[:4]

	case bytes.Equal(w.Vbytes, w.pubMagic[:4]):
		mac := hmac.New(sha512.New, w.Chaincode)
		if i >= uint32(0x80000000) {
			return &HDWallet{}, errors.New("can't do private derivation on public key")
		}
		_, writeErr := mac.Write(append(w.Key, uint32ToByte(i)...))
		if writeErr != nil {
			return nil, writeErr
		}

		I = mac.Sum(nil)
		iL := new(big.Int).SetBytes(I[:32])
		if iL.Cmp(curve.N) >= 0 || iL.Sign() == 0 {
			return &HDWallet{}, errors.New("invalid child")
		}

		newkey = addPubKeys(privToPub(I[:32]), w.Key)
		raw, err := hash160(w.Key)
		if err != nil {
			return nil, err
		}

		fingerprint = raw[:4]
	}
	return &HDWallet{w.prvMagic, w.pubMagic, w.Vbytes, w.Depth + 1, fingerprint, uint32ToByte(i), I[32:], newkey}, nil
}

// Serialize returns the serialized form of the hdwallet.
func (w *HDWallet) Serialize() ([]byte, error) {
	depth := uint16ToByte(w.Depth % 256)

	bindata := make([]byte, 78)
	copy(bindata, w.Vbytes)
	copy(bindata[4:], depth)
	copy(bindata[5:], w.Fingerprint)
	copy(bindata[9:], w.I)
	copy(bindata[13:], w.Chaincode)
	copy(bindata[45:], w.Key)
	raw, err := dblSha256(bindata)
	if err != nil {
		return nil, err
	}

	chksum := raw[:4]

	return append(bindata, chksum...), nil
}

// String returns the base58-encoded string form of the hdwallet.
func (w *HDWallet) String() (string, error) {
	serialized, err := w.Serialize()
	if err != nil {
		return "", err
	}

	return base58.Encode(serialized), nil
}

// StringWallet returns a hdwallet given a base58-encoded extended key
func StringWallet(data string, prvMagic, pubMagic [4]byte) (*HDWallet, error) {
	dbin := base58.Decode(data)
	if err := ByteCheck(dbin, prvMagic, pubMagic); err != nil {
		return &HDWallet{}, err
	}

	sha256calc, err := dblSha256(dbin[:(len(dbin) - 4)])
	if err != nil {
		return &HDWallet{}, err
	}

	if !bytes.Equal(sha256calc[:4], dbin[(len(dbin)-4):]) {
		return &HDWallet{}, errors.New("invalid checksum")
	}

	vbytes := dbin[0:4]
	depth := byteToUint16(dbin[4:5])
	fingerprint := dbin[5:9]
	i := dbin[9:13]
	chaincode := dbin[13:45]
	key := dbin[45:78]

	return &HDWallet{
		prvMagic:    prvMagic,
		pubMagic:    pubMagic,
		Vbytes:      vbytes,
		Depth:       depth,
		Fingerprint: fingerprint,
		I:           i,
		Chaincode:   chaincode,
		Key:         key,
	}, nil
}

// Pub returns a new hdwallet which is the public key version of w.
// If w is a public key, Pub returns a copy of w
func (w *HDWallet) Pub() *HDWallet {
	if bytes.Equal(w.Vbytes, w.pubMagic[:4]) {
		return &HDWallet{w.prvMagic, w.pubMagic, w.Vbytes, w.Depth, w.Fingerprint, w.I, w.Chaincode, w.Key}
	}

	return &HDWallet{w.prvMagic, w.pubMagic, w.pubMagic[:4], w.Depth, w.Fingerprint, w.I, w.Chaincode, privToPub(w.Key)}
}

// StringChild returns the ith base58-encoded extended key of a base58-encoded extended key.
func StringChild(data string, i uint32, prvMagic, pubMagic [4]byte) (string, error) {
	w, err := StringWallet(data, prvMagic, pubMagic)
	if err != nil {
		return "", err
	}

	w, err = w.Child(i)
	if err != nil {
		return "", err
	}

	str, err := w.String()
	if err != nil {
		return "", err
	}

	return str, nil
}

// StringAddress returns the Bitcoin address of a base58-encoded extended key.
func StringAddress(data string, prvMagic, pubMagic [4]byte) (string, error) {
	w, err := StringWallet(data, prvMagic, pubMagic)
	if err != nil {
		return "", err
	}

	addr, err := w.Address()
	if err != nil {
		return "", err
	}

	return addr, nil
}

// Address returns bitcoin address represented by hdwallet w.
func (w *HDWallet) Address() (string, error) {
	x, y := expand(w.Key)
	paddedKey, err := hex.DecodeString("04")
	if err != nil {
		return "", err
	}

	paddedKey = append(paddedKey, append(x.Bytes(), y.Bytes()...)...)
	addr1, err := hex.DecodeString("00")
	if err != nil {
		return "", err
	}

	raw, err := hash160(paddedKey)
	if err != nil {
		return "", err
	}

	addr1 = append(addr1, raw...)
	chkSum, err := dblSha256(addr1)
	if err != nil {
		return "", err
	}

	return base58.Encode(append(addr1, chkSum[:4]...)), nil
}

// GenSeed returns a random seed with a length measured in bytes.
// The length must be at least 128.
func GenSeed(length int) ([]byte, error) {
	b := make([]byte, length)
	if length < 128 {
		return b, errors.New("length must be at least 128 bits")
	}
	_, err := rand.Read(b)
	return b, err
}

// MasterKey returns a new hdwallet given a random seed.
func MasterKey(seed []byte, prvMagic, pubMagic [4]byte) (*HDWallet, error) {
	key := []byte("Bitcoin seed")
	mac := hmac.New(sha512.New, key)
	_, err := mac.Write(seed)
	if err != nil {
		return nil, err
	}

	I := mac.Sum(nil)
	secret := I[:len(I)/2]
	chainCode := I[len(I)/2:]
	depth := 0
	i := make([]byte, 4)
	fingerprint := make([]byte, 4)
	zero := make([]byte, 1)

	return &HDWallet{
		prvMagic:    prvMagic,
		pubMagic:    pubMagic,
		Vbytes:      prvMagic[:4],
		Depth:       uint16(depth),
		Fingerprint: fingerprint,
		I:           i,
		Chaincode:   chainCode,
		Key:         append(zero, secret...),
	}, nil
}

func ByteCheck(dBin []byte, prvMagic, pubMagic [4]byte) error {
	// check proper length
	if len(dBin) != 82 {
		return errors.New("invalid string")
	}
	// check for correct Public or Private vbytes
	if !bytes.Equal(dBin[:4], pubMagic[:4]) && !bytes.Equal(dBin[:4], prvMagic[:4]) {
		return errors.New("invalid string")
	}

	// if Public, check x coord is on curve
	x, y := expand(dBin[45:78])
	if bytes.Equal(dBin[:4], pubMagic[:4]) {
		if !onCurve(x, y) {
			return errors.New("invalid string")
		}
	}
	return nil
}
