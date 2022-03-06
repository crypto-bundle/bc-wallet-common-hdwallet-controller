package vault

import (
	b64 "encoding/base64"
)

const (
	plainTxt    = "plaintext"
	cipherTxt   = "ciphertext"
	encryptPath = "transit/encrypt/"
	decryptPath = "transit/decrypt/"
)

// Encrypt get encrypted ciphertext bytes via vault transit secret engine.
func (s *service) Encrypt(toEncrypt []byte) ([]byte, error) {
	b64Val := b64.StdEncoding.EncodeToString(toEncrypt)
	path := encryptPath + s.cfg.TransitKey

	secret, err := s.client.Logical().Write(path, map[string]interface{}{
		plainTxt: b64Val,
	})
	if err != nil {
		return nil, err
	}

	encrVal := secret.Data[cipherTxt]
	encrValStr, ok := encrVal.(string)
	if !ok {
		return nil, ErrTransitSecretFormat
	}

	return []byte(encrValStr), nil
}

// Decrypt get decrypted value from ciphertext via vault transit secret engine.
func (s *service) Decrypt(cipherBytes []byte) ([]byte, error) {
	path := decryptPath + s.cfg.TransitKey

	secret, err := s.client.Logical().Write(path, map[string]interface{}{
		cipherTxt: string(cipherBytes),
	})
	if err != nil {
		return nil, err
	}

	decrVal := secret.Data[plainTxt]
	decrValStr, ok := decrVal.(string)
	if !ok {
		return nil, ErrTransitSecretFormat
	}

	return b64.StdEncoding.DecodeString(decrValStr)
}
