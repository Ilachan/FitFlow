# ✅ **1) Student Registers for a Course**
**User Story**:  
As a student,  
I want to register for a course within the registration window,  
so that I can participate in the class.

**Endpoint**: `POST /classes/register`  
**Functionality**: The student registers for a course, **with status written as `enrolled`** (replaces the original `registered`).

**Request Body (JSON)**:
```json
{
  "course_id": 7
}
```

[FE]:
- Create a registration form for students to input the `course_id`.
- Call the API with the JSON payload and handle the response (e.g., display success or failure messages).
- Add a "Register Now" button for students on the course detail page.

[BE]:
- Validate whether the `course_id` exists and ensure the course is open for registration.
- Check if the user meets registration criteria (e.g., user exists, no time conflicts).
- Return appropriate status codes and messages for success or failure.

---

# ✅ **2) Course Time Conflict Detection**
**User Story**:  
As a student,  
when I register for a course,  
I want the system to detect if the course timing overlaps with my existing schedule,  
so that I can avoid conflicts.

**Added Logic**:
- Added weekday normalization and time range overlap checks to prevent enrollment with conflicting time slots.

[FE]:
- Display a user-friendly error message if a time conflict occurs and suggest students adjust their schedule.

[BE]:
- Integrate conflict detection logic into the course registration API: compare new registration time slots with existing enrolled courses.
- Return an error response with details of the conflicting course if applicable.

---

# ✅ **3) Generate Course Sessions**
**User Story**:  
As a course admin,  
I want to auto-generate future course sessions weekly,  
so that I can efficiently manage recurring courses.

**New Functionality**:
- Support a `GenerateClassSessions` function to generate and upsert future sessions based on a predefined weekly schedule.

[FE]:
- Provide a "Generate Sessions" button on the course management page to trigger the backend process.
- Show the result of the session generation (success or failure).

[BE]:
- Implement the `GenerateClassSessions` logic to accept course plans and create recurring sessions in the database.
- Ensure a batch operation inserts or updates the generated sessions.

---

# ✅ **4) Instructor Retrieves Their Course List**
**User Story**:  
As an instructor,  
I want to view the list of courses I am responsible for,  
so that I can clearly understand my teaching load.

**Endpoint**: `GET /instructor/courses`  
**Functionality**: Returns the list of courses the instructor is responsible for.

**Request Parameters**: None (GET request has no body).  

**Header**:
```
Authorization: Bearer <token>
```

[FE]:
- Create an interface for instructors to view their courses in a list format.
- Call this API and dynamically render the courses on the frontend.
- Handle session expiration by prompting instructors to re-login when the token is invalid.

[BE]:
- Verify the instructor's identity and ensure only their courses are returned.
- Query the database for courses associated with the logged-in instructor.
- Optimize query performance to handle large datasets (e.g., add pagination support).

---

# ✅ **5) Instructor Views Course Students**
**User Story**:  
As an instructor,  
I want to view all students enrolled in a specific course,  
so that I can monitor their enrollment progress.

**Endpoint**: `GET /instructor/courses/:id/enrollments`  
**Functionality**: Returns all enrolled students for a specific course.

**Request Parameters**: None (GET request has no body).  

**Header**:
```
Authorization: Bearer <token>
```

[FE]:
- Add a "View Students" interface to the course details page that displays a list of enrolled students.
- Support search functionality to allow instructors to quickly find specific students.
- Display a clear error message if the server fails to respond.

[BE]:
- Fetch the list of students for the specified course from the database.
- Verify the instructor's permission to access this data.
- Return relevant student information (e.g., name, enrollment status).

---

# ✅ **6) Instructor Marks Student Attendance**
**User Story**:  
As an instructor,  
I want to update the attendance status of students in my course,  
so that I can accurately record their participation.

**Endpoint**: `PATCH /instructor/courses/:id/enrollments`  
**Functionality**: Updates the attendance status of an enrolled student.

**Request Body (JSON)**:
```json
{
  "user_id": 15,
  "status": "attended"
}
```

**Allowed Status Values**:
```
enrolled | attended | missed
```

**Header**:
```
Authorization: Bearer <token>
```

[FE]:
- Provide a form or inline buttons for instructors to update attendance statuses dynamically.
- Call the API with the updated `user_id` and `status` values, then refresh the displayed status accordingly.
- Display success or error messages based on the API response.

[BE]:
- Validate that the instructor has permission to update attendance for the specified course.
- Verify that the student is enrolled in the course.
- Persist the updated status in the database and return a success or failure message.

---

# ✅ **7) SuperManager Assigns User Roles**
**User Story**:  
As a SuperManager,  
I want to assign roles (e.g., Instructor, Manager) to users,  
so that I can define their permissions accordingly.

**Endpoint**: `POST /auth/roles/assign`  
**Functionality**: SuperManager assigns a role to a user.

**Request Body (JSON)**:
```json
{
  "user_id": 20,
  "role_name": "Instructor"
}
```

**Header**:
```
Authorization: Bearer <token>
```

[FE]:
- Provide a role-assignment interface for SuperManagers with a dropdown to select roles and a search bar for users.
- Call this API upon submitting the role data and handle success or error responses.

[BE]:
- Ensure that only SuperManagers (Role == 2) can access this endpoint.
- Validate both the `user_id` and the `role_name` before updating the database.
- Avoid duplicate role assignments by checking existing user-role mappings.

---

# ✅ **8) Manager Retrieves User List (With Pagination)**
**User Story**:  
As a Manager,  
I want to retrieve a paginated list of users,  
so that I can manage user data more efficiently.

**Endpoint**: `GET /manager/users?page=1&limit=20`  
**Functionality**: Returns a paginated user list.

[FE]:
- Render user information into a dynamic list on the user management page.
- Add "Previous" and "Next" buttons for page navigation and update the interface as users paginate.
- Show a spinner or loading indicator while fetching data.

[BE]:
- Allow Managers to query for a paginated list of all system users.
- Implement pagination logic by using `page` and `limit` parameters to determine database offsets.
- Return a paginated response format, such as:
```json
{
  "users": [{ "id": 1, "name": "Alice" }, ...],
  "total_count": 100,
  "current_page": 1
}
```

---

# ✅ **9) Manager Views User Enrollment Data**
**User Story**:  
As a Manager,  
I want to view the list of courses and states for a specific user,  
so that I can track their academic progress.

**Endpoint**: `GET /manager/users/:id/enrollments`  
**Functionality**: Returns the list of courses and enrollment states for the specified user.

[FE]:
- Create a dedicated interface to display a user’s course enrollments in a table.
- Support filtering and sorting by enrollment status, course name, etc.
- Show error messages when no data is found or if the API request fails.

[BE]:
- Validate Manager permissions and ensure the user exists.
- Fetch the user’s courses and enrollment statuses from the database.
- Respond with detailed course and state information.

---

# ✅ **10) Manager Adds a Course for the User**
**User Story**:  
As a Manager,  
I want to assign a course to a user,  
so that they can enroll in specific classes as required.

**Endpoint**: `POST /manager/users/:id/enrollments`  
**Functionality**: Adds a course to the specified user.

**Request Body (JSON)**:
```json
{
  "course_id": 7
}
```

[FE]:
- Create a modal or form for Managers to select course(s) for a user.
- Call this API to assign the course, and refresh the user’s enrollment table upon success.

[BE]:
- Verify that the course and user data are valid.
- Ensure the user is not already enrolled in the course.
- Insert the new enrollment into the database.

---

# ✅ **11) Manager Removes a User from a Course**
**User Story**:  
As a Manager,  
I want to remove a course from a user’s enrollments,  
so that their schedule can be adjusted.

**Endpoint**: `DELETE /manager/users/:id/enrollments/:course_id`  
**Functionality**: Removes the specified course from the user’s enrollments.

[FE]:
- Add a "Remove" button on the user’s course table.
- Provide a confirmation dialog before calling the API to delete the enrollment.

[BE]:
- Validate Manager permissions and ensure the specified course and user exist.
- Remove the corresponding enrollment from the database.

