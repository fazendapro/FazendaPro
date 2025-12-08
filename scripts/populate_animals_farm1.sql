-- Script SQL para popular a farm_id 1 com animais
-- Execute este script no banco de dados para adicionar animais de teste

-- Inserir vacas (fêmeas adultas)
INSERT INTO animals (
    farm_id, 
    ear_tag_number_local, 
    ear_tag_number_register, 
    animal_name, 
    sex, 
    breed, 
    type, 
    birth_date, 
    animal_type, 
    status, 
    confinement, 
    fertilization, 
    castrated, 
    purpose, 
    current_batch, 
    created_at, 
    updated_at
) VALUES
-- Vacas leiteiras
(1, 1001, 1001, 'Daisy', 0, 'Holandesa', 'vaca', '2018-03-15', 1, 1, false, true, false, 1, 1, NOW(), NOW()),
(1, 1002, 1002, 'Bella', 0, 'Holandesa', 'vaca', '2019-05-20', 1, 1, false, false, false, 1, 1, NOW(), NOW()),
(1, 1003, 1003, 'Luna', 0, 'Jersey', 'vaca', '2017-11-10', 1, 1, false, true, false, 1, 1, NOW(), NOW()),
(1, 1004, 1004, 'Maggie', 0, 'Holandesa', 'vaca', '2018-08-22', 1, 1, false, false, false, 1, 1, NOW(), NOW()),
(1, 1005, 1005, 'Rosie', 0, 'Gir', 'vaca', '2019-01-30', 1, 1, false, true, false, 1, 1, NOW(), NOW()),
(1, 1006, 1006, 'Molly', 0, 'Holandesa', 'vaca', '2018-07-05', 1, 1, false, false, false, 1, 1, NOW(), NOW()),
(1, 1007, 1007, 'Penny', 0, 'Jersey', 'vaca', '2017-12-18', 1, 1, false, true, false, 1, 1, NOW(), NOW()),
(1, 1008, 1008, 'Coco', 0, 'Gir', 'vaca', '2019-04-12', 1, 1, false, false, false, 1, 1, NOW(), NOW()),

-- Touros (machos adultos)
(1, 2001, 2001, 'Thor', 1, 'Nelore', 'touro', '2016-06-10', 1, 1, false, false, false, 2, 1, NOW(), NOW()),
(1, 2002, 2002, 'Zeus', 1, 'Angus', 'touro', '2017-02-25', 1, 1, false, false, false, 2, 1, NOW(), NOW()),
(1, 2003, 2003, 'Max', 1, 'Nelore', 'touro', '2016-11-08', 1, 1, false, false, false, 2, 1, NOW(), NOW()),
(1, 2004, 2004, 'Rocky', 1, 'Brahman', 'touro', '2017-09-14', 1, 1, false, false, false, 2, 1, NOW(), NOW()),

-- Novilhos (machos jovens)
(1, 3001, 3001, 'Bolinha', 1, 'Nelore', 'novilho', '2022-01-20', 1, 1, false, false, true, 0, 1, NOW(), NOW()),
(1, 3002, 3002, 'Pipoca', 1, 'Angus', 'novilho', '2022-03-15', 1, 1, false, false, true, 0, 1, NOW(), NOW()),
(1, 3003, 3003, 'Trovão', 1, 'Nelore', 'novilho', '2022-05-10', 1, 1, false, false, false, 0, 1, NOW(), NOW()),
(1, 3004, 3004, 'Raio', 1, 'Brahman', 'novilho', '2022-07-22', 1, 1, false, false, true, 0, 1, NOW(), NOW()),
(1, 3005, 3005, 'Foguinho', 1, 'Nelore', 'novilho', '2022-09-05', 1, 1, false, false, false, 0, 1, NOW(), NOW()),
(1, 3006, 3006, 'Tempestade', 1, 'Angus', 'novilho', '2022-11-18', 1, 1, false, false, true, 0, 1, NOW(), NOW()),

-- Bezerros (machos e fêmeas jovens)
(1, 4001, 4001, 'Pequeno', 1, 'Nelore', 'bezerro', '2023-08-10', 1, 1, false, false, false, 0, 1, NOW(), NOW()),
(1, 4002, 4002, 'Mini', 0, 'Holandesa', 'bezerro', '2023-09-15', 1, 1, false, false, false, 1, 1, NOW(), NOW()),
(1, 4003, 4003, 'Tiny', 1, 'Angus', 'bezerro', '2023-10-20', 1, 1, false, false, false, 0, 1, NOW(), NOW()),
(1, 4004, 4004, 'Pequenino', 0, 'Jersey', 'bezerro', '2023-11-25', 1, 1, false, false, false, 1, 1, NOW(), NOW()),

-- Mais vacas para produção
(1, 5001, 5001, 'Estrela', 0, 'Holandesa', 'vaca', '2018-02-14', 1, 1, false, true, false, 1, 1, NOW(), NOW()),
(1, 5002, 5002, 'Flor', 0, 'Gir', 'vaca', '2019-06-08', 1, 1, false, false, false, 1, 1, NOW(), NOW()),
(1, 5003, 5003, 'Mel', 0, 'Jersey', 'vaca', '2018-10-30', 1, 1, false, true, false, 1, 1, NOW(), NOW()),
(1, 5004, 5004, 'Doce', 0, 'Holandesa', 'vaca', '2019-03-22', 1, 1, false, false, false, 1, 1, NOW(), NOW()),
(1, 5005, 5005, 'Açúcar', 0, 'Gir', 'vaca', '2018-12-05', 1, 1, false, true, false, 1, 1, NOW(), NOW()),

-- Mais novilhos
(1, 6001, 6001, 'Forte', 1, 'Nelore', 'novilho', '2021-04-18', 1, 1, false, false, true, 0, 1, NOW(), NOW()),
(1, 6002, 6002, 'Bravo', 1, 'Angus', 'novilho', '2021-07-25', 1, 1, false, false, false, 0, 1, NOW(), NOW()),
(1, 6003, 6003, 'Valente', 1, 'Nelore', 'novilho', '2021-10-12', 1, 1, false, false, true, 0, 1, NOW(), NOW()),

-- Mais bezerros
(1, 7001, 7001, 'Bebê', 0, 'Holandesa', 'bezerro', '2024-01-10', 1, 1, false, false, false, 1, 1, NOW(), NOW()),
(1, 7002, 7002, 'Filhote', 1, 'Nelore', 'bezerro', '2024-02-15', 1, 1, false, false, false, 0, 1, NOW(), NOW()),
(1, 7003, 7003, 'Cria', 0, 'Jersey', 'bezerro', '2024-03-20', 1, 1, false, false, false, 1, 1, NOW(), NOW());

-- Verificar quantos animais foram inseridos
SELECT COUNT(*) as total_animais FROM animals WHERE farm_id = 1;

-- Listar todos os animais inseridos
SELECT 
    id,
    ear_tag_number_local,
    animal_name,
    type,
    sex,
    breed,
    birth_date,
    purpose,
    status
FROM animals 
WHERE farm_id = 1 
ORDER BY ear_tag_number_local;



