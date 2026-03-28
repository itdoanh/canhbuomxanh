# Part 5 (Part E) - Voting and Badge Engine Delivery

## Scope Completed
1. Weighted vote logic by voter role implemented in backend.
2. Vote freeze rule for current students implemented via merge_status.
3. Semester release flow implemented to merge pending votes.
4. Badge milestone recompute engine implemented.
5. Ghost vote safety signaling implemented with forced-vote alert endpoint.

## Backend Engine Endpoints
- POST /api/v1/voting/votes
- POST /api/v1/voting/votes/release
- POST /api/v1/voting/alerts/forced-vote
- POST /api/v1/voting/badges/recompute
- GET /api/v1/voting/badges/teachers/:id

## Core Rule Mapping
- Alumni weight > current student weight:
  - student_alumni = 1.4
  - student_current = 1.0
  - parent = 1.1
- Current student votes are not merged immediately:
  - merge_status = pending_freeze
- Semester release moves pending votes to merged.
- Badge assignment is milestone-based from merged weighted score totals.
- Ghost vote mode is accepted and stored in vote_mode.
- Forced voting pressure can be reported to moderation queue (flags table).

## Frontend Deliverables
- voting.html
- assets/css/pages/voting.css
- assets/js/voting.js
- navigation link added across core pages

## What Is Intentionally Deferred
- Automatic scheduler for semester release.
- Advanced anti-abuse detection scoring.
- Moderator UI for processing forced-vote alerts.

## Next Part (Wait for "tiep")
Part F - Teacher, Moderator, Admin Portals:
- claim profile + opt-out controls
- appeal handling dashboards
- moderation and admin operation panels
