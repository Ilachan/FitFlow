### Sprint 3 Report: FitFlow Team
**Team Members**:  

Frontend: Forrest Yan Sun, Ila Adhikari  
Backend: Qing Li, Yingzhu Chen  

**Project Links**:  

🔗 GitHub Repository: https://github.com/Ilachan/FitFlow  
📺 Demo Video: https://youtu.be/fK1popveTx0 


---

### User Stories  

#### **User Story 1 – Student Registers for a Course**
**As a student, I want to register for a course within the registration window, so that I can participate in the class.**  

**Acceptance Criteria & Tasks**  
- **Acceptance Criteria**: 
    - Students can register for courses only within the registration time window.
    - Registration should check if the course exists, the student exists, and avoid duplicate enrollments.
    - Time conflicts with other enrolled courses should be avoided.
    - Successful registration should return feedback and redirect the user to the course list.

- [FE-1] Create a registration form allowing students to input `course_id`.  
- [FE-2] Handle API calls for course registration and display feedback messages.  
- [FE-3] Implement "Register Now" functionality with dynamic buttons on the course detail page.  
- [BE-1] Validate course and user data during registration.  
- [BE-2] Implement a new POST endpoint: `/classes/register`.  
- [BE-3] Write logic to return conflict/errors for duplicate enrollments or time clashes.  

---

#### **User Story 2 – Course Time Conflict Detection**
**As a student, I want the system to detect if the course timing overlaps with my existing schedule, so that I can avoid conflicts.**  

**Acceptance Criteria & Tasks**  
- **Acceptance Criteria**: 
    - Time conflicts must be checked before enrollment is finalized.
    - For enrollments with conflicts, the system should provide detailed feedback.  

- [FE-4] Show error messages and explain the conflicting courses if time overlaps occur.  
- [BE-4] Add backend logic for time conflict detection via weekday normalization and range overlap checks.  

---

#### **User Story 3 – Generate Course Sessions**
**As a course admin, I want to auto-generate future course sessions weekly, so that I can efficiently manage recurring courses.**  

**Acceptance Criteria & Tasks**  
- **Acceptance Criteria**: 
    - Admin can schedule course sessions automatically by specifying recurrence patterns (e.g., weekly).
    - The backend updates the database with all new sessions in bulk.

- [BE-5] Build and expose the `GenerateClassSessions` logic using a new backend function.  
- [BE-6] Add error handling for database bulk inserts or validation mismatches.  

---

#### **User Story 4 – Instructor Retrieves Their Courses**
**As an instructor, I want to view the list of courses I am responsible for, so that I can clearly understand my teaching load.**  

**Acceptance Criteria & Tasks**  
- **Acceptance Criteria**: 
    - Instructors can view all courses they are assigned to.
    - Only courses for the authenticated instructor are returned.  

- [FE-7] Create a course list page dynamically populated with API data.  
- [FE-8] Add navigation to the course list page from the dashboard.  
- [BE-7] Implement the `GET /instructor/courses` endpoint.  
- [BE-8] Ensure role validation so only instructors can access their data.  

---

#### **User Story 5 – Instructor Views Enrolled Students**
**As an instructor, I want to view all students enrolled in a specific course, so that I can monitor their enrollment progress.**  

**Acceptance Criteria & Tasks**  
- **Acceptance Criteria**: 
    - Displays student data (e.g., name, email, enrollment status) for the selected course.
    - API request must validate the instructor's permission.  

- [FE-9] Add a "View Students" section under the course detail page.  
- [FE-10] Handle pagination or large data loads in the student list table.  
- [BE-9] Create the `GET /instructor/courses/:id/enrollments` endpoint and add permission checks.  

---

#### **User Story 6 – Instructor Marks Student Attendance**
**As an instructor, I want to update the attendance status of students in my course, so that I can accurately record their participation.**  

**Acceptance Criteria & Tasks**  
- **Acceptance Criteria**: 
    - Attendance options include: "enrolled", "attended", and "missed".
    - Updates are valid only for current course participants.  

- [FE-11] Add an inline attendance update form for each student row.  
- [FE-12] Include filtering by attendance status in the UI.  
- [BE-10] Implement `PATCH /instructor/courses/:id/enrollments` endpoint with validation logic.  
- [BE-11] Write persistence logic for updating attendance in the database.  

---

### Frontend Unit Tests

- `src/lib/api.test.ts`
- `src/lib/validation.test.ts`
- `src/store/authStore.test.ts`
- `src/pages/Login.test.tsx`
- `src/pages/Register.test.tsx`
- `src/pages/Profile.test.tsx`
- `src/pages/Browse.test.tsx` (Sprint 3: enrollment window display logic)
- `src/pages/InstructorDashboard.test.tsx` (Sprint 3: attendance and walk-in flows)
- `src/pages/InstructorProfile.test.tsx` (Sprint 3: bio and availability quick actions)

### Backend Unit Tests

- Overall status: PASS
- Packages with tests: `routes`, `service`
- Packages without tests: `api`, `dao`, `db`, `model`, `script/seeder/demo_data`, `tmpdbmaint`
- Runtime: about 3.7 seconds for the full suite

Notable coverage shown in the suite:

- class registration and dropping
- class listing and analytics
- instructor course listing and enrollment updates
- manager user and enrollment management
- super manager invite code creation

The negative-path tests print expected `record not found` database logs, but the assertions still pass.

### Backend API Documentation Updates

## Auth

- `POST /auth/login` - log in and return a JWT plus `role_id`
- `GET /auth/profile` - get the authenticated user's profile
- `PUT /auth/profile` - update the authenticated user's profile
- `POST /auth/manager/register` - register a manager with an invite code
- `POST /auth/manager/invite-codes` - create a manager invite code, super manager only
- `POST /auth/roles/assign` - assign a role, super manager only

## Users

- `DELETE /users/:id` - delete a user
- `GET /users/:id/enrollments` - list the user's enrolled classes
- `GET /users/:id/analytics?range=7d|1m|3m` - view user analytics

## Classes

- `GET /classes` - list classes with pagination
- `GET /classes/:id` - get a single class
- `GET /classes/:id/enrollments` - list class enrollments, manager only
- `POST /classes/register` - enroll in a class
- `POST /classes/drop` - drop a class
- `POST /classes` - create a class, manager only
- `PUT /classes/:id` - update a class, manager only
- `DELETE /classes/:id` - delete a class, manager only

## Instructor

- `GET /instructor/courses` - list instructor courses
- `GET /instructor/courses/:id/enrollments` - list enrollments for one instructor course
- `POST /instructor/courses/:id/enrollments` - enroll a user in an instructor course
- `PATCH /instructor/courses/:id/enrollments` - update enrollment status as `attended` or `missed`

## Manager

- `GET /manager/users?page=1&limit=20` - list users
- `GET /manager/users/:id/enrollments` - list one user's enrollments
- `POST /manager/users/:id/enrollments` - add a user enrollment
- `DELETE /manager/users/:id/enrollments/:course_id` - delete a user enrollment

## Notes

- Authentication uses the `Authorization: Bearer <token>` header.
- Public endpoints are limited to class listing, class details, registration, and login/register flows.
- Role checks are enforced in the handlers, so the same route can return `401`, `403`, `404`, `409`, or `201` depending on the request state.
