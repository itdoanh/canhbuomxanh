package voting

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CastVote(ctx context.Context, userID uint64, role string, req castVoteRequest) (fiber.Map, error) {
	if req.TeacherID == 0 {
		return nil, errors.New("teacherId is required")
	}
	if req.RawScore < 1 || req.RawScore > 5 {
		return nil, errors.New("rawScore must be between 1 and 5")
	}
	if req.SemesterKey == "" {
		return nil, errors.New("semesterKey is required")
	}
	if req.VoteMode == "" {
		req.VoteMode = "normal"
	}
	if req.VoteMode != "normal" && req.VoteMode != "ghost" {
		return nil, errors.New("voteMode must be normal or ghost")
	}

	weight := voteWeightByRole(role)
	weightedScore := float64(req.RawScore) * weight

	mergeStatus := "merged"
	if role == "student_current" {
		mergeStatus = "pending_freeze"
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO votes (teacher_id, voter_user_id, raw_score, weighted_score, vote_mode, semester_key, merge_status)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE
		 raw_score = VALUES(raw_score),
		 weighted_score = VALUES(weighted_score),
		 vote_mode = VALUES(vote_mode),
		 merge_status = VALUES(merge_status)`,
		req.TeacherID,
		userID,
		req.RawScore,
		weightedScore,
		req.VoteMode,
		req.SemesterKey,
		mergeStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to save vote: %w", err)
	}

	return fiber.Map{
		"teacherId":      req.TeacherID,
		"rawScore":       req.RawScore,
		"weightedScore":  weightedScore,
		"voteMode":       req.VoteMode,
		"semesterKey":    req.SemesterKey,
		"mergeStatus":    mergeStatus,
		"voterRoleWeight": weight,
	}, nil
}

func (s *Service) ReleaseSemesterVotes(ctx context.Context, semesterKey string) (int64, error) {
	if semesterKey == "" {
		return 0, errors.New("semesterKey is required")
	}

	result, err := s.db.ExecContext(ctx,
		"UPDATE votes SET merge_status = 'merged' WHERE semester_key = ? AND merge_status = 'pending_freeze'",
		semesterKey,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to release votes: %w", err)
	}

	affected, _ := result.RowsAffected()
	return affected, nil
}

func (s *Service) ReportForcedVote(ctx context.Context, reporterID uint64, req forcedVoteAlertRequest) error {
	if req.TeacherID == 0 {
		return errors.New("teacherId is required")
	}

	detail := req.Detail
	if detail == "" {
		detail = "forced vote pressure reported"
	}

	reasonDetail := fmt.Sprintf("reporter_user_id=%d; detail=%s", reporterID, detail)
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO flags (source_type, source_id, risk_level, reason_code, reason_detail, status)
		 VALUES ('teacher_profile', ?, 'red', 'forced_vote_alarm', ?, 'queued')`,
		req.TeacherID,
		reasonDetail,
	)
	if err != nil {
		return fmt.Errorf("failed to create forced vote alert: %w", err)
	}

	return nil
}

func (s *Service) RecomputeBadges(ctx context.Context, teacherID uint64) (fiber.Map, error) {
	if err := s.ensureDefaultBadges(ctx); err != nil {
		return nil, err
	}

	teacherIDs, err := s.targetTeacherIDs(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	touched := 0
	for _, tID := range teacherIDs {
		if err := s.recomputeTeacherBadge(ctx, tID); err != nil {
			return nil, err
		}
		touched++
	}

	return fiber.Map{"updatedTeachers": touched}, nil
}

func (s *Service) GetTeacherBadgeSummary(ctx context.Context, teacherID uint64) (fiber.Map, error) {
	if teacherID == 0 {
		return nil, errors.New("teacher id is required")
	}

	var totalScore float64
	if err := s.db.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(weighted_score), 0) FROM votes WHERE teacher_id = ? AND merge_status = 'merged'",
		teacherID,
	).Scan(&totalScore); err != nil {
		return nil, fmt.Errorf("failed to load merged score: %w", err)
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT b.code, b.display_name, tb.awarded_at
		 FROM teacher_badges tb
		 JOIN badges b ON b.id = tb.badge_id
		 WHERE tb.teacher_id = ?
		 ORDER BY b.min_score ASC`,
		teacherID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query badge summary: %w", err)
	}
	defer rows.Close()

	badges := make([]fiber.Map, 0)
	for rows.Next() {
		var code, displayName string
		var awardedAt time.Time
		if err := rows.Scan(&code, &displayName, &awardedAt); err != nil {
			return nil, fmt.Errorf("failed to scan badge row: %w", err)
		}
		badges = append(badges, fiber.Map{
			"code":      code,
			"displayName": displayName,
			"awardedAt": awardedAt.Format(time.RFC3339),
		})
	}

	return fiber.Map{
		"teacherId":    teacherID,
		"mergedScore":  totalScore,
		"badgeCount":   len(badges),
		"badgeMilestones": badges,
	}, nil
}

func (s *Service) targetTeacherIDs(ctx context.Context, teacherID uint64) ([]uint64, error) {
	if teacherID != 0 {
		return []uint64{teacherID}, nil
	}

	rows, err := s.db.QueryContext(ctx, "SELECT id FROM teachers")
	if err != nil {
		return nil, fmt.Errorf("failed to list teachers: %w", err)
	}
	defer rows.Close()

	ids := make([]uint64, 0)
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan teacher id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *Service) recomputeTeacherBadge(ctx context.Context, teacherID uint64) error {
	var totalScore float64
	if err := s.db.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(weighted_score), 0) FROM votes WHERE teacher_id = ? AND merge_status = 'merged'",
		teacherID,
	).Scan(&totalScore); err != nil {
		return fmt.Errorf("failed to aggregate teacher score: %w", err)
	}

	_, err := s.db.ExecContext(ctx, "DELETE FROM teacher_badges WHERE teacher_id = ?", teacherID)
	if err != nil {
		return fmt.Errorf("failed to clear teacher badges: %w", err)
	}

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO teacher_badges (teacher_id, badge_id)
		 SELECT ?, b.id FROM badges b WHERE b.min_score <= ?`,
		teacherID,
		totalScore,
	)
	if err != nil {
		return fmt.Errorf("failed to assign teacher badges: %w", err)
	}

	return nil
}

func (s *Service) ensureDefaultBadges(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO badges (code, display_name, min_score)
		 VALUES
		 ('CANH_BUOM_BAC', 'Canh Buom Bac', 120),
		 ('CANH_BUOM_VANG', 'Canh Buom Vang', 220),
		 ('NHA_GIAO_KIEN_TAO', 'Nha Giao Kien Tao', 360)
		 ON DUPLICATE KEY UPDATE display_name = VALUES(display_name), min_score = VALUES(min_score)`,
	)
	if err != nil {
		return fmt.Errorf("failed to seed badges: %w", err)
	}
	return nil
}

func voteWeightByRole(role string) float64 {
	switch role {
	case "student_alumni":
		return 1.4
	case "parent":
		return 1.1
	case "student_current":
		return 1.0
	default:
		return 1.0
	}
}
