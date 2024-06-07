/*
 *
 *
 * MIT-License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package mocks

import (
	"bytes"
	"context"
	"sync"
)

type obscurityDataStore struct {
	mu sync.Mutex

	identifierByWallets map[string][]byte
}

func (s *obscurityDataStore) GetLastObscurityData(ctx context.Context,
	walletUUID string,
	accessTokenHash string,
) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, isExists := s.identifierByWallets[accessTokenHash]
	if !isExists {
		return nil, nil
	}

	return data, nil
}

func (s *obscurityDataStore) AddLastObscurityData(ctx context.Context,
	walletUUID string,
	accessTokenHash string,
	obscurityData []byte,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.identifierByWallets[accessTokenHash] = bytes.Clone(obscurityData)

	return nil
}

func NewObscurityDataStoreStore(identifiers map[string][]byte) *obscurityDataStore {
	return &obscurityDataStore{
		mu:                  sync.Mutex{},
		identifierByWallets: identifiers,
	}
}
