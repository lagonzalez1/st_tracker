-- Permissions for assessments
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:admin', 'Delete admins users', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:locations', 'Delete locations', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:assessments', 'Permission to delete assessments', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:assessments', 'Permission to view assessments', 'tutor', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:assessments', 'Permission to write assessments', 'root', 'admin');
-- Permissions for district
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:district', 'Permission to delete district', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:district', 'Permission to view district', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:district', 'Permission to write district', 'root', 'admin');

-- Permissions for location
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:location', 'Permission to delete location', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:location', 'Permission to view location', 'tutor', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:location', 'Permission to write location', 'root', 'admin');

-- Permissions for program
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:program', 'Permission to delete program', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:program', 'Permission to view program', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:program', 'Permission to write program', 'root', 'admin');

-- Permissions for program-location
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:program-location', 'Permission to delete program location', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:program-location', 'Permission to view program location', 'tutor', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:program-location', 'Permission to write program location', 'root', 'admin');

-- Permissions for semester
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:semester', 'Permission to delete semester', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:semester', 'Permission to view semester', 'view', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:semester', 'Permission to write semester', 'root', 'admin');

-- Permissions for students
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:students', 'Permission to delete students', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:students', 'Permission to view students', 'tutor', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:students', 'Permission to write students', 'root', 'admin');

-- Permissions for subject
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:subject', 'Permission to delete subject', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:subject', 'Permission to view subject', 'tutor', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:subject', 'Permission to write subject', 'root', 'admin');

-- Permissions for tutors
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:tutors', 'Permission to delete tutors', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:tutors', 'Permission to view tutors', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:tutors', 'Permission to write tutors', 'root', 'admin');

-- Permissions for tutor-locations
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:tutor-locations', 'Permission to delete tutor locations', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:tutor-locations', 'Permission to view tutor locations', 'tutor', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:tutor-locations', 'Permission to write tutor locations', 'root', 'admin');

-- Permissions for material
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:material', 'Delete material', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:material', 'Write material', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:material', 'View material', 'tutor', 'admin');

-- Admin permissions
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:admin', 'Write admin fields', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:locations', 'View location', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:admin', 'View administrators', 'root', 'admin');

-- Announcements permissions
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:announcements', 'delete announcements', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:announcements', 'Write announcements', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:announcements', 'view announcements', 'tutor', 'admin');

-- Session permissions
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:session', 'Write session', 'tutor', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:session', 'Update session', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:session', 'Update session', 'tutor', 'admin');

-- Permissions management
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:permissions', 'Write permissions', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:permissions', 'Delete permissions', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:permissions', 'View permissions', 'root', 'admin');

-- Tutor data permissions
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:tutor-data', 'write big tutor data', 'tutor', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:tutor-data', 'view big tutor data', 'tutor', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:tutor-data', 'delete big tutor data', 'tutor', 'admin');

-- Student data permissions
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:student-data', 'write big student data', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:student-data', 'view big student data', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:student-data', 'delete big student data', 'root', 'admin');


-- Semester locations
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:semester-location', 'Write semester locations', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:semester-location', 'Delete semester locations','root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:semester-location', 'View semester locations', 'root', 'admin');

-- Subject locations
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:subject-location', 'Write subject locations', 'root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:subject-location', 'Delete subject locations','root', 'admin');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:subject-location', 'View subject locations', 'root', 'admin');


-- Password
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:password', 'Change passwords from user accounts', 'root', 'root');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:password', 'Delete passwords from user accounts','root', 'root');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:password', 'View passwords from user accouts', 'root', 'root');


-- Emails
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:email', 'Change email from user accounts', 'root', 'root');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:email', 'Delete email from user accounts','root', 'root');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:email', 'View email from user accouts', 'root', 'root');

-- Root permissions
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'write:*', 'Root permission for write', 'root', 'root');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'delete:*', 'Root permission to delete','root', 'root');
INSERT INTO stu_tracker.Permissions( name, description, primary_role, secondary_role) VALUES ( 'view:*', 'Root permission to view', 'root', 'root');