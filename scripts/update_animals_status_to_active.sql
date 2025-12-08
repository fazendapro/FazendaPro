-- Script SQL para atualizar todos os animais de uma fazenda para status 0 (Ativo)
-- Execute este script no banco de dados para resetar o status de todos os animais
-- 
-- IMPORTANTE: Substitua o valor de farm_id abaixo pela ID da fazenda desejada
-- Exemplo: Se quiser atualizar a fazenda com ID 1, use: WHERE farm_id = 1;

-- Atualizar todos os animais da fazenda 1 para status 0 (Ativo)
UPDATE animals 
SET status = 0, 
    updated_at = CURRENT_TIMESTAMP
WHERE farm_id = 1;

-- Verificar quantos animais foram atualizados
SELECT COUNT(*) as animais_atualizados
FROM animals 
WHERE farm_id = 1 AND status = 0;

-- Para atualizar TODAS as fazendas (use com cuidado!):
-- UPDATE animals 
-- SET status = 0, 
--     updated_at = CURRENT_TIMESTAMP;

-- Para atualizar uma fazenda específica por ID (substitua X pelo ID da fazenda):
-- UPDATE animals 
-- SET status = 0, 
--     updated_at = CURRENT_TIMESTAMP
-- WHERE farm_id = X;

-- Para atualizar apenas animais vendidos (status = 1) de uma fazenda:
-- UPDATE animals 
-- SET status = 0, 
--     updated_at = CURRENT_TIMESTAMP
-- WHERE farm_id = 1 AND status = 1;

-- Para atualizar apenas animais falecidos (status = 2) de uma fazenda:
-- UPDATE animals 
-- SET status = 0, 
--     updated_at = CURRENT_TIMESTAMP
-- WHERE farm_id = 1 AND status = 2;

