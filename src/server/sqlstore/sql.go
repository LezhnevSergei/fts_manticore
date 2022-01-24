package sqlstore

const ListRules = `
select r.rule_id,
       r.display_name,
       r.description
from rules as r
limit 50;
`

const IncidentsList = `
select incident_id, display_name, description
from (
	values
		%v
) as t (id, rule_id)
    left join incidents on incidents.rule_id = t.rule_id
    left join rules on rules.rule_id = t.rule_id
limit 100;
`
