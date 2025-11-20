CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

SET timezone = 'America/Sao_Paulo';

DO $$
BEGIN
    RAISE NOTICE 'Banco de dados FazendaPro inicializado com sucesso!';
END $$; 