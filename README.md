# Authora 🔐

**Authora** is the authentication and authorization service in the **ose-micro** ecosystem.  
It provides secure login, tenant management, role-based access control (RBAC), and permission handling for multi-tenant
systems.

---

## 🚀 Features

- **User Authentication**
    - Login with JWT (access + refresh tokens)
    - Password hashing (bcrypt/argon2)
    - Token revocation & refresh
- **Tenant Management**
    - Each organization (org) is mapped to a tenant
    - Tenant auto-created when an org is registered
    - Users can belong to multiple tenants
- **Role-Based Access Control (RBAC)**
    - Default roles: `admin`, `member`
    - Assign multiple roles per user per tenant
    - Extendable for custom roles
- **Permission Management**
    - JSON/Map-based permission definitions
    - Example:
      ```json
      {
        "campaign:create": "allow",
        "campaign:delete": "deny",
        "report:view": "allow"
      }
      ```
- **Event-Driven**
    - Publishes events when tenants, users, or roles are created/updated
    - Works seamlessly with other `ose-micro` services

---

## 🏗 Architecture

            +-------------------+
            |    Authora API    |
            +---------+---------+
                      |
    +-----------------+-----------------+
    |                                   |
    +-------v-------+ +-------v-------+
    | Authentication| | Authorization |
    +---------------+ +---------------+
    | JWT / Refresh | | RBAC / ACL |
    | Password Hash | | Role Assign |
    +---------------+ +---------------+
    |
    +---------v---------+
    | Tenant Service |
    +-------------------+
    |
    +---------v---------+
    | Event Bus (NATS/Redis) |
    +-------------------+

---

## 📦 Installation

```bash
# Clone the repo
git clone https://github.com/ose-micro/authora.git
cd authora

# Install dependencies
go mod tidy

# Generate protobufs
make proto

```

## ⚡ Usage

Run Authora

```go
go run cmd/authora/main.go
```

## Example gRPC Flow

- User signs up → creates User in Authora
- Org created → Authora publishes TenantCreated event
- User assigned roles → stored as UserRoleAssignment

- Services query Authora to validate permissions

## 🔑 Permission Format

Permissions are stored as map[string]string, example:

```json
{
  "campaign:create": "allow",
  "campaign:delete": "deny",
  "report:view": "allow"
}

```

## 🛠 Development

- Protobufs live in: proto/ose/micro/authora/
- Generated code in: internal/interface/grpc/gen/
- Business logic in: internal/core/
- DB Layer in: internal/infra/db/

## 📚 Related Services

- ose-cqrs - Event sourcing & messaging
- ose-nats - NATS integration
- ose-mongo - Database wrapper

## 📝 License

MIT © ose-micro