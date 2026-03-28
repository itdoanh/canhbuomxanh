# Part 1 - Tech Stack

## Developer Environment
- VS Code extensions:
  - Go (Google)
  - Live Server
  - Prettier
  - Thunder Client or Postman
  - PHP Intelephense (optional for database tooling ecosystem)

## Frontend
- Core: HTML5, CSS3, Vanilla JavaScript
- UI direction: Spatial Design + Glassmorphism
- Responsive layout: CSS Grid/Flex utility structure
- Motion and interaction:
  - GSAP
  - Lenis (smooth scroll)
  - VanillaTilt (3D card tilt)
  - Odometer (animated numeric counters)
- API integration: Axios

## Backend
- Language: Go
- Framework: Fiber
- Auth: JWT-based session tokens
- Modules: Routes, Controllers, Services, Models, Middleware

## Data Layer
- MySQL
- Suggested local admin tool in early phase: phpMyAdmin

## AI/Moderation Layer
- Internal AIOps only.
- No external GenAI API dependency for moderation decisions.
- Dictionary/rule/score based Vietnamese content checks as initial baseline.
