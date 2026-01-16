package constants

const (
	GetOrganizationLimits = `SELECT p.plan_id, p.key, p.limit_value, p.enabled, p.enterprise 
	FROM stu_tracker.plan_entitlement p
	JOIN stu_tracker.organization_subscription os ON os.plan_id = p.plan_id
	WHERE os.organization_id = $1 AND p.enabled = TRUE;`
)
