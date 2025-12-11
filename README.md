Identity Service (Clean Architecture Template)
===============================================

Overview
--------
Lightweight Go service structured around clean architecture: domain is pure, use cases orchestrate behavior, transports map I/O, and infrastructure implements interfaces.

Project Structure
-----------------
```
cmd/api/                   # app entrypoint (wiring only)
internal/domain/           # entities + domain behavior (no deps)
internal/usecase/          # application logic; depends on interfaces
internal/repository/       # repository interfaces
internal/repository/inmemory/   # sample infra implementation
internal/transport/httptransport/ # HTTP handlers + router
internal/transport/dto/    # transport-facing DTOs
```

Flow: POST /users/register
--------------------------
1) HTTP handler decodes JSON into DTO and maps to use case input.  
2) Use case builds domain `User`, runs `Validate`, sets IDs/timestamps, persists via repository interface.  
3) Repository implementation stores data (in-memory sample).  
4) Use case returns output; handler maps to response DTO and writes JSON.

How it satisfies clean architecture
-----------------------------------
- Dependency rule: transport → usecase → repository interface → infra; domain has no outward deps.  
- Domain purity: `internal/domain` has only Go stdlib; validation (`User.Validate`) lives here.  
- Use case isolation: `internal/usecase` works with interfaces and domain types, unaware of HTTP/JSON.  
- Interface adapters: transport maps DTOs; repositories implement interfaces in `internal/repository/<impl>`.  
- Composition root: `cmd/api/main.go` wires concrete implementations.

Key Code References
-------------------
Handler mapping:
```
23:55:internal/transport/httptransport/user_handler.go
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    var input dto.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid Input"+err.Error(), http.StatusBadRequest)
        return
    }
    ucInput := &usecase.CreateUserInput{...}
    ucOutput, err := h.userUC.CreateUser(r.Context(), *ucInput)
    if err != nil {
        http.Error(w, "server error", http.StatusInternalServerError)
        return
    }
    resp := &dto.UserResponse{...}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(resp)
}
```

Use case (logic + repository interface):
```
12:65:internal/usecase/user_usecase.go
type UserUsecase struct { userRepo repository.UserRepository }
...
func (uc *UserUsecase) CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
    user := &domain.User{ ID: uuid.NewString(), ... }
    createdUser, err := uc.userRepo.CreateUser(ctx, user)
    if err != nil { return nil, err }
    return &CreateUserOutput{ ID: createdUser.ID, ... }, nil
}
```

Domain entity + validation:
```
8:26:internal/domain/user.go
type User struct { ID string; Email string; ... }
func (u *User) Validate() error {
    if u.Email == "" { return errors.New("email can't be empty") }
    return nil
}
```

Repository interface and in-memory adapter:
```
9:11:internal/repository/repository.go
type UserRepository interface {
    CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}
```
```
10:27:internal/repository/inmemory/user_repository.go
type InMemoryUserRepository struct { mu sync.RWMutex; users []*domain.User }
func (r *InMemoryUserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    r.users = append(r.users, user)
    return user, nil
}
```

Router and composition:
```
9:14:internal/transport/httptransport/router.go
func NewRouter(h *UserHandler) http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/users/register/", h.Create)
    return mux
}
```
```
11:22:cmd/api/main.go
func main() {
    userRepo := inmemory.NewInMemoryUserRepository()
    userUC := usecase.NewUserUsecase(userRepo)
    userHandler := httptransport.NewUserHandler(userUC)
    r := httptransport.NewRouter(userHandler)
    log.Println("🚀 Server running on :8080")
    http.ListenAndServe(":8080", r)
}
```

Run locally
-----------
- `go run ./cmd/api` (or `go run ./cmd/api/main.go`)

Sample request
--------------
```
curl -X POST http://localhost:8080/users/register/ \
  -H "Content-Type: application/json" \
  -d '{
    "cognito_sub": "sub-123",
    "email": "user@example.com",
    "phone": "+10000000000",
    "tenant_id": "t1",
    "is_active": true
  }'
```

Extending
---------
- Add new use cases in `internal/usecase`; keep inputs/outputs transport-agnostic.
- Add new transport adapters (e.g., gRPC) under `internal/transport/<adapter>`.
- Add new repo implementations (e.g., `repository/postgres`) implementing `UserRepository`.


