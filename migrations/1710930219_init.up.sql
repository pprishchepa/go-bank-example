CREATE TABLE wallet
(
    id         BIGSERIAL   NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE wallet_balance
(
    wallet_id BIGINT NOT NULL PRIMARY KEY,
    amount    INT    NOT NULL DEFAULT 0,
    CONSTRAINT fk_wallet FOREIGN KEY (wallet_id) REFERENCES wallet (id) ON DELETE CASCADE,
    CONSTRAINT amount_nonnegative CHECK (amount >= 0)
);

CREATE TABLE wallet_entry
(
    id            BIGSERIAL   NOT NULL PRIMARY KEY,
    wallet_id     BIGINT      NOT NULL,
    debit_amount  INT                  DEFAULT NULL,
    credit_amount INT                  DEFAULT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_wallet FOREIGN KEY (wallet_id) REFERENCES wallet (id) ON DELETE CASCADE,
    CONSTRAINT either_debit_or_credit CHECK (
        (debit_amount IS NULL AND credit_amount IS NOT NULL
            OR debit_amount IS NOT NULL AND credit_amount IS NULL) AND
        (COALESCE(debit_amount, 0) >= 0 AND COALESCE(credit_amount, 0) >= 0))
);

INSERT INTO wallet (id)
VALUES (101),
       (102);

INSERT INTO wallet_balance (wallet_id, amount)
VALUES (101, 1860043),
       (102, 4277);

INSERT INTO wallet_entry (wallet_id, debit_amount, credit_amount)
VALUES (101, 2000000, NULL),
       (101, NULL, 139957),
       (102, 4000, NULL),
       (102, 277, NULL);

