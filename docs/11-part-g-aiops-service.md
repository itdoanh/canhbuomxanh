# Part 7 (Part G) - Internal AIOps Service Delivery

## Scope Completed
1. Vietnamese lexicon ingestion from internal file source.
2. Text risk scoring pipeline with severity-weighted matching.
3. Red/yellow/none classifier output.
4. Moderator handoff integration with enriched queue payloads.

## Backend AIOps Endpoints
- POST /api/v1/aiops/lexicon/reload
- GET /api/v1/aiops/lexicon/summary
- POST /api/v1/aiops/analyze
- POST /api/v1/aiops/analyze-and-flag
- GET /api/v1/aiops/queue/enriched

## Lexicon and Scoring Model
- Lexicon file: backend/internal/modules/aiops/data/lexicon_vi.txt
- Entry format: term,severity,category
- Score = sum(term_severity * hit_count)
- Risk thresholds:
  - red: score >= 8
  - yellow: score >= 3
  - none: score < 3

## Moderation Handoff
- analyze-and-flag inserts queued flags to flags table.
- enriched queue adds handoffHint for moderator action guidance.
- moderator UI now reads enriched queue endpoint.

## Frontend Deliverables
- aiops.html
- assets/js/aiops.js
- shared portal style reused (assets/css/pages/portal.css)

## Compliance Notes
- Service is internal-only and does not call any third-party GenAI API.
- Processing remains inside project infrastructure boundaries.

## What Is Intentionally Deferred
- Advanced NLP normalization for Vietnamese tone/variants.
- Model-based semantic moderation.
- Persistent lexicon editing UI with DB storage.

## Follow-up
Part H hardening and release delivery is completed in docs/12-part-h-hardening-release.md.
