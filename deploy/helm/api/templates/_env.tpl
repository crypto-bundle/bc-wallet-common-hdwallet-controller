{{- define "_env" }}
- name: APP_ENV
  value: {{ pluck .Values.global.env .Values.app.environment | first | default .Values.app.environment._default | quote }}
- name: APP_DEBUG
  value: {{ pluck .Values.global.env .Values.app.debug_mode | first | default .Values.app.debug_mode._default | quote }}
- name: APP_STAGE
  value: {{ pluck .Values.global.env .Values.app.stage.name | first | default .Values.app.stage.name._default | quote }}

- name: LOGGER_LEVEL
  value: {{ pluck .Values.global.env .Values.app.logger.minimal_level | first | default .Values.app.logger.minimal_level._default | quote }}
- name: LOGGER_STACKTRACE_ENABLE
  value: {{ pluck .Values.global.env .Values.app.logger.enabled_stack_trace | first | default .Values.app.logger.enabled_stack_trace._default | quote }}

- name: API_GRPC_PORT
  value: {{ pluck .Values.global.env .Values.app.api.grpc_port | first | default .Values.app.api.grpc_port._default | quote }}

- name: DB_HOST
  value: {{ pluck .Values.global.env .Values.app.db.host | first | default .Values.app.db.host._default | quote }}
- name: DB_PORT
  value: {{ pluck .Values.global.env .Values.app.db.port | first | default .Values.app.db.port._default | quote }}
- name: DB_MAX_OPEN_CONNECTIONS
  value: {{ pluck .Values.global.env .Values.app.db.open_connections | first | default .Values.app.db.open_connections._default | quote }}
- name: DB_MAX_IDLE_CONNECTIONS
  value: {{ pluck .Values.global.env .Values.app.db.idle_connections | first | default .Values.app.db.idle_connections._default | quote }}
- name: DB_SSL_MODE
  value: {{ pluck .Values.global.env .Values.app.db.ssl_mode | first | default .Values.app.db.ssl_mode._default | quote }}
- name: DB_DATABASE
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: db_name
      optional: false
- name: DB_USERNAME
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: db_username
      optional: false
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: db_password
      optional: false

- name: NATS_ADDRESSES
  value: {{ pluck .Values.global.env .Values.app.nats.hosts | first | default .Values.app.nats.hosts._default | join "," | quote }}
- name: NATS_CONNECTION_RETRY
  value: {{ pluck .Values.global.env .Values.app.nats.connection_retry | first | default .Values.app.nats.connection_retry._default | quote }}
- name: NATS_CONNECTION_RETRY_COUNT
  value: {{ pluck .Values.global.env .Values.app.nats.connection_retry_count | first | default .Values.app.nats.connection_retry_count._default | quote }}
- name: NATS_CONNECTION_RETRY_TIMEOUT
  value: {{ pluck .Values.global.env .Values.app.nats.connection_retry_timeout | first | default .Values.app.nats.connection_retry_timeout._default | quote }}
- name: NATS_FLUSH_TIMEOUT
  value: {{ pluck .Values.global.env .Values.app.nats.flush_timeout | first | default .Values.app.nats.flush_timeout._default | quote }}
- name: NATS_WORKER_PER_CONSUMER
  value: {{ pluck .Values.global.env .Values.app.nats.workers | first | default .Values.app.nats.workers._default | quote }}
- name: NATS_USER
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: nats_username
      optional: false
- name: NATS_PASSWORD
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: nats_password
      optional: false

- name: REDIS_HOST
  value: {{ pluck .Values.global.env .Values.app.redis.host | first | default .Values.app.redis.host._default | quote }}
- name: REDIS_PORT
  value: {{ pluck .Values.global.env .Values.app.redis.port | first | default .Values.app.redis.port._default | quote }}
- name: REDIS_CONNECTION_RETRY_TIMEOUT
  value: {{ pluck .Values.global.env .Values.app.redis.connection_retry_timeout | first | default .Values.app.redis.connection_retry_timeout._default | quote }}
- name: REDIS_CONNECTION_RETRY_COUNT
  value: {{ pluck .Values.global.env .Values.app.redis.connection_retry_count | first | default .Values.app.redis.connection_retry_count._default | quote }}
- name: REDIS_MAX_RETRY_COUNT
  value: {{ pluck .Values.global.env .Values.app.redis.max_retry_count | first | default .Values.app.redis.max_retry_count._default | quote }}
- name: REDIS_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.app.redis.read_timeout | first | default .Values.app.redis.read_timeout._default | quote }}
- name: REDIS_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.app.redis.write_timeout | first | default .Values.app.redis.write_timeout._default | quote }}
- name: REDIS_MIN_IDLE_CONNECTIONS
  value: {{ pluck .Values.global.env .Values.app.redis.min_idle_connections | first | default .Values.app.redis.min_idle_connections._default | quote }}
- name: REDIS_IDLE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.app.redis.idle_timeout | first | default .Values.app.redis.idle_timeout._default | quote }}
- name: REDIS_MAX_CONNECTION_AGE
  value: {{ pluck .Values.global.env .Values.app.redis.connection_age | first | default .Values.app.redis.connection_age._default | quote }}
- name: REDIS_POOL_SIZE
  value: {{ pluck .Values.global.env .Values.app.redis.pool_size | first | default .Values.app.redis.pool_size._default | quote }}
- name: REDIS_POOL_TIMEOUT
  value: {{ pluck .Values.global.env .Values.app.redis.pool_timeout | first | default .Values.app.redis.pool_timeout._default | quote }}
- name: REDIS_DIAL_TIMEOUT
  value: {{ pluck .Values.global.env .Values.app.redis.dial_timeout | first | default .Values.app.redis.dial_timeout._default | quote }}
- name: REDIS_USER
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: redis_username
      optional: false
- name: REDIS_PASSWORD
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: redis_password
      optional: false

- name: VAULT_HOST
  value: {{ pluck .Values.global.env .Values.app.vault.host | first | default .Values.app.vault.host._default | quote }}
- name: VAULT_PORT
  value: {{ pluck .Values.global.env .Values.app.vault.port | first | default .Values.app.vault.port._default | quote }}
- name: VAULT_USE_HTTPS
  value: {{ pluck .Values.global.env .Values.app.vault.use_https | first | default .Values.app.vault.use_https._default | quote }}
- name: VAULT_DATA_PATH
  value: {{ pluck .Values.global.env .Values.app.vault.data_path | first | default .Values.app.vault.data_path._default }}
- name: VAULT_AUTH_METHOD
  value: {{ pluck .Values.global.env .Values.app.vault.auth_method | first | default .Values.app.vault.auth_method._default }}
- name: VAULT_AUTH_TOKEN
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: vault_auth_path
      optional: false
- name: VAULT_TRANSIT_KEY
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: vault_transit_secret_key
      optional: false

- name: HDWALLET_WORDS_COUNT
  value: {{ pluck .Values.global.env .Values.app.mnemonic.words_count | first | default .Values.app.mnemonic.words_count._default | quote }}
{{- end }}