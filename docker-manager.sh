#!/bin/bash

# Order Tracking System - Docker Management Script

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Functions
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Main menu
show_menu() {
    clear
    echo -e "${BLUE}╔══════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  Order Tracking System - Docker Manager ║${NC}"
    echo -e "${BLUE}╚══════════════════════════════════════════╝${NC}"
    echo ""
    echo "1. 🚀 Start all services"
    echo "2. 🛑 Stop all services"
    echo "3. 🔄 Restart all services"
    echo "4. 🏗️  Rebuild and start"
    echo "5. 📊 View service status"
    echo "6. 📋 View logs (all services)"
    echo "7. 🔍 View logs (specific service)"
    echo "8. 🗑️  Clean up (stop and remove volumes)"
    echo "9. 🧪 Insert sample data"
    echo "10. 📦 Access PostgreSQL shell"
    echo "11. 🐰 Open RabbitMQ Management UI"
    echo "12. 🔍 Open Elasticsearch"
    echo "0. ❌ Exit"
    echo ""
}

# Start services
start_services() {
    print_header "Starting All Services"
    docker compose up -d
    print_success "All services started!"
    echo ""
    print_warning "Wait 10-15 seconds for services to be ready..."
    sleep 3
    docker compose ps
}

# Stop services
stop_services() {
    print_header "Stopping All Services"
    docker compose down
    print_success "All services stopped!"
}

# Restart services
restart_services() {
    print_header "Restarting All Services"
    docker compose restart
    print_success "All services restarted!"
}

# Rebuild and start
rebuild_services() {
    print_header "Rebuilding and Starting Services"
    docker compose up --build -d
    print_success "Services rebuilt and started!"
}

# View status
view_status() {
    print_header "Service Status"
    docker compose ps
}

# View all logs
view_logs() {
    print_header "Viewing Logs (Press Ctrl+C to exit)"
    docker compose logs -f
}

# View specific service logs
view_service_logs() {
    echo ""
    echo "Select service:"
    echo "1. Order Service"
    echo "2. Stock Service"
    echo "3. Analytics Service"
    echo "4. PostgreSQL"
    echo "5. RabbitMQ"
    echo "6. Elasticsearch"
    echo ""
    read -p "Enter choice [1-6]: " service_choice
    
    case $service_choice in
        1) SERVICE="order-service" ;;
        2) SERVICE="stock-service" ;;
        3) SERVICE="analytics-service" ;;
        4) SERVICE="postgres" ;;
        5) SERVICE="rabbitmq" ;;
        6) SERVICE="elasticsearch" ;;
        *) print_error "Invalid choice"; return ;;
    esac
    
    print_header "Viewing $SERVICE Logs (Press Ctrl+C to exit)"
    docker compose logs -f $SERVICE
}

# Clean up
cleanup() {
    print_header "Cleaning Up"
    print_warning "This will stop all services and remove volumes!"
    read -p "Are you sure? (y/N): " confirm
    if [[ $confirm == [yY] ]]; then
        docker compose down -v
        print_success "Cleanup complete!"
    else
        print_warning "Cleanup cancelled"
    fi
}

# Insert sample data
insert_sample_data() {
    print_header "Inserting Sample Data"
    
    # Check if postgres is running
    if ! docker compose ps | grep -q "order-tracking-postgres.*running"; then
        print_error "PostgreSQL is not running. Please start services first."
        return
    fi
    
    print_warning "Waiting for PostgreSQL to be ready..."
    sleep 2
    
    docker exec -i order-tracking-postgres psql -U admin -d orders_db <<EOF
INSERT INTO customers (id, name) VALUES
('CUST-001', 'Md. Rakibul Kabir'),
('CUST-002', 'Arka Das'),
('CUST-003', 'Morshed Alam'),
('CUST-004', 'Emran Ahmed Emon')
ON CONFLICT (id) DO NOTHING;

INSERT INTO stocks (sku, quantity) VALUES
('LAPTOP-001', 50),
('MOUSE-001', 100),
('KEYBOARD-001', 75)
ON CONFLICT (sku) DO UPDATE SET quantity = EXCLUDED.quantity;

SELECT 'Sample data inserted successfully!' as status;
EOF
    
    print_success "Sample data inserted!"
}

# Access PostgreSQL shell
access_postgres() {
    print_header "PostgreSQL Shell"
    print_warning "Type \\q to exit"
    echo ""
    docker exec -it order-tracking-postgres psql -U admin -d orders_db
}

# Open RabbitMQ Management UI
open_rabbitmq() {
    print_header "Opening RabbitMQ Management UI"
    echo "URL: http://localhost:15673"
    echo "Username: admin"
    echo "Password: admin"
    echo ""
    if command -v xdg-open &> /dev/null; then
        xdg-open "http://localhost:15673"
    elif command -v open &> /dev/null; then
        open "http://localhost:15673"
    else
        print_warning "Please open http://localhost:15673 in your browser"
    fi
}

# Open Elasticsearch
open_elasticsearch() {
    print_header "Opening Elasticsearch"
    echo "URL: http://localhost:9201"
    echo ""
    curl -s "http://localhost:9201" | jq || print_error "Elasticsearch is not running or jq is not installed"
}

# Main loop
while true; do
    show_menu
    read -p "Enter your choice [0-12]: " choice
    
    case $choice in
        1) start_services ;;
        2) stop_services ;;
        3) restart_services ;;
        4) rebuild_services ;;
        5) view_status ;;
        6) view_logs ;;
        7) view_service_logs ;;
        8) cleanup ;;
        9) insert_sample_data ;;
        10) access_postgres ;;
        11) open_rabbitmq ;;
        12) open_elasticsearch ;;
        0) 
            print_success "Goodbye!"
            exit 0
            ;;
        *)
            print_error "Invalid choice. Please try again."
            sleep 2
            ;;
    esac
    
    echo ""
    read -p "Press Enter to continue..."
done
