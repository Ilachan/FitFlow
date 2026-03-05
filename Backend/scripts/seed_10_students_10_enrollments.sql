PRAGMA foreign_keys = ON;
BEGIN TRANSACTION;

-- 1) Ensure Student role exists (main.go also seeds this, but this keeps script standalone)
INSERT INTO Role (role_name)
SELECT 'Student'
WHERE NOT EXISTS (
  SELECT 1 FROM Role WHERE role_name = 'Student'
);

-- 2) Create 10 demo students (idempotent by unique email)
WITH student_role AS (
  SELECT id AS role_id FROM Role WHERE role_name = 'Student' LIMIT 1
)
INSERT INTO Student (name, email, password, avatar_url, role_id, created_at)
SELECT s.name, s.email, s.password, '', sr.role_id, datetime('now')
FROM (
  SELECT 'Demo Student 01' AS name, 'demo.student01@seed.local' AS email, '$2a$10$zzVv9WHqtRfcevUt1PncS.BXZ.68R0vOUtnFGpmdcAjtQl0KXXynS' AS password UNION ALL
  SELECT 'Demo Student 02', 'demo.student02@seed.local', '$2a$10$zzVv9WHqtRfcevUt1PncS.BXZ.68R0vOUtnFGpmdcAjtQl0KXXynS' UNION ALL
  SELECT 'Demo Student 03', 'demo.student03@seed.local', '$2a$10$zzVv9WHqtRfcevUt1PncS.BXZ.68R0vOUtnFGpmdcAjtQl0KXXynS' UNION ALL
  SELECT 'Demo Student 04', 'demo.student04@seed.local', '$2a$10$zzVv9WHqtRfcevUt1PncS.BXZ.68R0vOUtnFGpmdcAjtQl0KXXynS' UNION ALL
  SELECT 'Demo Student 05', 'demo.student05@seed.local', '$2a$10$zzVv9WHqtRfcevUt1PncS.BXZ.68R0vOUtnFGpmdcAjtQl0KXXynS' UNION ALL
  SELECT 'Demo Student 06', 'demo.student06@seed.local', '$2a$10$zzVv9WHqtRfcevUt1PncS.BXZ.68R0vOUtnFGpmdcAjtQl0KXXynS' UNION ALL
  SELECT 'Demo Student 07', 'demo.student07@seed.local', '$2a$10$zzVv9WHqtRfcevUt1PncS.BXZ.68R0vOUtnFGpmdcAjtQl0KXXynS' UNION ALL
  SELECT 'Demo Student 08', 'demo.student08@seed.local', '$2a$10$zzVv9WHqtRfcevUt1PncS.BXZ.68R0vOUtnFGpmdcAjtQl0KXXynS' UNION ALL
  SELECT 'Demo Student 09', 'demo.student09@seed.local', '$2a$10$zzVv9WHqtRfcevUt1PncS.BXZ.68R0vOUtnFGpmdcAjtQl0KXXynS' UNION ALL
  SELECT 'Demo Student 10', 'demo.student10@seed.local', '$2a$10$zzVv9WHqtRfcevUt1PncS.BXZ.68R0vOUtnFGpmdcAjtQl0KXXynS'
) s
CROSS JOIN student_role sr
WHERE NOT EXISTS (
  SELECT 1 FROM Student st WHERE st.email = s.email
);

-- 3) Create 10 enrollments (1 per demo student), spreading across available courses.
--    Enroll dates are staggered over recent days so analytics has a timeline.
WITH demo_students AS (
  SELECT id AS student_id,
         ROW_NUMBER() OVER (ORDER BY id) AS rn
  FROM Student
  WHERE email LIKE 'demo.student%@seed.local'
  ORDER BY id
  LIMIT 10
),
course_pool AS (
  SELECT id AS course_id,
         ROW_NUMBER() OVER (ORDER BY id) AS rn,
         (SELECT COUNT(*) FROM Course) AS total
  FROM Course
),
mapped AS (
  SELECT ds.student_id,
         (
           SELECT cp.course_id
           FROM course_pool cp
           WHERE cp.rn = ((ds.rn - 1) % (SELECT total FROM course_pool LIMIT 1)) + 1
           LIMIT 1
         ) AS course_id,
         datetime('now', printf('-%d day', ds.rn - 1), '09:00:00') AS enroll_time
  FROM demo_students ds
)
INSERT INTO StudentEnrollment (student_id, course_id, status, enroll_time)
SELECT m.student_id, m.course_id, 'registered', m.enroll_time
FROM mapped m
WHERE m.course_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1
    FROM StudentEnrollment se
    WHERE se.student_id = m.student_id
      AND se.course_id = m.course_id
  );

COMMIT;

-- Optional checks:
-- SELECT id, name, email FROM Student WHERE email LIKE 'demo.student%@seed.local' ORDER BY id;
-- SELECT student_id, course_id, status, enroll_time FROM StudentEnrollment
--   WHERE student_id IN (SELECT id FROM Student WHERE email LIKE 'demo.student%@seed.local')
--   ORDER BY student_id;
