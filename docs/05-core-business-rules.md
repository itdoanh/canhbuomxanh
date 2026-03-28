# Part 1 - Core Business Rules

These rules must be implemented in backend code, not only in UI.

## Badge System
- Do not expose local top-1/top-2/top-3 competition.
- Use milestone badges (example: Silver Sail, Golden Sail).

## Claim Profile
- Teacher profile remains private until ownership verification is successful.
- Verification channels can include email/SMS.

## Opt-out and Data Deletion
- Teacher can opt out at any time.
- Data deletion target SLA: within 24 hours.

## Anti-Manipulation
- Ghost vote mode for user safety under coercion.
- Alarm/report channel for forced-vote incidents.

## Weighted Voting
- Alumni votes have higher weight than current students.
- Current-student votes are frozen and merged after semester ends.

## Moderation and Appeals
- AIOps flags high-risk content to moderator queues.
- Teachers can submit appeals for defamatory content.
