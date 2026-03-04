-- =============================================================
-- stu_tracker Sample Data Init Script
-- Assumes organization_id = 1 and schema/tables already exist
-- Run order respects FK dependencies
-- =============================================================

-- ---------------------------------------------------------------
-- 0. ORGANIZATION (prerequisite for all FK references)
-- ---------------------------------------------------------------
INSERT INTO stu_tracker.Organization (id, name)
VALUES (1, 'Bright Futures Learning')
ON CONFLICT (id) DO NOTHING;


-- ---------------------------------------------------------------
-- 1. ADMIN ROOT
-- ---------------------------------------------------------------
INSERT INTO stu_tracker.Admin_root (id, first_name, last_name, email, organization_id)
VALUES (1, 'Sarah', 'Mitchell', 'sarah.mitchell@brightfutures.org', 1)
ON CONFLICT (id) DO NOTHING;


-- ---------------------------------------------------------------
-- 2. DISTRICT  (1 district, 2 locations within it)
-- ---------------------------------------------------------------
INSERT INTO stu_tracker.District (id, name, city, state, region, admin_id, organization_id)
VALUES (
    1,
    'Riverside Unified School District',
    'Riverside',
    'CA',
    'Southern California',
    1,
    1
)
ON CONFLICT (id) DO NOTHING;


-- ---------------------------------------------------------------
-- 3. LOCATIONS  (both belong to district 1)
-- ---------------------------------------------------------------
INSERT INTO stu_tracker.Locations (id, name, district_id, organization_id, address, city, state, zip_code)
VALUES
    (1, 'Riverside Elementary Center', 1, 1, '100 Oak Street',    'Riverside', 'CA', '92501'),
    (2, 'Riverside Middle Academy',   1, 1, '250 Maple Avenue',   'Riverside', 'CA', '92503')
ON CONFLICT (id) DO NOTHING;


-- ---------------------------------------------------------------
-- 4. PROGRAMS
-- ---------------------------------------------------------------
INSERT INTO stu_tracker.Programs (id, program_name, organization_id)
VALUES
    (1, 'After-School Enrichment', 1),
    (2, 'Summer Bridge Program',   1)
ON CONFLICT (id) DO NOTHING;


-- ---------------------------------------------------------------
-- 5. LOCATION_PROGRAMS  (both locations enrolled in both programs)
-- ---------------------------------------------------------------
INSERT INTO stu_tracker.Location_programs (location_id, program_id, organization_id)
VALUES
    (1, 1, 1),
    (1, 2, 1),
    (2, 1, 1),
    (2, 2, 1)
ON CONFLICT DO NOTHING;


-- ---------------------------------------------------------------
-- 6. SUBJECTS
-- ---------------------------------------------------------------
INSERT INTO stu_tracker.Subjects (id, title, description, organization_id)
VALUES
    (1, 'Mathematics',       'Arithmetic, algebra, and problem-solving fundamentals', 1),
    (2, 'Reading & Literacy','Comprehension, fluency, and writing skills',            1)
ON CONFLICT (id) DO NOTHING;


-- ---------------------------------------------------------------
-- 7. TUTORS  (both assigned to Location 1 / After-School Enrichment)
--    password_hash is a bcrypt placeholder — replace before production
-- ---------------------------------------------------------------
INSERT INTO stu_tracker.Tutors
    (id, organization_id, location_id, email, password_hash, first_name, last_name, active)
VALUES
    (1, 1, 1, 'james.carter@brightfutures.org',
     '$2b$12$placeholderHashForJamesCarter000000000000000',
     'James', 'Carter', TRUE),
    (2, 1, 1, 'maria.lopez@brightfutures.org',
     '$2b$12$placeholderHashForMariaLopez0000000000000000',
     'Maria', 'Lopez',  TRUE)
ON CONFLICT (id) DO NOTHING;


-- ---------------------------------------------------------------
-- 8. TUTOR_LOCATIONS  (both tutors ↔ Location 1, org 1)
-- ---------------------------------------------------------------
INSERT INTO stu_tracker.Tutor_locations (tutor_id, location_id, organization_id, attendance_link)
VALUES
    (1, 1, 1, 'https://attend.brightfutures.org/james-carter-loc1'),
    (2, 1, 1, 'https://attend.brightfutures.org/maria-lopez-loc1')
ON CONFLICT DO NOTHING;


-- ---------------------------------------------------------------
-- 9. STUDENTS  (sample students at Location 1)
-- ---------------------------------------------------------------
INSERT INTO stu_tracker.Students
    (id, first_name, last_name, grade, location_id)
VALUES
    (1, 'Alex',    'Turner',  3, 1),
    (2, 'Brianna', 'Nguyen',  4, 1),
    (3, 'Carlos',  'Reyes',   3, 1),
    (4, 'Diana',   'Park',    5, 1)
ON CONFLICT (id) DO NOTHING;


-- ---------------------------------------------------------------
-- 10. SESSIONS
--     Each session ties: location, program, subject, tutor, date
-- ---------------------------------------------------------------
INSERT INTO stu_tracker.Sessions
    (id, organization_id, location_id, program_id, subject_id, tutor_id,
     session_date, start_time, end_time, notes)
VALUES
    -- James Carter · Math · After-School Enrichment · Location 1
    (1, 1, 1, 1, 1, 1,
     '2025-03-03', '15:00', '16:00',
     'Introduction to fractions'),

    -- Maria Lopez · Reading · After-School Enrichment · Location 1
    (2, 1, 1, 1, 2, 2,
     '2025-03-03', '15:00', '16:00',
     'Guided reading – chapter books'),

    -- James Carter · Math · After-School Enrichment · Location 1 (next week)
    (3, 1, 1, 1, 1, 1,
     '2025-03-10', '15:00', '16:00',
     'Multiplying fractions'),

    -- Maria Lopez · Reading · Summer Bridge · Location 2
    (4, 1, 2, 2, 2, 2,
     '2025-06-16', '09:00', '10:00',
     'Vocabulary building workshop')
ON CONFLICT (id) DO NOTHING;


-- ---------------------------------------------------------------
-- 11. SESSION_STUDENTS  (attendance join table)
-- ---------------------------------------------------------------
INSERT INTO stu_tracker.Session_students (session_id, student_id, attended)
VALUES
    (1, 1, TRUE),
    (1, 2, TRUE),
    (1, 3, FALSE),
    (2, 1, TRUE),
    (2, 4, TRUE),
    (3, 1, TRUE),
    (3, 2, TRUE),
    (3, 3, TRUE),
    (4, 4, TRUE)
ON CONFLICT DO NOTHING;


-- =============================================================
-- Done.  Sequence counters — bump past seeded IDs so future
-- INSERTs don't collide.
-- =============================================================
SELECT setval(pg_get_serial_sequence('stu_tracker.District',        'id'), 10);
SELECT setval(pg_get_serial_sequence('stu_tracker.Locations',       'id'), 10);
SELECT setval(pg_get_serial_sequence('stu_tracker.Programs',        'id'), 10);
SELECT setval(pg_get_serial_sequence('stu_tracker.Subjects',        'id'), 10);
SELECT setval(pg_get_serial_sequence('stu_tracker.Tutors',          'id'), 10);
SELECT setval(pg_get_serial_sequence('stu_tracker.Students',        'id'), 10);
SELECT setval(pg_get_serial_sequence('stu_tracker.Sessions',        'id'), 10);
SELECT setval(pg_get_serial_sequence('stu_tracker.Admin_root',      'id'), 10);