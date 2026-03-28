# Canh Buom Xanh - AI Working Context

## Mission
Build a student-safe teacher feedback platform with transparent honor badges, anti-manipulation voting, and strict privacy controls.

## Mandatory Product Constraints
1. Frontend and Backend are fully separated.
2. Architecture follows a basic microservices direction.
3. Deployment target is cloud infrastructure (Cloud VPS/AWS/GCP), not shared local hosting.
4. Internal AIOps must run on project-owned servers.
5. Do NOT use third-party GenAI APIs for moderation data (no OpenAI API, no Google Cloud AI API).
6. Data privacy for students and teachers is highest priority.

## Core Domains
- Guest/Auth
- User (student/alumni/parent)
- Teacher Portal
- Moderator (CTV)
- Admin
- Internal AIOps

## Non-Negotiable Backend Business Rules
- Badge system replaces local rank leaderboard.
- Claim Profile required before teacher profile becomes public.
- Teacher opt-out can be requested anytime; user data deletion SLA: within 24h.
- Ghost vote flow to protect students under voting pressure.
- Weighted voting: alumni weight > current student weight.
- Current-student votes are frozen until semester end before score merge.

## Tech Baseline
- Frontend: HTML5, CSS3, Vanilla JS, GSAP, Lenis, VanillaTilt, Odometer, Axios
- Backend API: Go + Fiber
- Data: MySQL

## Collaboration Protocol For AI Chats
1. Work by large parts only.
2. Complete one part end-to-end before moving on.
3. Wait for explicit user command "tiep" before starting the next part.
4. Keep outputs aligned with files in docs/ as source of truth.

## Current Delivery Status
- Part 1 (Architecture + Tech Stack + Project decomposition): COMPLETED
- Part 2 / Part B (Data + Backend Skeleton): COMPLETED
- Part 3 / Part C (Frontend baseline and UI system): COMPLETED
- Part 4 / Part D (Core feature integration): COMPLETED
- Part 5 / Part E (Voting and badge engine): COMPLETED
- Part 6 / Part F (Teacher, Moderator, Admin portals): COMPLETED
- Part 7 / Part G (Internal AIOps service): COMPLETED
- Part 8 / Part H (Hardening and release): COMPLETED
