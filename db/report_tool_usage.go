package db

import "context"

// ReportToolUsage describes recorded calls, not inferred test success.
// Returned includes errors; Pending means no matching result was recorded.
type ReportToolUsage struct {
	Tool string
	Total, Returned, Errors, Pending int
}

// ReportToolUsages counts each call once within a single exploration. Pair
// results by session and call id, and stop at the next reuse of that id.
// Never include raw input/output or another exploration's results.
func (d *DB) ReportToolUsages(ctx context.Context,expID int64)([]ReportToolUsage,error){
	rows,err:=d.QueryContext(ctx,`
WITH calls AS (
 SELECT id,tool,tool_use_id,worker,node_id,main_seg,
 LEAD(id) OVER (PARTITION BY tool_use_id,worker,node_id,main_seg ORDER BY id) AS next_id
 FROM activity WHERE exploration_id=$1 AND kind='tool_use'
)
SELECT COALESCE(NULLIF(u.tool,''),'-'),COUNT(*),
 COUNT(r.id),COUNT(*) FILTER (WHERE r.is_error IS TRUE),
 COUNT(*) FILTER (WHERE r.id IS NULL)
FROM calls u
LEFT JOIN LATERAL (
 SELECT r.id,r.is_error FROM activity r
 WHERE r.exploration_id=$1 AND r.kind='tool_result'
 AND NULLIF(u.tool_use_id,'') IS NOT NULL AND r.tool_use_id=u.tool_use_id
 AND r.worker IS NOT DISTINCT FROM u.worker
 AND r.node_id IS NOT DISTINCT FROM u.node_id
 AND r.main_seg IS NOT DISTINCT FROM u.main_seg
 AND r.id>u.id AND (u.next_id IS NULL OR r.id<u.next_id)
 ORDER BY r.id DESC LIMIT 1
) r ON true
GROUP BY 1 ORDER BY COUNT(*) DESC,1`,expID)
	if err!=nil{return nil,err}
	defer rows.Close()
	out:=[]ReportToolUsage{}
	for rows.Next(){var s ReportToolUsage;if err:=rows.Scan(&s.Tool,&s.Total,&s.Returned,&s.Errors,&s.Pending);err!=nil{return nil,err};out=append(out,s)}
	return out,rows.Err()
}
