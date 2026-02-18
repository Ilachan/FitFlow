# Sprint 1 Report: FitFlow Team

**Team Members:**
* **Frontend:** Forrest Yan Sun, Ila Adhikari
* **Backend:** Qing Li, Yingzhu Chen

**Project Links:**
* 🔗 **GitHub Repository:** [https://github.com/Ilachan/FitFlow](https://github.com/Ilachan/FitFlow)
* 📺 **Frontend Demo Video:** [https://youtu.be/GINi2uiHeqY](https://youtu.be/GINi2uiHeqY)
* 📺 **Backend Demo Video:** [https://youtu.be/yC_YCglcZHM](https://youtu.be/yC_YCglcZHM)

---

# User Stories
## successfully completed
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
- [BE-24] Check if enrollment existed.

**Frontend Issue**
- [FE-31] Implement "Book Now" button.


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

### User Story 8 – View My Enrolled Classes

As a student, I want to see the classes I am enrolled in, so that I can manage my schedule.

#### Acceptance Criteria & Tasks

**Backend Issue**
- [BE-26] Implement GET /api/my-courses endpoint.

**Frontend Issue**
- [FE-34] Display enrolled courses in dashboard.
- 
### User Story 9 – View spots left for a class
#### Acceptance Criteria & Tasks

**Backend Issue**
- [BE-27] Implement API endpoint to return available spots for a class.

**Frontend Issue**
- [FE-35] Display the remaining spots for each class in the course listing UI.

## Didn't Completed User Stories (Sprint 1 - Removed)
The following user stories were deprecated during Sprint 1 development. Based on Agile development principles and our subsequent team meetings, we determined that these features do not align with current user needs or are lower priority.
---

### ~~User Story X1 – Social Media Login Integration~~
**Status**: ❌ Deprecated  
**Reason**: User feedback indicated a preference for traditional email registration due to privacy concerns with social media authorization. The team decided to prioritize core functionality first.

As a new user, I want to register and login using my social media accounts (Google, Facebook, Apple), so that I can quickly access the platform without creating a new password.

#### Acceptance Criteria & Tasks

**Acceptance Criteria:**
- Users can register using Google, Facebook, or Apple accounts
- OAuth 2.0 authentication integration
- Automatically create user profile with social media data
- Link social accounts to existing email accounts

**Backend Issues**
- [BE-X1] Integrate OAuth 2.0 for Google authentication.
- [BE-X2] Integrate OAuth 2.0 for Facebook authentication.
- [BE-X3] Integrate OAuth 2.0 for Apple Sign-In.
- [BE-X4] Handle social media token validation and refresh.
- [BE-X5] Link social accounts to existing user accounts.

**Frontend Issues**
- [FE-X1] Add Google login button with OAuth flow.
- [FE-X2] Add Facebook login button with OAuth flow.
- [FE-X3] Add Apple Sign-In button.
- [FE-X4] Handle OAuth callback and token storage.

---

### ~~User Story X2 – Friend System & Social Network~~
**Status**: ❌ Deprecated  
**Reason**: User interviews revealed that students are more focused on the courses themselves rather than social features. This functionality adds system complexity but is expected to have very low usage.

As a student user, I want to add other students as friends and see their enrolled classes, so that I can join classes with my friends.

#### Acceptance Criteria & Tasks

**Acceptance Criteria:**
- Users can search for other students by username or email
- Send and accept friend requests
- View friends' enrolled classes on their profile
- Receive notifications for friend requests
- Friends list displays with avatars and online status

**Backend Issues**
- [BE-X6] Create friendships table (user_id, friend_id, status, created_at).
- [BE-X7] Implement POST /api/friends/request endpoint.
- [BE-X8] Implement PUT /api/friends/accept/:id endpoint.
- [BE-X9] Implement DELETE /api/friends/:id endpoint.
- [BE-X10] Implement GET /api/friends endpoint.
- [BE-X11] Implement GET /api/friends/:id/classes endpoint.
- [BE-X12] Add privacy settings for profile visibility.

**Frontend Issues**
- [FE-X5] Build friend search functionality.
- [FE-X6] Create friend request UI with accept/decline buttons.
- [FE-X7] Display friends list in user dashboard.
- [FE-X8] Show friend's enrolled classes on their profile.
- [FE-X9] Add friend request notifications.

---

### ~~User Story X3 – Class Review & Rating System~~
**Status**: ❌ Deprecated  
**Reason**: After communicating with instructors and administrators, they expressed concerns that negative reviews could impact course promotion. The decision was made to use an internal feedback system instead, not publicly visible to students.

As a student user, I want to rate and review classes I've completed, so that I can share my experience and help other students make informed decisions.

#### Acceptance Criteria & Tasks

**Acceptance Criteria:**
- Students can only review classes they have completed
- Rating scale: 1-5 stars
- Text review (optional, 500 character limit)
- Reviews display average rating on course list
- Students can edit or delete their own reviews
- Display reviewer name and date

**Backend Issues**
- [BE-X13] Create reviews table (id, user_id, class_id, rating, comment, created_at).
- [BE-X14] Implement POST /api/classes/:id/reviews endpoint.
- [BE-X15] Implement PUT /api/reviews/:id endpoint.
- [BE-X16] Implement DELETE /api/reviews/:id endpoint.
- [BE-X17] Implement GET /api/classes/:id/reviews endpoint.
- [BE-X18] Calculate and cache average rating for each class.
- [BE-X19] Add validation: user must have completed the class.

**Frontend Issues**
- [FE-X10] Create review submission form with star rating component.
- [FE-X11] Display average rating on course cards.
- [FE-X12] Build reviews section on class detail page.
- [FE-X13] Add edit/delete buttons for user's own reviews.
- [FE-X14] Implement character counter for review text.

---

### ~~User Story X4 – Advanced Class Filtering & Preferences~~
**Status**: ❌ Deprecated  
**Reason**: The current course catalog is relatively small (approximately 20-30 classes), and simple search and categorization are sufficient. Advanced filtering features are overly complex and would actually reduce user experience.

As a student user, I want to set my fitness preferences and use advanced filters, so that I can quickly find classes that match my interests and schedule.

#### Acceptance Criteria & Tasks

**Acceptance Criteria:**
- Save user preferences: intensity level, preferred time slots, favorite instructors
- Filter classes by: date range, time of day, duration, difficulty level, instructor, location
- Apply multiple filters simultaneously
- Save custom filter presets
- Sort results by: popularity, rating, upcoming time, difficulty

**Backend Issues**
- [BE-X20] Create user_preferences table.
- [BE-X21] Implement PUT /api/user/preferences endpoint.
- [BE-X22] Extend GET /api/courses with query parameters for all filters.
- [BE-X23] Implement filter preset saving (saved_filters table).
- [BE-X24] Add sorting logic for multiple criteria.
- [BE-X25] Optimize database queries with proper indexing.

**Frontend Issues**
- [FE-X15] Build advanced filter panel with all filter options.
- [FE-X16] Create user preferences settings page.
- [FE-X17] Implement filter preset save/load functionality.
- [FE-X18] Add sorting dropdown in course list.
- [FE-X19] Display active filters with clear buttons.
- [FE-X20] Implement filter state persistence in URL.

---

### ~~User Story X5 – Waiting List for Full Classes~~
**Status**: ❌ Deprecated  
**Reason**: The development team determined that implementation complexity is high, requiring automatic notification and position management systems. User feedback indicated they prefer to directly browse classes at other time slots.

As a student user, I want to join a waiting list for full classes, so that I can automatically enroll if a spot becomes available.

#### Acceptance Criteria & Tasks

**Acceptance Criteria:**
- Students can join waiting list when class is at capacity
- Waiting list shows user's position in queue
- Automatic enrollment when spot opens (FIFO order)
- Email/in-app notification when enrolled from waiting list
- Users can leave waiting list at any time
- Waiting list expires 24 hours before class start time

**Backend Issues**
- [BE-X26] Create waiting_list table (id, user_id, class_id, position, created_at).
- [BE-X27] Implement POST /api/classes/:id/waitlist endpoint.
- [BE-X28] Implement DELETE /api/waitlist/:id endpoint.
- [BE-X29] Create background job to process waiting list when spots open.
- [BE-X30] Implement notification system for automatic enrollment.
- [BE-X31] Add logic to expire waiting list entries.
- [BE-X32] Update position numbers when users leave waiting list.

**Frontend Issues**
- [FE-X21] Add "Join Waiting List" button for full classes.
- [FE-X22] Display waiting list position in user dashboard.
- [FE-X23] Add "Leave Waiting List" button.
- [FE-X24] Show waiting list count on class cards.
- [FE-X25] Display notification for automatic enrollment.

---

### ~~User Story X6 – Class Recommendation Engine~~
**Status**: ❌ Deprecated  
**Reason**: Requires substantial user data and machine learning algorithms. Too advanced for MVP stage, and team resources are insufficient. Decided to defer to future iterations.

As a student user, I want to receive personalized class recommendations based on my history and preferences, so that I can discover new classes that match my interests.

#### Acceptance Criteria & Tasks

**Acceptance Criteria:**
- System analyzes user's enrollment history and completed classes
- Recommends similar classes based on type, instructor, time preferences
- "Recommended for You" section on dashboard
- Recommendations update after each class completion
- Option to dismiss or hide recommendations

**Backend Issues**
- [BE-X33] Design recommendation algorithm based on user behavior.
- [BE-X34] Implement GET /api/recommendations endpoint.
- [BE-X35] Create background job to calculate recommendations.
- [BE-X36] Store recommendation scores in cache (Redis).
- [BE-X37] Track user interactions with recommendations.
- [BE-X38] Implement collaborative filtering logic.

**Frontend Issues**
- [FE-X26] Build "Recommended for You" component on dashboard.
- [FE-X27] Display recommendation cards with reasoning text.
- [FE-X28] Add dismiss button for recommendations.
- [FE-X29] Implement carousel for multiple recommendations.
- [FE-X30] Track clicks on recommended classes for analytics.

