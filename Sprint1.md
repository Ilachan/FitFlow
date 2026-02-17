# User Stories

### User Story 1 – User Registration (Student)

As a new student user, I want to create an account to use the website, so that I can browse and enroll in fitness classes.

#### Acceptance Criteria & Tasks

**Acceptance Criteria:**
- Students can register with email, username, and password
- Email must be unique and follow valid email format
- Password must be at least 8 characters with at least one uppercase letter, one lowercase letter, and one number
- Users receive immediate feedback on validation errors
- Upon successful registration, users are automatically logged in and redirected to the course list page
- User role is automatically set to "student" by default
- Confirmation message is displayed after successful registration

**Backend Issues**
- [BE-1] Design the user database table with role support (id, username, email, password_hash, role, created_at, updated_at).
- [BE-2] Implement signup logic for students with email uniqueness validation and bcrypt password hashing.
- [BE-3] Add password strength validation (minimum 8 characters, uppercase, lowercase, number).
- [BE-4] Implement POST /api/signup endpoint returning JWT token and user data.
- [BE-5] Add email format validation using regex.
- [BE-6] Set default role to "student" for new registrations.

**Frontend Issues**
- [FE-1] Create Angular AuthService signup method for student registration.
- [FE-2] Build student registration UI using reactive forms with validation.
- [FE-3] Add real-time validation feedback for email format and password strength.
- [FE-4] Display password strength indicator (weak/medium/strong).
- [FE-5] Implement "Show/Hide Password" toggle button.
- [FE-6] Add confirmation dialog for successful registration.
- [FE-7] Implement auto-login after successful registration with JWT token storage.
- [FE-8] Add redirect to course list page after registration.

### User Story 2 – User Login

As a registered user, I want to log into my account, so that I can access my personalized features and enrolled classes.

#### Acceptance Criteria & Tasks

**Backend Issues**
- [BE-7] Implement login endpoint POST /api/login with email and password validation.
- [BE-8] Generate and return JWT token upon successful authentication.
- [BE-9] Implement password verification using bcrypt.

**Frontend Issues**
- [FE-9] Create Angular AuthService login method with JWT token storage.
- [FE-10] Build login UI form with email and password fields.
- [FE-11] Implement error handling for invalid credentials.
- [FE-12] Add redirect to dashboard upon successful login.

### User Story 3 – View User Profile (Student)

As a student user, I want to view my profile information, so that I can see my account details and personal information.

#### Acceptance Criteria & Tasks

**Acceptance Criteria:**
- Students can view their complete profile information including:
  - Name
  - Avatar/profile picture
  - Date of birth
  - Gender
  - Phone number
  - Address
  - Email and registration date
- Profile page displays user's enrolled class count
- Users must be authenticated to access their profile
- Profile information is displayed in a clear and organized layout
- Default avatar is displayed if user hasn't uploaded one

**Backend Issues**
- [BE-10] Extend user database table to include name, avatar_url, date_of_birth, gender, phone_number, address fields.
- [BE-11] Implement GET /api/user/profile endpoint with JWT authentication.
- [BE-12] Return complete user data excluding password_hash.
- [BE-13] Include enrolled classes count in the response.
- [BE-14] Add middleware to verify JWT token for protected routes.

**Frontend Issues**
- [FE-13] Create ProfileService to fetch user profile data.
- [FE-14] Build user profile page component with all profile fields displayed.
- [FE-15] Display avatar image with fallback to default avatar.
- [FE-16] Format date of birth and phone number for display.
- [FE-17] Add navigation link to profile page in header/menu.
- [FE-18] Implement loading state while fetching profile data.

### User Story 4 – Update User Profile (Student)

As a student user, I want to update my profile information, so that I can keep my personal details current and accurate.

#### Acceptance Criteria & Tasks

**Acceptance Criteria:**
- Students can update the following profile fields:
  - Name
  - Avatar URL (profile picture)
  - Date of birth
  - Gender (dropdown: Male/Female/Other/Prefer not to say)
  - Phone number (with format validation)
  - Address
- Phone number must follow valid format: +1 (XXX) XXX-XXXX
- Date of birth must be a valid date and user must be at least 13 years old
- Avatar URL must be a valid URL format
- Users receive confirmation message after successful update
- Form displays current profile information as default values
- Users can cancel changes and revert to original values
- All fields are optional except name

**Backend Issues**
- [BE-15] Implement PUT /api/user/profile endpoint with JWT authentication.
- [BE-16] Add validation for phone number format.
- [BE-17] Add validation for date of birth
- [BE-18] Add validation for avatar URL format.
- [BE-19] Validate gender field against allowed values.
- [BE-20] Return updated user data (excluding password_hash).
- [BE-21] Add transaction support to ensure data integrity.

**Frontend Issues**
- [FE-19] Create profile edit form with reactive validation.
- [FE-20] Pre-populate form fields with current user data.
- [FE-21] Add avatar image preview with upload URL input field.
- [FE-22] Implement date picker for date of birth with age validation.
- [FE-23] Add dropdown for gender selection.
- [FE-24] Implement phone number input with format mask (+1 (XXX) XXX-XXXX).
- [FE-25] Add real-time validation for phone number and avatar URL.
- [FE-26] Display validation errors for each field.
- [FE-27] Add confirmation dialog before saving changes.
- [FE-28] Show success message after update and refresh profile data.
- [FE-29] Add "Cancel" button to discard changes and return to view mode.
- [FE-30] Implement error handling for validation failures.

### User Story 5 – Enroll in a Class

As a student, I want to enroll in a class, so that I can participate in the course.

#### Acceptance Criteria & Tasks

**Backend Issues**
- [BE-22] Implement POST /api/register with transaction support.
- [BE-23] Check max capacity before enrollment.

**Frontend Issue**
- [FE-31] Implement "Book Now" button.

## 2. Ability to search for different classes
I can check classes based on different types of fitness formats.

### User Story 6 – Drop a Class

As a student, I want to cancel my enrollment, so that I can adjust my schedule.

#### Acceptance Criteria & Tasks

**Backend Issue**
- [BE-24] Implement DELETE /api/register/:class_id endpoint.

**Frontend Issue**
- [FE-32] Add "Cancel Registration" button in dashboard.

### User Story 7 – View Course List

As a user, I want to view a list of available courses, so that I can discover classes.

#### Acceptance Criteria & Tasks

**Backend Issue**
- [BE-25] Implement GET /api/courses endpoint.

**Frontend Issue**
- [FE-33] Display course list in the UI.

### User Story 8 – View My Enrolled Classes (Student)

As a student, I want to see the classes I am enrolled in, so that I can manage my schedule.

#### Acceptance Criteria & Tasks

**Backend Issue**
- [BE-26] Implement GET /api/my-courses endpoint.

**Frontend Issue**
- [FE-34] Display enrolled courses in dashboard.

#### Acceptance Criteria & Tasks

**Backend Issue**
- [BE-27] Implement API endpoint to return available spots for a class.

**Frontend Issue**
- [FE-35] Display the remaining spots for each class in the course listing UI.

## 6. Ability to "like" a class
I want to save classes that I am interested.
