

INSERT INTO stu_tracker.subscription_plan( code, name, stripe_price_id, is_active, cost_monthly, cost_yearly) VALUES ('trial', 'Trial', null, TRUE, 0, 0);
INSERT INTO stu_tracker.subscription_plan( code, name, stripe_price_id, is_active, cost_monthly, cost_yearly) VALUES ('pilot', 'Pilot plan', 'price_id_pilot', TRUE, 99, 950 );
INSERT INTO stu_tracker.subscription_plan( code, name, stripe_price_id, is_active, cost_monthly, cost_yearly) VALUES ('district', 'District plan', 'price_id_district', TRUE,399, 3800 );
INSERT INTO stu_tracker.subscription_plan( code, name, stripe_price_id, is_active, cost_monthly, cost_yearly) VALUES ('Enterprise', 'Enterprise plan', 'price_id_enterprise', TRUE, 0, 0);


INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (1, 'max_districts', 1,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (1, 'max_locations_per_district', 1,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (1, 'max_admin_per_district', 1,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (1, 'max_tutors_per_location', 2,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (1, 'max_students_per_location', 20,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (1, 'max_llm_tokens', 1000000,TRUE);

INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (2, 'max_districts', 1,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (2, 'max_locations_per_district', 5,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (2, 'max_admin_per_district', 10,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (2, 'max_tutors_per_location', 5,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (2, 'max_students_per_location', 100,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (2, 'max_llm_tokens', 2000000,TRUE);


INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (3, 'max_districts', 10,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (3, 'max_locations_per_district', 50,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (3, 'max_admin_per_district', 25,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (3, 'max_tutors_per_location', 10,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (3, 'max_students_per_location', 200,TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled) VALUES (3, 'max_llm_tokens', 3000000,TRUE);


INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled, enterprise) VALUES (4, 'max_districts', 10,TRUE, TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled, enterprise) VALUES (4, 'max_locations_per_district', 50,TRUE, TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled, enterprise) VALUES (4, 'max_admin_per_district', 50,TRUE, TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled, enterprise) VALUES (4, 'max_tutors_per_location', 20,TRUE, TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled, enterprise) VALUES (4, 'max_students_per_location', 100,TRUE, TRUE);
INSERT INTO stu_tracker.plan_entitlement (plan_id, key, limit_value, enabled, Enterprise) VALUES (4, 'max_llm_tokens', 5000000,TRUE,TRUE);
