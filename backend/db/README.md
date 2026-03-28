# Database Setup (Part B)

## Quick Start
1. Create a MySQL 8 database server.
2. Run schema file:
   - `mysql -u root -p < backend/db/schema.sql`
3. Confirm tables are created successfully.

## Notes
- The schema includes baseline tables for Auth, Teacher profile, Forum, Voting, Moderation, and Badge logic.
- Business rules such as claim status, opt-out timestamps, weighted vote columns, and vote freeze state are represented at data level.
