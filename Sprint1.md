# User Stories

### User Story 1 – User Registration (Student)

As a new student user, I want to create an account to use the website.

#### Acceptance Criteria & Tasks

**Backend Issues**
- [BE-1] Design the user database table with role support (id, username, email, password_hash, role).
- [BE-2] Implement signup logic for students with email uniqueness validation and bcrypt password hashing.

**Frontend Issues**
- [FE-1] Create Angular AuthService signup method for student registration.
- [FE-2] Build student registration UI using reactive forms with validation.
- 
### User Story 3 – Enroll in a Class
As a student, I want to enroll in a class, so that I can participate in the course.

#### Acceptance Criteria & Tasks
**Backend Issues**

[BE-17] Implement POST /api/register with transaction support.
[BE-18] Check max capacity before enrollment.
**Frontend Issue**

[FE-14] Implement "Book Now" button.

## 2. Ability to search for different classes
I can check classes based on different types of fitness formats.

### User Story 4 – Drop a Class

As a student, I want to cancel my enrollment, so that I can adjust my schedule.

#### Acceptance Criteria & Tasks

**Backend Issue**
- [BE-20] Implement DELETE /api/register/:class_id endpoint.

**Frontend Issue**
- [FE-18] Add "Cancel Registration" button in dashboard.


### User Story 2 – View Course List

As a user, I want to view a list of available courses, so that I can discover classes.

#### Acceptance Criteria & Tasks

**Backend Issue**
- [BE-14] Implement GET /api/courses endpoint.

**Frontend Issue**
- [FE-11] Display course list in the UI.
- 

### User Story 6 – View My Enrolled Classes (Student)

As a student, I want to see the classes I am enrolled in, so that I can manage my schedule.

#### Acceptance Criteria & Tasks

**Backend Issue**
- [BE-21] Implement GET /api/my-courses endpoint.

**Frontend Issue**
- [FE-19] Display enrolled courses in dashboard.


#### Acceptance Criteria & Tasks

**Backend Issue**
- [BE-22] Implement API endpoint to return available spots for a class.

**Frontend Issue**
- [FE-20] Display the remaining spots for each class in the course listing UI.

##6. Ability to "like" a class
I want to save classes that I am interested.

