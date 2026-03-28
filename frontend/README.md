# Canh Buom Xanh Frontend (Part C)

## Scope
- Baseline UI system with Spatial + Glassmorphism direction.
- Home Hero foundation page.
- Honor page baseline with animated counters.

## Stack
- HTML5, CSS3, Vanilla JS
- GSAP
- Lenis
- VanillaTilt
- Odometer
- Axios

## Structure
- index.html: Hero landing baseline
- honor.html: Honor page baseline
- auth.html: register/login flow
- forum.html: forum posts/comments UI
- teacher.html: teacher public list/detail UI
- profile.html: personal profile UI
- voting.html: weighted vote and badge engine UI
- teacher-portal.html: claim, opt-out, appeal, dashboard UI
- moderator.html: AI queue, appeal review, violation action UI
- admin.html: system monitor, aiops config, access control UI
- aiops.html: internal lexicon, analysis, and handoff UI
- assets/css/tokens.css: design tokens
- assets/css/base.css: global foundations
- assets/css/components.css: reusable components
- assets/css/pages/*.css: page-specific styles
- assets/js/app.js: shared runtime setup
- assets/js/api.js: Axios helper + token injection
- assets/js/home.js: Hero motion
- assets/js/honor.js: Honor motion and counters
- assets/js/auth.js: auth interactions
- assets/js/forum.js: forum API interactions
- assets/js/teacher.js: teacher API interactions
- assets/js/profile.js: profile API interactions
- assets/js/voting.js: vote/badge engine API interactions
- assets/js/teacher-portal.js: teacher portal API interactions
- assets/js/moderator.js: moderator portal API interactions
- assets/js/admin.js: admin portal API interactions
- assets/js/aiops.js: internal aiops API interactions

## Run
- Open frontend/index.html using Live Server extension.
- Or use any static file host.

## Notes
- Frontend is fully separated from backend API.
- Core API connectivity for auth/forum/teacher/profile is available in Part D.
