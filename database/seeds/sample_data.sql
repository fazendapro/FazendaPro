-- FazendaPro Database Seed Data
-- Description: Test data for the farm management system
-- Created: 2024-01-21
-- 
-- This file contains sample data for testing all the functionalities

-- Insert test users (RF01)
INSERT INTO users (name, email, password_hash) VALUES
('João Silva', 'joao@fazendapro.com', '$2b$12$LQv3c1yqBwqVzqqqBBYAqOE2w0eKfVlS9owU3Ql1nq5Pvu8qNv.Q2'),
('Maria Santos', 'maria@fazenda.com', '$2b$12$LQv3c1yqBwqVzqqqBBYAqOE2w0eKfVlS9owU3Ql1nq5Pvu8qNv.Q2'),
('Pedro Oliveira', 'pedro@rancho.com', '$2b$12$LQv3c1yqBwqVzqqqBBYAqOE2w0eKfVlS9owU3Ql1nq5Pvu8qNv.Q2');

-- Insert production lots (RF08)
INSERT INTO lots (name, description, min_milk_production, max_milk_production) VALUES
('Lote A - Alta Produção', 'Vacas de alta produção leiteira', 25.00, 50.00),
('Lote B - Produção Média', 'Vacas de produção média', 10.00, 24.99),
('Lote C - Novilhas', 'Animais jovens em desenvolvimento', 0.00, 9.99),
('Lote D - Reprodutores', 'Touros e reprodutores', 0.00, 0.00);

-- Insert sample animals (RF02, RF03)
INSERT INTO animals (user_id, lot_id, name, ear_tag, birth_date, race, gender, mother_id, father_id, status, current_weight, daily_milk_production, notes) VALUES
-- Touros (reprodutores)
(1, 4, 'Touro Nelore Chief', 'T001', '2020-03-15', 'Nelore', 'male', NULL, NULL, 'active', 850.00, 0.00, 'Touro reprodutor principal'),
(1, 4, 'Touro Angus Thunder', 'T002', '2019-05-20', 'Angus', 'male', NULL, NULL, 'active', 780.00, 0.00, 'Touro de elite para melhoramento genético'),

-- Vacas matrizes (sem mãe e pai conhecidos)
(1, 1, 'Vaca Estrela', 'V001', '2019-01-10', 'Holandesa', 'female', NULL, NULL, 'active', 650.00, 32.50, 'Vaca de alta produção, matriz principal'),
(1, 1, 'Vaca Luna', 'V002', '2019-03-22', 'Gir Leiteiro', 'female', NULL, NULL, 'active', 580.00, 28.00, 'Excelente produtora, resistente ao calor'),
(1, 2, 'Vaca Bonita', 'V003', '2020-07-08', 'Girolando', 'female', NULL, NULL, 'active', 520.00, 18.75, 'Boa produtora, fácil manejo'),
(1, 2, 'Vaca Serena', 'V004', '2020-09-15', 'Jersey', 'female', NULL, NULL, 'active', 380.00, 22.00, 'Leite com alta concentração de gordura'),

-- Descendentes (filhos das matrizes)
(1, 3, 'Novilha Esperança', 'N001', '2022-02-14', 'Holandesa', 'female', 3, 1, 'active', 420.00, 8.50, 'Filha da Estrela, grande potencial genético'),
(1, 3, 'Novilho Forte', 'N002', '2022-04-18', 'Gir Leiteiro', 'male', 4, 1, 'active', 480.00, 0.00, 'Filho da Luna, candidato a reprodutor'),
(1, 3, 'Bezerro Júnior', 'B001', '2023-01-20', 'Girolando', 'male', 5, 2, 'active', 280.00, 0.00, 'Filho da Bonita, em crescimento'),
(1, 3, 'Bezerra Princesa', 'B002', '2023-03-10', 'Jersey', 'female', 6, 2, 'active', 250.00, 0.00, 'Filha da Serena, promissora');

-- Insert sample vaccines (RF11)
INSERT INTO vaccines (user_id, name, manufacturer, description, dosage, application_method, validity_period_days, notes) VALUES
(1, 'Vacina contra Febre Aftosa', 'Boehringer Ingelheim', 'Proteção contra febre aftosa tipos O, A e C', '2ml', 'Subcutânea', 180, 'Aplicar a cada 6 meses'),
(1, 'Vacina Brucelose B19', 'MSD Saúde Animal', 'Prevenção da brucelose bovina', '2ml', 'Subcutânea', 365, 'Aplicar apenas em fêmeas de 3-8 meses'),
(1, 'Vacina Raiva', 'Ourofino', 'Profilaxia da raiva em bovinos', '2ml', 'Intramuscular', 365, 'Aplicação anual obrigatória'),
(1, 'Clostridioses', 'Zoetis', 'Proteção contra clostrídios (gangrena gasosa)', '5ml', 'Subcutânea', 180, 'Reforço aos 30 dias da primeira aplicação'),
(1, 'IBR/BVD', 'MSD Saúde Animal', 'Rinotraqueíte infecciosa e diarreia viral', '2ml', 'Intramuscular', 365, 'Importante para reprodução');

-- Insert vaccination records (RF06)
INSERT INTO animal_vaccines (animal_id, vaccine_id, application_date, next_application_date, batch_number, administered_by, notes) VALUES
-- Vacinações da Vaca Estrela
(3, 1, '2024-01-15', '2024-07-15', 'FMD240115', 'João Silva', 'Primeira dose do ano'),
(3, 3, '2024-02-20', '2025-02-20', 'RAB240220', 'João Silva', 'Vacinação anual'),
(3, 5, '2024-03-10', '2025-03-10', 'IBR240310', 'João Silva', 'Pré-cobertura'),

-- Vacinações da Vaca Luna
(4, 1, '2024-01-15', '2024-07-15', 'FMD240115', 'João Silva', 'Junto com o lote'),
(4, 3, '2024-02-20', '2025-02-20', 'RAB240220', 'João Silva', 'Vacinação de rotina'),
(4, 4, '2024-01-30', '2024-07-30', 'CLS240130', 'Veterinário', 'Primeira dose'),

-- Vacinações das novilhas
(7, 1, '2024-01-20', '2024-07-20', 'FMD240115', 'João Silva', 'Primeira vacinação'),
(7, 2, '2023-08-15', '2024-08-15', 'BRU230815', 'Veterinário', 'Aplicada na idade correta'),
(10, 2, '2023-09-10', '2024-09-10', 'BRU230910', 'Veterinário', 'Bezerra dentro da idade');

-- Insert weight records (RF07)
INSERT INTO weight_records (animal_id, weight, measurement_date, measurement_type, notes) VALUES
-- Registros da Vaca Estrela
(3, 620.00, '2024-01-01', 'monthly', 'Peso inicial do ano'),
(3, 635.00, '2024-02-01', 'monthly', 'Ganho de peso satisfatório'),
(3, 650.00, '2024-03-01', 'monthly', 'Peso atual - em boa forma'),

-- Registros da Novilha Esperança
(7, 380.00, '2024-01-01', 'monthly', 'Início do acompanhamento'),
(7, 400.00, '2024-02-01', 'monthly', 'Bom desenvolvimento'),
(7, 420.00, '2024-03-01', 'monthly', 'Crescimento dentro do esperado'),

-- Registros do Bezerro Júnior (mais frequentes)
(9, 220.00, '2024-01-01', 'weekly', 'Início do controle'),
(9, 235.00, '2024-01-08', 'weekly', 'Primeira semana'),
(9, 250.00, '2024-01-15', 'weekly', 'Crescimento acelerado'),
(9, 265.00, '2024-01-22', 'weekly', 'Mantendo o ritmo'),
(9, 280.00, '2024-01-29', 'weekly', 'Peso atual');

-- Insert pregnancy records (RF09)
INSERT INTO pregnancy_records (animal_id, father_id, pregnancy_date, expected_birth_date, actual_birth_date, pregnancy_status, notification_sent, notes) VALUES
(3, 1, '2023-08-15', '2024-05-15', '2024-05-18', 'completed', TRUE, 'Parto normal, bezerra saudável'),
(4, 1, '2023-09-10', '2024-06-10', '2024-06-08', 'completed', TRUE, 'Parto antecipado mas sem complicações'),
(5, 2, '2023-10-20', '2024-07-20', '2024-07-22', 'completed', TRUE, 'Bezerro macho nasceu bem'),
(6, 2, '2023-12-05', '2024-09-05', '2024-09-10', 'completed', TRUE, 'Bezerra fêmea, primeira cria'),
(3, 1, '2024-01-20', '2024-10-20', NULL, 'confirmed', FALSE, 'Nova gestação confirmada, notificação agendada'),
(7, 2, '2024-02-14', '2024-11-14', NULL, 'confirmed', FALSE, 'Primeira gestação da novilha');

-- Insert suppliers (RF12)
INSERT INTO suppliers (user_id, name, company_name, cnpj, cpf, email, phone, address, city, state, supplier_type, notes) VALUES
(1, 'Ração Premium Ltda', 'Ração Premium Comércio e Indústria Ltda', '12.345.678/0001-90', NULL, 'vendas@racaopremium.com.br', '(31) 3456-7890', 'Rua das Indústrias, 123', 'Belo Horizonte', 'MG', 'feed', 'Fornecedor principal de ração'),
(1, 'VetMinas', 'Veterinária Minas Gerais S.A.', '98.765.432/0001-10', NULL, 'comercial@vetminas.com.br', '(31) 2345-6789', 'Av. Veterinários, 456', 'Contagem', 'MG', 'medicine', 'Medicamentos e vacinas'),
(1, 'Equipamentos Rurais Cerrado', 'Equipamentos Rurais Cerrado Ltda', '11.222.333/0001-44', NULL, 'vendas@cerradoequip.com.br', '(31) 4567-8901', 'Rod. MG-040, Km 25', 'Ribeirão das Neves', 'MG', 'equipment', 'Ordenhadeiras e equipamentos'),
(1, 'José Santos - Técnico', NULL, NULL, '123.456.789-00', 'josesantos.tecnico@email.com', '(31) 99876-5432', 'Rua do Campo, 789', 'Sete Lagoas', 'MG', 'services', 'Inseminação artificial e manejo'),
(1, 'Maria Oliveira - Veterinária', NULL, NULL, '987.654.321-00', 'dra.maria@email.com', '(31) 98765-4321', 'Av. Principal, 321', 'Pedro Leopoldo', 'MG', 'services', 'Consultas e cirurgias veterinárias');

-- Insert sales records (RF10)
INSERT INTO sales (animal_id, buyer_name, buyer_document, buyer_contact, sale_date, sale_price, payment_method, payment_status, notes) VALUES
-- Venda fictícia de um animal que depois foi marcado como vendido
(8, 'Fazenda São João', '12.345.678/0001-00', '(31) 99999-0000', '2024-02-15', 3500.00, 'bank_transfer', 'completed', 'Novilho vendido para reprodução'),
-- Venda pendente
(9, 'Roberto Silva', '123.456.789-00', 'roberto@email.com', '2024-03-10', 2800.00, 'installments', 'partial', 'Pagamento em 3x, primeira parcela paga');

-- Update status of sold animals
UPDATE animals SET status = 'sold' WHERE id = 8;