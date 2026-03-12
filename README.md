
# FitFlow
- Project Name: FitFlow
- Project description: FitFlow is a full-stack web application designed to manage the lifecycle of fitness class registration. The system serves as a bridge between fitness enthusiasts and studio offerings, providing a real-time interface for class discovery, enrollment management, and schedule coordination. From an engineering standpoint, the project focuses on high-performance concurrency, stateless authentication, and modular front-end architecture.A full-stack web application built with React and Go. This project demonstrates a fitness class registration website. The website presents different types of fitness classes, and allows users to register/cancel classes. Other functions like setting a reminder to register specific class. 
- Members:
  - Frontend: Ila Adhikari, Forrest Yan Sun
  - Backend: Yingzhu Chen, Qing Li


---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Installation](#installation)
- [Backend Test Demo](#backend-test-demo)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- User authentication (signup/login/logout)
- CRUD operations for [main resource, e.g., tasks, products, posts]
- Responsive design for desktop and mobile
- RESTful API integration between frontend and backend
- [Any other special feature]

---

## Tech Stack

**Frontend:**  
* **React 19 (TypeScript):** UI development with strict typing.
* **Vite:** High-performance build tool and dev server.
* **Tailwind CSS 4:** Modern utility-first styling for responsive design.
* **React Router 7:** Handling client-side navigation and protected routes.
* **Zustand:** Lightweight state management for authentication.
* **Axios:** Consuming RESTful APIs from the Go backend.
* **Lucide React:** Consistent and clean UI iconography.

**Backend:**  
- Go 
- SQLite
- JWT for authentication  

---

## Installation

1. Clone the repository:  
```bash
git clone https://github.com/your-username/project-name.git
```

## Backend Test Demo

Use these two files side by side during the demo:

- `Backend/service/class_service_test.go` for business-logic tests
- `Backend/routes/class_routes_test.go` for HTTP endpoint tests

Recommended demo flow:

1. Open the two test files side by side in VS Code.
2. Show that the service tests cover enroll, drop, list classes, and analytics logic.
3. Show that the route tests hit the real Gin endpoints and verify HTTP status codes plus JSON responses.
4. Run the tests from the `Backend` folder.

Commands:

```bash
go test ./service -v
go test ./routes -v
go test ./service ./routes -cover
```

Focused demo command:

```bash
go test ./service ./routes -run "Test(RegisterClass|DropClass|ListClassesPaged|GetUserAnalytics|RegisterClassEndpoint|DropClassEndpoint|ListClassesEndpoint|GetUserAnalyticsEndpoint)" -v
```
