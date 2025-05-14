# FazendaPro - Soluções Agropecuárias

- **Título do Projeto**: [Título claro e conciso que reflete a essência do produto ou ferramenta].
- **Nome do Estudante**: Gustavo Henrique Dias.
- **Curso**: Engenharia de Software.
- **Data de Entrega**: [Data].

# Resumo

O projeto FazendaPro é uma solução agropecuária que visa facilitar a gestão de fazendas e a produção de leite. O sistema oferece uma interface intuitiva para gerenciar animais, pastagens e produção de leite, além de fornecer insights para tomada de decisão. Uma das principais funcionalidades é a monitoração de vacas em lactação, permitindo acompanhar a produção de leite e identificar possíveis problemas, assim como manter seu histórico, como genitora, filho, etc.

## 1. Introdução

- **Contexto**: o software se mostra útil no contexto pecuário, resolvendo problemas e facilitando o gerenciamento de animais.
- **Justificativa**: a ideia surgiu para resolver uma dor real de um fazendeiro do norte de Minas Gerais, que não tinha bons faturamentos na venda de seus animais, pois não tira o histórico de cada animal.
- **Objetivos**: o objetivo principal do projeto é criar um sistema que permita a gestão de fazendas e a produção de leite de forma eficiente e ágil, guardando histórico dos animais e calculando gastos e faturamento.

## 2. Descrição do Projeto

- **Tema do Projeto**: o tema do projeto é a gestão de fazendas, o histórico do gado e a produção de leite.
- **Problemas a Resolver**: o principal problema a ser resolvido é a garantia da valorização de um gado no mercado, por meio do seu histórico, desde o nascimento, genética, vacinas, alimentação entre outras informações. Além de oferecer um sistema baixo custo para produtores e fazendeiros que não tem acesso a tecnologias semelhantes por causa do altos preços do softwares existentes no mercado.
- **Limitações**: Delimitação dos problemas que o projeto não abordará.

## 3. Especificação Técnica

## Requisitos Funcionais (RF)

**1. RF01 - Acessar o Sistema**

- RF01.01 O sistema deve permitir que o usuário faça o login na plataforma com seus credenciais.
- RF01.02 O sistema deve validar as credenciais do usuário e conceder acesso apenas para os usuários que geraram o token.

**2. RF02 - Adicionar um Animal**

- RF02.01 sistema deve permitir que o usuário cadastre um novo animal no sistema.
- RF02.02 O sistema deve permitir incluir dados do animal como, no mínimo: identificação (nome e número do brinco), data de nascimento, genitora, filho (caso exista), raça, sexo e informações de saúde (vacinas).

**RF03 - Gerenciar o Animal**

- RF03.01 O sistema deve permitir que o usuário edite ou exclua as informações de um animal já cadastrado.
- RF03.02 O sistema deve oferecer a opção de exportar o histórico do animal em formato PDF.

**RF05 - Analisar Dashboards**

- RF05.01 O sistema deve fornecer dashboards com informações analíticas sobre os animais, como produção de leite, saúde geral, e tendências de desempenho.

**RF06 - Inserir Informações do Animal**

- RF06.01 O sistema deve permitir que o usuário insira informações adicionais sobre o animal, como registros de vacinas, alimentação, tratamentos ou eventos, como nascimento de filhotes.

**RF07 - Registrar Peso do Animal por Mês/Semana**

- RF07.01 O sistema deve permitir que o usuário registre o peso do animal em intervalos regulares (mensal ou semanal).
- RF07.02 O sistema deve armazenar esses registros para acompanhamento do desenvolvimento do animal.
- RF07.03 O sistema deve permitir a edição ou exclusão desses registros.

**RF08 - Mudar de Lote**

- RF08.01 O sistema deve mudar automaticamente o lote ao qual um animal pertence dependendo da sua produção de leite.

**RF09 - Definir Data de Prenhez**

- RF09.01 sistema deve permitir que o usuário registre a data de prenhez de uma vaca.
- RF09.02 sistema deve notificar o usuário (via WhatsApp) quando a data de prenhez estiver próxima do parto, 20 dias antes.

**RF10 - Vender o Animal**

- RF10.01 O sistema deve permitir que o usuário registre a venda de um animal.
- RF10.02 O sistema deve atualizar o status do animal para "vendido" e registrar a data da venda.
- RF10.03 O sistema deve oferecer a opção de exportar o histórico do animal em PDF no momento da venda.
- RF10.04 O sistema deve permitir verificar o histórico de todas as vendas dentro do módulo de vendas

**RF11 - Cadastrar Vacinas**

- RF11.01 O sistema deve permitir que o usuário cadastre a vacina para que depois ela seja vinculada ao animal
- RF11.02 O sistema deve permitir a pesquisa de vacinas por datas

**RF12 - Separar Módulo**

- RF12.01 O sistema deve organizar as informações através de módulos dentro de um menu lateral (Dashboard, Vacas, Fornecedores, Vendas, Estoque)

**RF13 - Sair da Plataforma**

- RF13.01 O sistema deve permitir que o usuário faça o logout da plataforma

## Requisitos Não Funcionais (RNF)

**1. RNF01 - Estilização**

- RNF01.01 A estilização da aplicação deve seguir os padrões de estilo do Figma
- RNF01.02 Para facilitar a estilização deve ser usado Tailwind ou outra biblioteca de CSS
- RNF01.03 Componentes padrões devem ser criados para seguir um padrão geral
- RNF01.04 As cores da aplicação devem apresentar-se de forma agradável

**2. RNF02 - Ferramentas**

- RNF02.01 Para o Frontend deve-se utilizar React com bibliotecas para facilitar o fetch das informações
- RNF02.02 Para o Backend será usado NestJS para serviços de autenticação e notificações, contudo para todo os resto será usado Go

**3. RNF03 - Idiomas**

- RNF03.01 Todo o desenvolvimento deve ser feito respeitando variáveis de idioma
- RNF02.02 O idioma principal será PT-BR, posteriormente pode ser implementado EN-US e ES-ES

**4. RNF04 - Mobile**

- RNF04.01 - O desenvolvimento deve respeitar os casos de mobile, respeitando um design responsivo e agradável

## Diagrama de Casos de uso

![Casos de uso](images/cases-of-use.drawio.png)

## Diagrama de Classes

![Diagrama de Classes](images/classesDiagram.drawio.png)

## Diagram de Relacionamento

![Diagram de Relacionamento](images/entityRelationshipDiagram.drawio.png)

## Diagrama de Estados

![Diagrama de Estados](images/stateDiagram.drawio.png)

### 3.2. Considerações de Design

**Visão Inicial da Arquitetura**: Foi decidido usar a arquitetura modular para o projeto. A arquitetura modular oferece um equilíbrio entre a simplicidade de um monolito e a flexibilidade dos microserviços. O NestJS facilita essa abordagem através de seus módulos bem definidos, permitindo escalabilidade sem a complexidade inicial dos microserviços. Esta escolha permite que o sistema cresça naturalmente, com a possibilidade de extrair módulos para microserviços no futuro (como é o caso das notificações no futuro).

**Padrões de Arquitetura**: A ideia seria usar uma arquitetura limpa como o DDD (Domain-Driven Design) com arquitetura hexagonal.

```bash
fazendapro-api/
├── api/                    # Definições de handlers e endpoints HTTP
│   ├── handlers/           # Funções que lidam com requisições HTTP
│   │   ├── user.go         # Ex.: handler para rotas de usuário
│   │   └── product.go      # Ex.: handler para rotas de produtos
│   └── middleware/         # Middlewares (ex.: autenticação JWT)
│       └── auth.go         # Middleware de validação de token
├── cmd/                    # Ponto de entrada do projeto
│   └── app/                # Aplicação principal
│       └── main.go         # Arquivo principal que inicializa o servidor
├── config/                 # Configurações (ex.: variáveis de ambiente)
│   └── config.go           # Carrega .env ou outras configs
├── internal/               # Código interno (não exposto para outros pacotes)
│   ├── models/             # Estruturas de dados (ex.: User, Product)
│   │   ├── user.go
│   │   └── product.go
│   ├── repository/         # Acesso a dados (ex.: banco de dados)
│   │   ├── user_repository.go
│   │   └── product_repository.go
│   └── service/            # Lógica de negócio
│       ├── user_service.go
│       └── product_service.go
├── pkg/                    # Código reutilizável (se necessário)
│   └── jwt/                # Funções utilitárias para JWT
│       └── jwt.go
├── scripts/                # Scripts úteis (ex.: para build ou deploy)
├── tests/                  # Testes unitários e de integração
│   ├── handlers/
│   ├── repository/
│   └── service/
├── .env                    # Variáveis de ambiente (ex.: JWT_SECRET)
├── go.mod                  # Definição do módulo Go
├── go.sum                  # Dependências
└── README.md               # Documentação do projeto
```

### Modelos C3:

```mermaid
classDiagram
    %% Diagrama de Contexto (C4 Nível 1)
    class FazendaPro {
        +Gerencia fazendas e produção de leite
    }
    class Fazendeiro {
        +Gerencia animais e vendas
    }
    class Sistema_de_Pagamento {
        +Processa transações (futuro)
    }
    class WhatsApp {
        +Envia notificações
    }
    Fazendeiro --> FazendaPro : Usa
    FazendaPro --> Sistema_de_Pagamento : Integra (futuro)
    FazendaPro --> WhatsApp : Envia notificações de prenhez

    %% Diagrama de Containers (C4 Nível 2)
    class WebApp {
        +Interface do usuário
        +React, TypeScript, Tailwind
    }
    class API_Server {
        +Gerencia lógica de negócios
        +NestJS (autenticação, notificações), Go (demais serviços)
    }
    class Banco_de_Dados {
        +Armazena dados
        +MySQL (JawsDB)
    }
    class Redis {
        +Cache em memória
        +Redis
    }
    class New_Relic {
        +Monitoramento de performance
        +APM
    }
    class Sentry {
        +Monitoramento de erros
        +Error Tracking
    }
    class Heroku {
        +Hospedagem e CI/CD
        +Plataforma PaaS
    }

    Fazendeiro --> WebApp : Acessa via navegador/mobile
    WebApp --> API_Server : Faz chamadas REST
    API_Server --> Banco_de_Dados : Lê/Escreve via TypeORM
    API_Server --> Redis : Cache de dados
    API_Server --> WhatsApp : Envia notificações
    API_Server --> New_Relic : Envia métricas
    API_Server --> Sentry : Reporta erros
    API_Server --> Heroku : Deploy via CI/CD

    %% Diagrama de Componentes (C4 Nível 3 - Foco na API_Server)
    class Autenticacao {
        +Gerencia login/logout
        +JWT, Bcrypt
    }
    class Gerenciador_de_Animais {
        +Cadastra, edita, exclui animais
        +Historico, PDF
    }
    class Gerenciador_de_Vacinas {
        +Cadastra, vincula vacinas
    }
    class Gerenciador_de_Lotes {
        +Move animais entre lotes
    }
    class Gerenciador_de_Prenhez {
        +Registra datas, notifica
    }
    class Gerenciador_de_Vendas {
        +Registra vendas, exporta PDF
    }
    class Dashboard {
        +Exibe analises
    }
    API_Server *--> Autenticacao
    API_Server *--> Gerenciador_de_Animais
    API_Server *--> Gerenciador_de_Vacinas
    API_Server *--> Gerenciador_de_Lotes
    API_Server *--> Gerenciador_de_Prenhez
    API_Server *--> Gerenciador_de_Vendas
    API_Server *--> Dashboard
    Autenticacao --> Banco_de_Dados : Valida credenciais
    Gerenciador_de_Animais --> Banco_de_Dados : Lê/Escreve dados de animais
    Gerenciador_de_Animais --> Redis : Cache de históricos
    Gerenciador_de_Vacinas --> Banco_de_Dados : Lê/Escreve vacinas
    Gerenciador_de_Lotes --> Banco_de_Dados : Atualiza lotes
    Gerenciador_de_Prenhez --> Banco_de_Dados : Registra prenhez
    Gerenciador_de_Prenhez --> WhatsApp : Envia notificações
    Gerenciador_de_Vendas --> Banco_de_Dados : Registra vendas
    Dashboard --> Banco_de_Dados : Lê dados analíticos
    Dashboard --> Redis : Cache de métricas
```

![Diagrama de Estados](images/architecture.png)

**Aplicação Web:** React.

**Api Server:** servidor NestJS em container no Heroku, funcionando como núcleo do sistema.

Inclui: Redis rodando no mesmo container para caching em memória.

**Armazenamento persistente de dados (MySql)**

**Interações:**
A Aplicação Web faz requisições HTTP (REST) ao API Server.
O API Server usa Redis internamente para caching e consulta o JawsDB MySQL via conexão SQL.

### 3.3. Stack Tecnológica

- **Linguagens de Programação**: Justificativa para a escolha de linguagens específicas.
- **Frameworks e Bibliotecas**:
  - React
  - Nest.js
  - TypeORM
  - JWT
  - Bcrypt
  - Express
  - Styled Components (em poucos casos)
  - React Router
  - React Hook Form
  - React Query
  - React Toastify
  - React Icons
  - Yup
  - Jest
  - Cypress
  - Docker
  - MySql
  - Docker
  - Docker Compose
- **Ferramentas de Desenvolvimento e Gestão de Projeto**: Para a gestão do projeto foi utilizado o Github Projects para criar as atividades. Algumas atividades já foram criadas e podem ser vistas neste [link](https://github.com/orgs/fazendapro/projects/1).

### 3.4. Considerações de Segurança

#### Autenticação e Autorização (as rotas serão protegidas)

- Credenciais expostas (senhas fracas ou vazamento de tokens).
  - Vai ser utilizado hash para senhas com bcrypt
- Falta de proteção contra ataques de força bruta.
  - Será usado limite de tentativas de login (rate limiting) com @nestjs/throttler

#### Exposição de Dados Sensíveis

- Vazamento de informações em respostas da API
  - Será usado DTO para retornar apenas o necessário pela API.
- Configurações inadequadas de CORS permitindo acesso não autorizado.
- Logs com informações sensíveis
  - O Heroku usa automaticamente o HTTPS para criptografar a comunicação.

#### Injeção de Código (SQL Injection, XSS, etc.)

- Consultas SQL no backend
  - será usado o TypeORM e não queries brutas
- Scripts maliciosos injetados no frontend React via entradas de usuário
  - será implementado Content Security Policy (CSP) no frontend para limitar fontes de scripts.

## 4. Próximos Passos

Descrição dos passos seguintes após a conclusão do documento, com uma visão geral do cronograma para Portfólio I e II.

## 5. Referências

Listagem de todas as fontes de pesquisa, frameworks, bibliotecas e ferramentas que serão utilizadas.

## 6. Apêndices (Opcionais)

Informações complementares, dados de suporte ou discussões detalhadas fora do corpo principal.

## 7. Avaliações de Professores

Adicionar três páginas no final do RFC para que os Professores escolhidos possam fazer suas considerações e assinatura:

- Considerações Professor/a:
- Considerações Professor/a:
- Considerações Professor/a:

escrever sobre:

- new relic
- redis para cache
- sentry
- integração ci/cd com Heroku
