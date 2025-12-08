-- Script para popular a tabela de produção de leite (milk_collections)
-- Este script cria registros de produção de leite para os animais existentes

-- Verificar se há animais cadastrados
DO $$
DECLARE
    animal_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO animal_count FROM animals;
    
    IF animal_count = 0 THEN
        RAISE EXCEPTION 'Nenhum animal encontrado na tabela. Crie animais primeiro.';
    END IF;
    
    RAISE NOTICE 'Encontrados % animais. Iniciando população de dados...', animal_count;
END $$;

-- Limpar dados existentes (opcional - descomente se necessário)
-- TRUNCATE TABLE milk_collections RESTART IDENTITY CASCADE;

-- Inserir registros de produção de leite
-- Varia a produção diária entre 15L e 45L por animal
-- Insere para TODOS os animais (remova o filtro de sexo se necessário)

DO $$
DECLARE
    animal_record RECORD;
    days_ago INTEGER;
    production_liters NUMERIC;
    record_date DATE;
    records_per_animal INTEGER;
    i INTEGER;
    j INTEGER;
    inserted_count INTEGER := 0;
    total_animals INTEGER := 0;
BEGIN
    -- Contar animais que serão processados
    SELECT COUNT(*) INTO total_animals 
    FROM animals 
    WHERE farm_id IS NOT NULL;
    
    RAISE NOTICE 'Processando % animais...', total_animals;
    
    -- Para cada animal, criar múltiplos registros
    FOR animal_record IN 
        SELECT id, animal_name, ear_tag_number_local, sex, farm_id
        FROM animals 
        WHERE farm_id IS NOT NULL
        ORDER BY id
    LOOP
        -- Cada animal terá entre 3 e 6 registros nos últimos 30 dias
        records_per_animal := 3 + FLOOR(RANDOM() * 4)::INTEGER;
        
        FOR i IN 1..records_per_animal LOOP
            -- Tentar até encontrar uma data sem registro para este animal
            FOR j IN 1..10 LOOP
                -- Data aleatória nos últimos 30 dias
                days_ago := FLOOR(RANDOM() * 30)::INTEGER;
                record_date := CURRENT_DATE - days_ago;
                
                -- Verificar se já existe registro para este animal nesta data
                IF NOT EXISTS (
                    SELECT 1 
                    FROM milk_collections 
                    WHERE animal_id = animal_record.id 
                      AND date = record_date
                ) THEN
                    -- Produção variando entre 15L e 45L
                    production_liters := ROUND((RANDOM() * 30 + 15)::NUMERIC, 2);
                    
                    -- Inserir registro
                    INSERT INTO milk_collections (animal_id, liters, date, created_at, updated_at)
                    VALUES (animal_record.id, production_liters, record_date, NOW(), NOW());
                    
                    inserted_count := inserted_count + 1;
                    EXIT; -- Sair do loop interno quando inserir com sucesso
                END IF;
            END LOOP;
        END LOOP;
    END LOOP;
    
    RAISE NOTICE 'Registros históricos inseridos: %', inserted_count;
END $$;

-- Inserir registros adicionais para os últimos 7 dias (produção mais recente)
DO $$
DECLARE
    animal_record RECORD;
    days_ago INTEGER;
    production_liters NUMERIC;
    record_date DATE;
    i INTEGER;
    j INTEGER;
    inserted_count INTEGER := 0;
BEGIN
    FOR animal_record IN 
        SELECT id, animal_name, ear_tag_number_local, sex, farm_id
        FROM animals 
        WHERE farm_id IS NOT NULL
        ORDER BY id
    LOOP
        -- Cada animal terá 1-3 registros nos últimos 7 dias
        FOR i IN 1..(1 + FLOOR(RANDOM() * 3)::INTEGER) LOOP
            -- Tentar até encontrar uma data sem registro
            FOR j IN 1..10 LOOP
                -- Data aleatória nos últimos 7 dias
                days_ago := FLOOR(RANDOM() * 7)::INTEGER;
                record_date := CURRENT_DATE - days_ago;
                
                -- Verificar se já existe registro
                IF NOT EXISTS (
                    SELECT 1 
                    FROM milk_collections 
                    WHERE animal_id = animal_record.id 
                      AND date = record_date
                ) THEN
                    -- Produção variando entre 15L e 45L
                    production_liters := ROUND((RANDOM() * 30 + 15)::NUMERIC, 2);
                    
                    -- Inserir registro
                    INSERT INTO milk_collections (animal_id, liters, date, created_at, updated_at)
                    VALUES (animal_record.id, production_liters, record_date, NOW(), NOW());
                    
                    inserted_count := inserted_count + 1;
                    EXIT; -- Sair do loop interno quando inserir com sucesso
                END IF;
            END LOOP;
        END LOOP;
    END LOOP;
    
    RAISE NOTICE 'Registros recentes inseridos: %', inserted_count;
END $$;

-- Verificar quantos registros foram inseridos
DO $$
DECLARE
    total_registros INTEGER;
    animais_com_producao INTEGER;
    data_antiga DATE;
    data_recente DATE;
    media_litros NUMERIC;
BEGIN
    SELECT COUNT(*), COUNT(DISTINCT animal_id), MIN(date), MAX(date), ROUND(AVG(liters)::NUMERIC, 2)
    INTO total_registros, animais_com_producao, data_antiga, data_recente, media_litros
    FROM milk_collections;
    
    RAISE NOTICE '========================================';
    RAISE NOTICE 'RESUMO DA POPULAÇÃO:';
    RAISE NOTICE 'Total de registros: %', total_registros;
    RAISE NOTICE 'Animais com produção: %', animais_com_producao;
    RAISE NOTICE 'Data mais antiga: %', data_antiga;
    RAISE NOTICE 'Data mais recente: %', data_recente;
    RAISE NOTICE 'Média de litros: %', media_litros;
    RAISE NOTICE '========================================';
END $$;

-- Estatísticas detalhadas
SELECT 
    COUNT(*) as total_registros,
    COUNT(DISTINCT animal_id) as animais_com_producao,
    MIN(date) as data_mais_antiga,
    MAX(date) as data_mais_recente,
    ROUND(AVG(liters)::numeric, 2) as media_litros,
    ROUND(MIN(liters)::numeric, 2) as menor_producao,
    ROUND(MAX(liters)::numeric, 2) as maior_producao
FROM milk_collections;

-- Mostrar produção por animal (top 15)
SELECT 
    a.id,
    a.animal_name,
    a.ear_tag_number_local,
    a.sex,
    COUNT(mc.id) as total_coletas,
    ROUND(SUM(mc.liters)::numeric, 2) as total_litros,
    ROUND(AVG(mc.liters)::numeric, 2) as media_litros,
    MIN(mc.date) as primeira_coleta,
    MAX(mc.date) as ultima_coleta
FROM animals a
INNER JOIN milk_collections mc ON a.id = mc.animal_id
GROUP BY a.id, a.animal_name, a.ear_tag_number_local, a.sex
ORDER BY total_litros DESC
LIMIT 15;
