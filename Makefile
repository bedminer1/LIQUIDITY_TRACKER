.PHONY: frontend
frontend:
	@cd frontend && npm run dev & open http://localhost:5173

.PHONY: backend
backend:
	@cd backend/cmd/api && go run .

.PHONY: init_frontend
init_frontend:
	@cd frontend && npm install
	@echo "Frontend dependencies installed successfully."

.PHONY: init_microservice
init_microservice:
	@cd ai_microservice && python3.11 -m venv .venv
	@cd ai_microservice && source .venv/bin/activate && .venv/bin/pip install -r requirements.txt
	@echo "Python virtual environment set up and dependencies installed."

.PHONY: microservice
microservice:
	@cd ai_microservice && . .venv/bin/activate && python microservice.py

.PHONY: init_all
init_all: init_frontend init_microservice
	@echo "All services initialized!"

.PHONY: run_all
run_all:
	@make -j3 frontend backend microservice