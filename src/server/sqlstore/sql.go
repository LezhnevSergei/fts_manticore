package sqlstore

const IncidentsList = `
select i.incident_id, r.rule_id, r.display_name, r.description, l.link_id, l.display_name, h.host_id, h.display_name
from (
         values %v
         ) as t (id, incident_id)
         left join roswell.incidents as i on i.incident_id = t.incident_id
		 left join roswell.rules as r on i.rule_id = r.rule_id
		 left join roswell.hosts as h on i.target->>'host_id' = h.host_id
		 left join lucie.links as l on i.target->>'link_id' = l.link_id
limit 50;
`

const IncidentsListFull = `
select
       i.incident_id,
       concat(
           r.rule_id, ' ',
           r.display_name, ' ',
           r.description, ' ',
           l.link_id, ' ',
           l.display_name, ' ',
           h.host_id, ' ',
           h.display_name
       ) as fields
from roswell.incidents as i
left join roswell.rules as r on i.rule_id = r.rule_id
left join roswell.hosts as h on i.target->>'host_id' = h.host_id
left join lucie.links as l on i.target->>'link_id' = l.link_id
limit 100;
`
