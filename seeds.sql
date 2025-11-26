-- SQL script para popular o banco de dados com dados iniciais para a aplicação Gestor Simples.

-- Observação: As senhas não são armazenadas em texto plano. A aplicação deve usar um algoritmo de hash seguro
-- como o bcrypt. Para o propósito deste script de semente, estamos inserindo um hash de espaço reservado.
-- A senha para todos os usuários é 'password123'.

-- Usuários
INSERT INTO Users (name, username, password_hash, role) VALUES
('Administrador', 'admin', '$2a$10$Q7.Y.aD5.d6t.9sR/6U6OOE/5.1J2A1QY6.Rz.F3B8c3j2a7C5y.O', 'admin'), -- password123
('João Silva', 'joao.silva', '$2a$10$Q7.Y.aD5.d6t.9sR/6U6OOE/5.1J2A1QY6.Rz.F3B8c3j2a7C5y.O', 'vendedor'), -- password123
('Maria Santos', 'maria.santos', '$2a$10$Q7.Y.aD5.d6t.9sR/6U6OOE/5.1J2A1QY6.Rz.F3B8c3j2a7C5y.O', 'vendedor'); -- password123

-- Produtos
INSERT INTO Products (name, description, quantity, price) VALUES
('Produto A', 'Descrição detalhada do Produto A.', 150, 29.99),
('Produto B', 'Descrição detalhada do Produto B.', 45, 199.90),
('Produto C', 'Novo produto adicionado.', 200, 50.00);

-- Vendas
-- Uma venda feita por João Silva (user_id = 2)
INSERT INTO Sales (user_id, date) VALUES
(2, '2025-11-20 14:30:00');

-- Itens de Venda para a venda acima
-- 2 unidades do Produto A e 1 unidade do Produto B
INSERT INTO Sales_Items (sale_id, product_id, quantity) VALUES
(1, 1, 2),
(1, 2, 1);
