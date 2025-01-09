.PHONY: frontend
frontend:
	@cd frontend && npm run dev & open http://localhost:5173

.PHONY: backend
backend:
	@cd backend/cmd/api && go run .