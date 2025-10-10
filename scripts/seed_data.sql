INSERT INTO companies (id, company_name, created_at, updated_at) 
VALUES (1, 'FazendaPro Demo', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO farms (id, company_id, name, logo, created_at, updated_at) 
VALUES (1, 1, 'Fazenda Demo', '', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO companies (id, company_name, created_at, updated_at) 
VALUES (2, 'Fazenda Teste', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO farms (id, company_id, name, logo, created_at, updated_at) 
VALUES (2, 2, 'Fazenda Teste', '', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
