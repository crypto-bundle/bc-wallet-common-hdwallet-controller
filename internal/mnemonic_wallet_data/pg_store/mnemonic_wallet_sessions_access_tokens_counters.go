/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package pg_store

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"

	"github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (s *pgRepository) GetCurrentAccessTokenCounterValue(ctx context.Context,
	tokenIdentifier uuid.UUID,
) (*entities.AccessTokenWalletSessionCounter, error) {
	var result *entities.AccessTokenWalletSessionCounter

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		const q = `
			SELECT * FROM "wallet_sessions_access_tokens_counters"
			WHERE "token_uuid" = $1
			FOR UPDATE;`

		row := stmt.QueryRowx(q, tokenIdentifier)

		item := &entities.AccessTokenWalletSessionCounter{}
		clbErr := row.StructScan(item)
		if clbErr != nil {
			return postgres.EmptyOrError(clbErr, "failed to select last counter values")
		}

		result = item

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pgRepository) GetNextWalletSessionCounterValue(ctx context.Context,
	tokenIdentifier uuid.UUID,
) (int32, error) {
	var nextNonce int32

	err := s.pgConn.MustWithTransaction(ctx, func(stmt *sqlx.Tx) error {
		var tmp int32

		q := `
			INSERT INTO wallet_sessions_access_tokens_counters (token_uuid, counter_value)
				VALUES ($1, 0)
			ON CONFLICT (token_uuid) 
				DO UPDATE
					SET counter_value = wallet_sessions_access_tokens_counters.counter_value + 1
			RETURNING counter_value`

		row := stmt.QueryRowx(q, tokenIdentifier)
		clbErr := row.Scan(&tmp)
		if clbErr != nil {
			return clbErr
		}

		nextNonce = tmp

		return nil
	})

	if err != nil {
		return -1, err
	}

	return nextNonce, nil
}
