# Documentação da API - Gestor Simples

Este documento detalha todos os endpoints da API necessários para o funcionamento do aplicativo **Gestor Simples**.

**URL Base da API:** `[URL_DA_SUA_API]/api/v1`

---

## 1. Autenticação

Recurso para autenticar usuários e gerar tokens de acesso.

### **`POST /auth/login`**

-   **Descrição:** Autentica um usuário com base em email e senha.
-   **Corpo da Requisição (`application/json`):**
    ```json
    {
      "email": "admin@example.com",
      "password": "password123"
    }
    ```
-   **Resposta de Sucesso (`200 OK`):** Retorna os dados do usuário e um token JWT para ser usado em requisições autenticadas.
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": 1,
        "name": "Administrador",
        "email": "admin@example.com",
        "role": "admin"
      }
    }
    ```
-   **Resposta de Erro (`401 Unauthorized`):**
    ```json
    {
      "message": "Credenciais inválidas."
    }
    ```

### **`POST /auth/register`**

-   **Descrição:** Registra um novo usuário (vendedor).
-   **Corpo da Requisição (`application/json`):**
    ```json
    {
      "name": "Novo Vendedor",
      "username": "novo.vendedor",
      "password": "password123"
    }
    ```
-   **Resposta de Sucesso (`201 Created`):** Retorna os dados do usuário recém-criado.
    ```json
    {
      "id": 5,
      "name": "Novo Vendedor",
      "username": "novo.vendedor",
      "role": "vendedor"
    }
    ```
-   **Resposta de Erro (`400 Bad Request`):**
    ```json
    {
      "error": "Invalid request payload"
    }
    ```
-   **Resposta de Erro (`500 Internal Server Error`):**
    ```json
    {
      "error": "Failed to create user"
    }
    ```

---

## 2. Usuários (Vendedores)

Endpoints para o gerenciamento de usuários (principalmente vendedores).

### **`GET /users`**

-   **Descrição:** Lista todos os usuários. Acesso restrito para `admin`. Pode ser filtrado por role.
-   **Query Params (Opcional):**
    -   `role` (string): Filtra usuários por perfil. Ex: `/users?role=vendedor`
-   **Resposta de Sucesso (`200 OK`):**
    ```json
    [
      {
        "id": 2,
        "name": "João Silva",
        "email": "joao.silva@example.com",
        "role": "vendedor"
      },
      {
        "id": 3,
        "name": "Maria Santos",
        "email": "maria.santos@example.com",
        "role": "vendedor"
      }
    ]
    ```

### **`GET /users/{id}`**

-   **Descrição:** Obtém os detalhes de um usuário específico.
-   **Resposta de Sucesso (`200 OK`):**
    ```json
    {
      "id": 2,
      "name": "João Silva",
      "email": "joao.silva@example.com",
      "role": "vendedor"
    }
    ```
-   **Resposta de Erro (`404 Not Found`):** Se o usuário não for encontrado.

### **`POST /users`**

-   **Descrição:** Cria um novo usuário (vendedor). Acesso restrito para `admin`.
-   **Corpo da Requisição (`application/json`):**
    ```json
    {
      "name": "Carlos Pereira",
      "email": "carlos.p@example.com",
      "password": "newSecurePassword",
      "role": "vendedor"
    }
    ```
-   **Resposta de Sucesso (`201 Created`):**
    ```json
    {
      "id": 4,
      "name": "Carlos Pereira",
      "email": "carlos.p@example.com",
      "role": "vendedor"
    }
    ```

### **`PUT /users/{id}`**

-   **Descrição:** Atualiza os dados de um usuário existente.
-   **Corpo da Requisição (`application/json`):**
    ```json
    {
      "name": "João da Silva",
      "email": "joao.silva.novo@example.com"
    }
    ```
-   **Resposta de Sucesso (`200 OK`):**
    ```json
    {
      "id": 2,
      "name": "João da Silva",
      "email": "joao.silva.novo@example.com",
      "role": "vendedor"
    }
    ```

### **`DELETE /users/{id}`**

-   **Descrição:** Remove um usuário do sistema.
-   **Resposta de Sucesso (`204 No Content`):** Nenhum corpo na resposta.

---

## 3. Produtos (Estoque)

Endpoints para gerenciamento do catálogo de produtos e estoque.

### **`GET /products`**

-   **Descrição:** Lista todos os produtos disponíveis.
-   **Resposta de Sucesso (`200 OK`):**
    ```json
    [
      {
        "id": 1,
        "name": "Produto A",
        "description": "Descrição detalhada do Produto A.",
        "price": 29.99,
        "quantity": 150
      },
      {
        "id": 2,
        "name": "Produto B",
        "description": "Descrição detalhada do Produto B.",
        "price": 199.9,
        "quantity": 45
      }
    ]
    ```

### **`POST /products`**

-   **Descrição:** Adiciona um novo produto ao estoque.
-   **Corpo da Requisição (`application/json`):**
    ```json
    {
      "name": "Produto C",
      "description": "Novo produto adicionado.",
      "price": 50.00,
      "quantity": 200
    }
    ```
-   **Resposta de Sucesso (`201 Created`):**
    ```json
    {
      "id": 3,
      "name": "Produto C",
      "description": "Novo produto adicionado.",
      "price": 50.00,
      "quantity": 200
    }
    ```

### **`PUT /products/{id}`**

-   **Descrição:** Atualiza um produto existente (preço, quantidade, etc.).
-   **Corpo da Requisição (`application/json`):**
    ```json
    {
      "price": 32.50,
      "quantity": 140
    }
    ```
-   **Resposta de Sucesso (`200 OK`):**
    ```json
    {
      "id": 1,
      "name": "Produto A",
      "description": "Descrição detalhada do Produto A.",
      "price": 32.50,
      "quantity": 140
    }
    ```

### **`DELETE /products/{id}`**

-   **Descrição:** Remove um produto do catálogo.
-   **Resposta de Sucesso (`204 No Content`):** Nenhum corpo na resposta.

---

## 4. Vendas

Endpoints para registrar e consultar vendas.

### **`GET /sales`**

-   **Descrição:** Retorna o histórico de vendas. Pode ser filtrado.
-   **Query Params (Opcional):**
    -   `userId` (number): Filtra vendas por um vendedor específico.
    -   `startDate` (date): Data de início do período (formato `YYYY-MM-DD`).
    -   `endDate` (date): Data de fim do período (formato `YYYY-MM-DD`).
-   **Resposta de Sucesso (`200 OK`):**
    ```json
    [
      {
        "id": 1,
        "userId": 2,
        "date": "2025-11-20T14:30:00Z",
        "items": [
          {
            "productId": 1,
            "productName": "Produto A",
            "quantity": 2,
            "unitPrice": 32.50
          },
          {
            "productId": 2,
            "productName": "Produto B",
            "quantity": 1,
            "unitPrice": 199.90
          }
        ],
        "totalPrice": 264.90
      }
    ]
    ```

### **`POST /sales`**

-   **Descrição:** Registra uma nova venda. O backend deve validar se há estoque suficiente e decrementar a quantidade do produto.
-   **Corpo da Requisição (`application/json`):**
    ```json
    {
      "userId": 2,
      "items": [
        {
          "productId": 1,
          "quantity": 2
        },
        {
          "productId": 2,
          "quantity": 1
        }
      ]
    }
    ```
-   **Resposta de Sucesso (`201 Created`):**
    ```json
    {
      "saleId": 2
    }
    ```
-   **Resposta de Erro (`400 Bad Request`):** Se o produto não tiver estoque suficiente.
    ```json
    {
      "message": "Estoque insuficiente para o produto: Produto B"
    }
    ```

---

## 5. Dashboards e Relatórios

Endpoints para obter dados consolidados.

### **`GET /dashboard/summary`**

-   **Descrição:** Obtém dados agregados para o dashboard. A resposta pode variar com base na `role` do usuário que faz a requisição.
-   **Resposta de Sucesso (`200 OK` para Admin):**
    ```json
    {
      "totalSalesMonth": 7580.50,
      "totalSellers": 15,
      "lowStockProducts": 8,
      "topSellingProduct": {
        "id": 2,
        "name": "Produto B"
      }
    }
    ```
-   **Resposta de Sucesso (`200 OK` para Vendedor):**
    ```json
    {
      "myTotalSalesMonth": 1250.75,
      "myRank": 3,
      "commissions": 125.07
    }
    ```
