/*
 * MIT License
 *
 * Copyright (c) 2021-2023 Aleksei Kotelnikov
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
 */

-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX uuid__idx ON mnemonic_wallets (
     "uuid"
);
CREATE UNIQUE INDEX wallet_uuid__idx ON mnemonic_wallets (
    "wallet_uuid"
);

CREATE UNIQUE INDEX wallet_uuid__mnemonic_wallet_uuid__idx ON mnemonic_wallets (
    "wallet_uuid", "uuid"
);

CREATE INDEX is_hot__idx ON mnemonic_wallets (
    "is_hot"
);

CREATE INDEX wallet_uuid__is_hot__idx ON mnemonic_wallets (
   "wallet_uuid", "is_hot"
);

CREATE UNIQUE INDEX hash__idx ON mnemonic_wallets ("hash");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX uuid__idx;
DROP INDEX wallet_uuid__idx;
DROP INDEX wallet_uuid__mnemonic_wallet_uuid__idx;
DROP INDEX is_hot__idx;
DROP INDEX wallet_uuid__is_hot__idx;
DROP INDEX hash__idx;
-- +goose StatementEnd
