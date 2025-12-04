backend/internal/api/handlers/animal.go


Define a constant instead of duplicating this literal "ID do animal inválido" 3 times.

Adaptability
Maintainability


4
High
design
+
Open
Gustavo Dias
Gustavo Dias
L185
6min effort
1 month ago
Code Smell
Critical
backend/internal/routes/routes.go


Define a constant instead of duplicating this literal "application/json" 3 times.

Adaptability
Maintainability


4
High
design
+
Open
Gustavo Dias
Gustavo Dias
L53
6min effort
23 hours ago
Code Smell
Critical


Define a constant instead of duplicating this literal "Content-Type" 3 times.

Adaptability
Maintainability


4
High
design
+
Open
Gustavo Dias
Gustavo Dias
L53
6min effort
23 hours ago
Code Smell
Critical
backend/internal/service/animal_service.go


Define a constant instead of duplicating this literal "animals:farm:%d" 4 times.

Adaptability
Maintainability


4
High
design
+
Open
Gustavo Dias
Gustavo Dias
L82
8min effort
23 hours ago
Code Smell
Critical


Define a constant instead of duplicating this literal "Erro ao invalidar cache (não crítico): %v" 3 times.

Adaptability
Maintainability


4
High
design
+
Open
Gustavo Dias
Gustavo Dias
L84
6min effort
23 hours ago
Code Smell
Critical
backend/internal/service/sale_service.go


Define a constant instead of duplicating this literal "animal not found" 3 times.

Adaptability
Maintainability


4
High
design
+
Open
Gustavo Dias
Gustavo Dias
L62
6min effort
23 hours ago
Code Smell
Critical


Define a constant instead of duplicating this literal "sale not found or does not belong to farm" 4 times.

Adaptability
Maintainability


4
High
design
+
Open
Gustavo Dias
Gustavo Dias
L197
8min effort
23 hours ago
Code Smell
Critical
frontend/src/components/services/axios/base.ts


Prefer `throw error` over `return Promise.reject(error)`.

Intentionality
Maintainability


4
High
async
confusing
...
+
Open
Gustavo Dias
Gustavo Dias
L54
5min effort
2 months ago
Code Smell
Major
frontend/src/contexts/__tests__/FarmContext.test.tsx


Refactor this code to not nest functions more than 4 levels deep.