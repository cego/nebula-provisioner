package store

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/openbao/openbao/sdk/v2/helper/shamir"
)

type storeEncryption struct {
	keyPartsUnseal [][]byte
}

func (e *storeEncryption) generate(numParts, threshold uint32) ([]byte, []string, error) {
	// Generating encryption key
	ek := make([]byte, 32)
	_, err := rand.Read(ek)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate random encryption key")
	}

	keyParts, err := e.getKeyParts(ek, numParts, threshold)
	if err != nil {
		return nil, nil, err
	}

	return ek, keyParts, nil
}

func (e *storeEncryption) getKeyParts(ek []byte, numParts, threshold uint32) ([]string, error) {
	secretParts, err := shamir.Split(ek, int(numParts), int(threshold))
	if err != nil {
		return nil, fmt.Errorf("failed to spilt encryption key using shamir: %v", err)
	}

	var keyParts []string
	for _, bytePart := range secretParts {
		keyParts = append(keyParts, hex.EncodeToString(bytePart))
	}
	return keyParts, nil
}

func (e *storeEncryption) unseal(keyPart string, removeExistingParts bool) ([]byte, error) {
	if removeExistingParts {
		e.keyPartsUnseal = nil
	}

	decodedPart, err := hex.DecodeString(keyPart)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key part")
	}

	e.appendIfMissing(decodedPart)

	ek, err := shamir.Combine(e.keyPartsUnseal)
	if err != nil {
		return nil, fmt.Errorf("failed to combine encryption key using shamir: %s", err)
	}

	return ek, nil
}

func (e *storeEncryption) clear() {
	e.keyPartsUnseal = nil
}

func (e *storeEncryption) appendIfMissing(b []byte) {
	for _, ele := range e.keyPartsUnseal {
		if bytes.Compare(ele, b) == 0 {
			return
		}
	}
	e.keyPartsUnseal = append(e.keyPartsUnseal, b)
}
