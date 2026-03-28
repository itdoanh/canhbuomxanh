# Part 6 (Part F) - Teacher, Moderator, Admin Portals Delivery

## Scope Completed
1. Teacher portal APIs and UI for claim profile, opt-out, appeal, and dashboard.
2. Moderator portal APIs and UI for AI queue, appeal queue, appeal review, and user violation actions.
3. Admin portal APIs and UI for system overview, spam-vote signals, AIOps config, and role management.

## Backend Portal Endpoints
### Teacher Portal
- POST /api/v1/teacher-portal/claim
- POST /api/v1/teacher-portal/opt-out
- POST /api/v1/teacher-portal/appeals
- GET /api/v1/teacher-portal/appeals
- GET /api/v1/teacher-portal/dashboard

### Moderator Portal
- GET /api/v1/moderator/queue
- GET /api/v1/moderator/appeals
- POST /api/v1/moderator/appeals/:id/review
- POST /api/v1/moderator/violations/user/:id

### Admin Portal
- GET /api/v1/admin/system/overview
- GET /api/v1/admin/spam-vote/teachers
- GET /api/v1/admin/aiops/config
- POST /api/v1/admin/aiops/config
- POST /api/v1/admin/access/role

## Frontend Portal Deliverables
- teacher-portal.html
- moderator.html
- admin.html
- assets/css/pages/portal.css
- assets/js/teacher-portal.js
- assets/js/moderator.js
- assets/js/admin.js

## Rule Mapping
- Claim profile turns claimed teacher profile to public.
- Opt-out request toggles profile hidden and records request timestamp.
- Teacher appeals are routed to moderator appeal queue.
- Moderator can accept/reject appeals and apply warn/suspend/ban actions.
- Admin can observe system volume, inspect vote concentration, tune AIOps config, and update user roles.

## What Is Intentionally Deferred
- Full 24h hard-delete worker for opt-out data deletion.
- Rich table UIs and filtering/sorting.
- Persistent AIOps config storage in dedicated table/service.

## Next Part (Wait for "tiep")
Part G - Internal AIOps Service:
- Vietnamese lexicon ingestion
- risk scoring pipeline
- red/yellow classifier service
- moderation handoff enrichments
