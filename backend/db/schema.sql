-- Canh Buom Xanh - Initial schema (Part B)
-- MySQL 8+

CREATE DATABASE IF NOT EXISTS canhbuomxanh CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE canhbuomxanh;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(190) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  full_name VARCHAR(120) NOT NULL,
  role ENUM('student_current', 'student_alumni', 'parent', 'teacher', 'moderator', 'admin') NOT NULL,
  status ENUM('active', 'suspended', 'deleted') NOT NULL DEFAULT 'active',
  avatar_url VARCHAR(255) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS teachers (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NULL,
  display_name VARCHAR(120) NOT NULL,
  school_name VARCHAR(190) NULL,
  subject_name VARCHAR(120) NULL,
  claim_status ENUM('unclaimed', 'pending', 'claimed', 'rejected') NOT NULL DEFAULT 'unclaimed',
  is_public TINYINT(1) NOT NULL DEFAULT 0,
  opt_out_requested_at DATETIME NULL,
  opt_out_executed_at DATETIME NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_teachers_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS posts (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  teacher_id BIGINT UNSIGNED NULL,
  title VARCHAR(190) NOT NULL,
  body TEXT NOT NULL,
  visibility ENUM('public', 'private') NOT NULL DEFAULT 'public',
  status ENUM('active', 'hidden', 'deleted') NOT NULL DEFAULT 'active',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_posts_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_posts_teacher_id FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS comments (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  post_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  parent_id BIGINT UNSIGNED NULL,
  body TEXT NOT NULL,
  status ENUM('active', 'hidden', 'deleted') NOT NULL DEFAULT 'active',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_comments_post_id FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  CONSTRAINT fk_comments_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_comments_parent_id FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS votes (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  teacher_id BIGINT UNSIGNED NOT NULL,
  voter_user_id BIGINT UNSIGNED NOT NULL,
  raw_score TINYINT UNSIGNED NOT NULL,
  weighted_score DECIMAL(6,2) NOT NULL,
  vote_mode ENUM('normal', 'ghost') NOT NULL DEFAULT 'normal',
  semester_key VARCHAR(20) NOT NULL,
  merge_status ENUM('pending_freeze', 'merged') NOT NULL DEFAULT 'pending_freeze',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT chk_votes_raw_score CHECK (raw_score BETWEEN 1 AND 5),
  CONSTRAINT fk_votes_teacher_id FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE,
  CONSTRAINT fk_votes_voter_user_id FOREIGN KEY (voter_user_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE KEY uq_votes_teacher_voter_semester (teacher_id, voter_user_id, semester_key)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS flags (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  source_type ENUM('post', 'comment', 'teacher_profile') NOT NULL,
  source_id BIGINT UNSIGNED NOT NULL,
  risk_level ENUM('yellow', 'red') NOT NULL,
  reason_code VARCHAR(80) NOT NULL,
  reason_detail VARCHAR(500) NULL,
  status ENUM('queued', 'reviewed', 'resolved') NOT NULL DEFAULT 'queued',
  assigned_moderator_id BIGINT UNSIGNED NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_flags_assigned_moderator_id FOREIGN KEY (assigned_moderator_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS appeals (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  teacher_id BIGINT UNSIGNED NOT NULL,
  requester_user_id BIGINT UNSIGNED NOT NULL,
  target_type ENUM('post', 'comment', 'profile') NOT NULL,
  target_id BIGINT UNSIGNED NOT NULL,
  reason TEXT NOT NULL,
  status ENUM('open', 'accepted', 'rejected') NOT NULL DEFAULT 'open',
  reviewed_by BIGINT UNSIGNED NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_appeals_teacher_id FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE,
  CONSTRAINT fk_appeals_requester_user_id FOREIGN KEY (requester_user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_appeals_reviewed_by FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS badges (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  code VARCHAR(40) NOT NULL UNIQUE,
  display_name VARCHAR(80) NOT NULL,
  min_score DECIMAL(8,2) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS teacher_badges (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  teacher_id BIGINT UNSIGNED NOT NULL,
  badge_id BIGINT UNSIGNED NOT NULL,
  awarded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uq_teacher_badges_teacher_badge (teacher_id, badge_id),
  CONSTRAINT fk_teacher_badges_teacher_id FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE,
  CONSTRAINT fk_teacher_badges_badge_id FOREIGN KEY (badge_id) REFERENCES badges(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE INDEX idx_posts_teacher_id ON posts (teacher_id);
CREATE INDEX idx_comments_post_id ON comments (post_id);
CREATE INDEX idx_votes_teacher_merge_status ON votes (teacher_id, merge_status);
CREATE INDEX idx_flags_status_risk ON flags (status, risk_level);
CREATE INDEX idx_teachers_claim_status ON teachers (claim_status);
