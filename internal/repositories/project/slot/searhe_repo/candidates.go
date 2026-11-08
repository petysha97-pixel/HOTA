package cand

import (
	"HOTA/internal/models"
	"fmt"
	"strconv"
	"strings"
)



func FindCandidates(role string, stackIDs []int, ownerID int, query string, grade string) ([]models.Candidate, error) {
	if len(stackIDs) == 0 {
		return []models.Candidate{}, nil
	}
	if grade == "" {
		grade = "all"
	}

	idsText := make([]string, len(stackIDs))
	for i, id := range stackIDs {
		idsText[i] = strconv.Itoa(id)
	}
	stackList := strings.Join(idsText, ",")
	total := len(stackIDs)

	sqlText := fmt.Sprintf(`
		SELECT id, nickname, name, role, matched_stack, stack_count,
		       matched_stack * 100 / %[1]d AS stack_score,
		       role_score,
		       (matched_stack * 100 / %[1]d * 75 + role_score * 25) / 100 AS match_percent
		FROM (
			SELECT u.id, u.nickname, u.name, u.role,
			       COUNT(DISTINCT CASE WHEN us.stack_id IN (%[2]s) THEN us.stack_id END) AS matched_stack,
			       COUNT(DISTINCT us.stack_id) AS stack_count,
			       CASE WHEN u.role = ? THEN 100 ELSE 0 END AS role_score
			FROM users u
			JOIN user_stacks us ON us.user_id = u.id
			WHERE u.id <> ?
			  AND (? = 'all' OR LOWER(u.grade) = ?)
			  AND (? = '' OR instr(LOWER(u.nickname), ?) > 0 OR instr(LOWER(u.name), ?) > 0)
			GROUP BY u.id
		)
		WHERE matched_stack > 0
		ORDER BY match_percent DESC, stack_count ASC, id ASC
		LIMIT 100`, total, stackList)

	rows, err := models.UserDB.Query(sqlText, role, ownerID, grade, grade, query, query, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.Candidate{}
	for rows.Next() {
		var c models.Candidate
		if err := rows.Scan(&c.ID, &c.Nickname, &c.Name, &c.Role,
			&c.MatchedStack, &c.StackCount, &c.StackScore, &c.RoleScore, &c.Match); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}
