#!/bin/bash

# FazendaPro Database Migration Runner
# Description: Runs database migrations in the correct order
# Usage: ./run_migrations.sh [database_name] [mysql_user] [mysql_password] [action]

DB_NAME="${1:-fazendapro}"
DB_USER="${2:-root}"
DB_PASS="${3:-}"
ACTION="${4:-migrate}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== FazendaPro Database Migration Runner ===${NC}"
echo -e "Database: $DB_NAME"
echo -e "User: $DB_USER"
echo -e "Action: $ACTION"
echo ""

# Function to execute SQL file
execute_migration() {
    local file="$1"
    local description="$2"
    
    echo -n "Running: $description... "
    
    if [ -z "$DB_PASS" ]; then
        MYSQL_CMD="mysql -u $DB_USER $DB_NAME"
    else
        MYSQL_CMD="mysql -u $DB_USER -p$DB_PASS $DB_NAME"
    fi
    
    if $MYSQL_CMD < "$file" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
        return 0
    else
        echo -e "${RED}✗${NC}"
        echo "Error in file: $file"
        return 1
    fi
}

case $ACTION in
    "migrate"|"up")
        echo -e "${YELLOW}=== Running Migrations ===${NC}"
        
        # Check if database exists, create if not
        echo -n "Ensuring database exists... "
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
        
        # Run migrations in order
        MIGRATIONS=(
            "database/migrations/001_create_users_table.sql:Create Users Table"
            "database/migrations/002_create_lots_table.sql:Create Lots Table"
            "database/migrations/003_create_animals_table.sql:Create Animals Table"
            "database/migrations/004_create_vaccines_table.sql:Create Vaccines Table"
            "database/migrations/005_create_animal_vaccines_table.sql:Create Animal Vaccines Table"
            "database/migrations/006_create_weight_records_table.sql:Create Weight Records Table"
            "database/migrations/007_create_pregnancy_records_table.sql:Create Pregnancy Records Table"
            "database/migrations/008_create_suppliers_table.sql:Create Suppliers Table"
            "database/migrations/009_create_sales_table.sql:Create Sales Table"
        )
        
        for migration in "${MIGRATIONS[@]}"; do
            IFS=':' read -r file description <<< "$migration"
            if [ -f "$file" ]; then
                execute_migration "$file" "$description"
                if [ $? -ne 0 ]; then
                    echo -e "${RED}Migration failed. Stopping.${NC}"
                    exit 1
                fi
            else
                echo -e "${RED}Migration file not found: $file${NC}"
                exit 1
            fi
        done
        
        echo -e "${GREEN}All migrations completed successfully!${NC}"
        ;;
        
    "seed")
        echo -e "${YELLOW}=== Loading Sample Data ===${NC}"
        execute_migration "database/seeds/sample_data.sql" "Loading Sample Data"
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}Sample data loaded successfully!${NC}"
        else
            echo -e "${RED}Failed to load sample data.${NC}"
            exit 1
        fi
        ;;
        
    "reset"|"fresh")
        echo -e "${YELLOW}=== Resetting Database ===${NC}"
        echo -n "Dropping and recreating database... "
        
        if [ -z "$DB_PASS" ]; then
            mysql -u $DB_USER -e "DROP DATABASE IF EXISTS $DB_NAME; CREATE DATABASE $DB_NAME;" > /dev/null 2>&1
        else
            mysql -u $DB_USER -p$DB_PASS -e "DROP DATABASE IF EXISTS $DB_NAME; CREATE DATABASE $DB_NAME;" > /dev/null 2>&1
        fi
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓${NC}"
        else
            echo -e "${RED}✗${NC}"
            exit 1
        fi
        
        # Run migrations after reset
        $0 "$DB_NAME" "$DB_USER" "$DB_PASS" "migrate"
        ;;
        
    "complete")
        echo -e "${YELLOW}=== Using Complete Schema ===${NC}"
        execute_migration "database/schema/complete_schema.sql" "Loading Complete Schema"
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}Complete schema loaded successfully!${NC}"
        else
            echo -e "${RED}Failed to load complete schema.${NC}"
            exit 1
        fi
        ;;
        
    "status")
        echo -e "${YELLOW}=== Database Status ===${NC}"
        
        if [ -z "$DB_PASS" ]; then
            MYSQL_CMD="mysql -u $DB_USER $DB_NAME"
        else
            MYSQL_CMD="mysql -u $DB_USER -p$DB_PASS $DB_NAME"
        fi
        
        echo "Checking tables..."
        echo "$MYSQL_CMD -e 'SHOW TABLES;'" | bash
        
        echo ""
        echo "Table row counts:"
        $MYSQL_CMD -e "
        SELECT 'users' as table_name, COUNT(*) as row_count FROM users
        UNION ALL SELECT 'animals', COUNT(*) FROM animals
        UNION ALL SELECT 'lots', COUNT(*) FROM lots
        UNION ALL SELECT 'vaccines', COUNT(*) FROM vaccines
        UNION ALL SELECT 'animal_vaccines', COUNT(*) FROM animal_vaccines
        UNION ALL SELECT 'weight_records', COUNT(*) FROM weight_records
        UNION ALL SELECT 'pregnancy_records', COUNT(*) FROM pregnancy_records
        UNION ALL SELECT 'suppliers', COUNT(*) FROM suppliers
        UNION ALL SELECT 'sales', COUNT(*) FROM sales;" 2>/dev/null
        ;;
        
    *)
        echo -e "${RED}Invalid action: $ACTION${NC}"
        echo ""
        echo "Usage: $0 [database_name] [mysql_user] [mysql_password] [action]"
        echo ""
        echo "Available actions:"
        echo "  migrate  - Run all migrations in order (default)"
        echo "  seed     - Load sample data"
        echo "  reset    - Drop and recreate database, then migrate"
        echo "  complete - Load complete schema in one go"
        echo "  status   - Show database status"
        echo ""
        echo "Examples:"
        echo "  $0 fazendapro root '' migrate"
        echo "  $0 fazendapro_test root password123 complete"
        echo "  $0 fazendapro root '' seed"
        exit 1
        ;;
esac