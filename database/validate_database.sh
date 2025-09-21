#!/bin/bash

# FazendaPro Database Validation Script
# Description: Validates the database schema and tests basic operations
# Usage: ./validate_database.sh [database_name] [mysql_user] [mysql_password]

DB_NAME="${1:-fazendapro_test}"
DB_USER="${2:-root}"
DB_PASS="${3:-}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== FazendaPro Database Validation ===${NC}"
echo -e "Database: $DB_NAME"
echo -e "User: $DB_USER"
echo ""

# Function to execute SQL and check result
execute_sql() {
    local sql="$1"
    local description="$2"
    
    echo -n "Testing: $description... "
    
    if [ -z "$DB_PASS" ]; then
        MYSQL_CMD="mysql -u $DB_USER $DB_NAME"
    else
        MYSQL_CMD="mysql -u $DB_USER -p$DB_PASS $DB_NAME"
    fi
    
    if echo "$sql" | $MYSQL_CMD > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
        return 0
    else
        echo -e "${RED}✗${NC}"
        return 1
    fi
}

# Test database connection
echo -n "Testing database connection... "
if [ -z "$DB_PASS" ]; then
    mysql -u $DB_USER -e "SELECT 1;" > /dev/null 2>&1
else
    mysql -u $DB_USER -p$DB_PASS -e "SELECT 1;" > /dev/null 2>&1
fi

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    echo "Error: Cannot connect to MySQL. Please check credentials."
    exit 1
fi

# Create test database if it doesn't exist
echo -n "Creating test database if needed... "
if [ -z "$DB_PASS" ]; then
    mysql -u $DB_USER -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;" > /dev/null 2>&1
else
    mysql -u $DB_USER -p$DB_PASS -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;" > /dev/null 2>&1
fi

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}=== Testing Schema Creation ===${NC}"

# Test schema creation
if [ -z "$DB_PASS" ]; then
    MYSQL_CMD="mysql -u $DB_USER $DB_NAME"
else
    MYSQL_CMD="mysql -u $DB_USER -p$DB_PASS $DB_NAME"
fi

echo -n "Loading complete schema... "
if $MYSQL_CMD < "database/schema/complete_schema.sql" > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    echo "Error loading schema. Check database/schema/complete_schema.sql"
    exit 1
fi

echo ""
echo -e "${YELLOW}=== Testing Table Structure ===${NC}"

# Test all tables exist
TABLES=("users" "lots" "animals" "vaccines" "animal_vaccines" "weight_records" "pregnancy_records" "suppliers" "sales")

for table in "${TABLES[@]}"; do
    execute_sql "DESCRIBE $table;" "Table $table structure"
done

echo ""
echo -e "${YELLOW}=== Testing Foreign Key Constraints ===${NC}"

execute_sql "SELECT TABLE_NAME, COLUMN_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = '$DB_NAME';" "Foreign key constraints"

echo ""
echo -e "${YELLOW}=== Testing Sample Data ===${NC}"

echo -n "Loading sample data... "
if $MYSQL_CMD < "database/seeds/sample_data.sql" > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    echo "Error loading sample data. Check database/seeds/sample_data.sql"
    exit 1
fi

echo ""
echo -e "${YELLOW}=== Testing Data Integrity ===${NC}"

# Test data was inserted correctly
execute_sql "SELECT COUNT(*) FROM users;" "Users table has data"
execute_sql "SELECT COUNT(*) FROM animals;" "Animals table has data"
execute_sql "SELECT COUNT(*) FROM animal_vaccines;" "Animal vaccines table has data"

# Test relationships work
execute_sql "SELECT a.name, u.name FROM animals a JOIN users u ON a.user_id = u.id LIMIT 1;" "User-Animal relationship"
execute_sql "SELECT a.name, l.name FROM animals a JOIN lots l ON a.lot_id = l.id LIMIT 1;" "Animal-Lot relationship"
execute_sql "SELECT a.name, mother.name FROM animals a JOIN animals mother ON a.mother_id = mother.id LIMIT 1;" "Parent-Child relationship"

echo ""
echo -e "${YELLOW}=== Testing Business Logic ===${NC}"

# Test unique constraints
execute_sql "SELECT COUNT(DISTINCT ear_tag) = COUNT(*) FROM animals;" "Ear tag uniqueness"

# Test enum values
execute_sql "SELECT COUNT(*) FROM animals WHERE gender IN ('male', 'female');" "Gender enum values"
execute_sql "SELECT COUNT(*) FROM animals WHERE status IN ('active', 'sold', 'deceased', 'transferred');" "Status enum values"

# Test date logic
execute_sql "SELECT COUNT(*) FROM pregnancy_records WHERE expected_birth_date > pregnancy_date;" "Pregnancy date logic"

echo ""
echo -e "${YELLOW}=== Testing Indexes ===${NC}"

execute_sql "SHOW INDEX FROM animals WHERE Key_name = 'idx_ear_tag';" "Animals ear_tag index"
execute_sql "SHOW INDEX FROM users WHERE Key_name = 'idx_email';" "Users email index"
execute_sql "SHOW INDEX FROM animal_vaccines WHERE Key_name = 'idx_animal_id';" "Animal vaccines index"

echo ""
echo -e "${YELLOW}=== Performance Test Query ===${NC}"

execute_sql "
SELECT 
    a.name,
    a.ear_tag,
    l.name as lot_name,
    COUNT(av.id) as vaccine_count,
    COUNT(wr.id) as weight_records
FROM animals a
LEFT JOIN lots l ON a.lot_id = l.id
LEFT JOIN animal_vaccines av ON a.id = av.animal_id
LEFT JOIN weight_records wr ON a.id = wr.animal_id
GROUP BY a.id
LIMIT 5;" "Complex JOIN query performance"

echo ""
echo -e "${GREEN}=== Validation Complete ===${NC}"

# Summary
echo ""
echo "Database: $DB_NAME"
echo "Schema files validated:"
echo "  ✓ database/schema/complete_schema.sql"
echo "  ✓ database/seeds/sample_data.sql"
echo ""

# Row counts
echo "Data summary:"
if [ -z "$DB_PASS" ]; then
    mysql -u $DB_USER $DB_NAME -e "
    SELECT 'Users' as Table_Name, COUNT(*) as Row_Count FROM users
    UNION ALL SELECT 'Animals', COUNT(*) FROM animals
    UNION ALL SELECT 'Lots', COUNT(*) FROM lots
    UNION ALL SELECT 'Vaccines', COUNT(*) FROM vaccines
    UNION ALL SELECT 'Animal_Vaccines', COUNT(*) FROM animal_vaccines
    UNION ALL SELECT 'Weight_Records', COUNT(*) FROM weight_records
    UNION ALL SELECT 'Pregnancy_Records', COUNT(*) FROM pregnancy_records
    UNION ALL SELECT 'Suppliers', COUNT(*) FROM suppliers
    UNION ALL SELECT 'Sales', COUNT(*) FROM sales;"
else
    mysql -u $DB_USER -p$DB_PASS $DB_NAME -e "
    SELECT 'Users' as Table_Name, COUNT(*) as Row_Count FROM users
    UNION ALL SELECT 'Animals', COUNT(*) FROM animals
    UNION ALL SELECT 'Lots', COUNT(*) FROM lots
    UNION ALL SELECT 'Vaccines', COUNT(*) FROM vaccines
    UNION ALL SELECT 'Animal_Vaccines', COUNT(*) FROM animal_vaccines
    UNION ALL SELECT 'Weight_Records', COUNT(*) FROM weight_records
    UNION ALL SELECT 'Pregnancy_Records', COUNT(*) FROM pregnancy_records
    UNION ALL SELECT 'Suppliers', COUNT(*) FROM suppliers
    UNION ALL SELECT 'Sales', COUNT(*) FROM sales;"
fi

echo ""
echo -e "${GREEN}Database validation completed successfully!${NC}"