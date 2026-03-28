# Part 1 - System Overview

## Architecture Style
The platform uses a separated Frontend + Backend API model and is designed to evolve into a microservices architecture.

## High-Level Components
1. Frontend Web App (SPA-style pages with Vanilla JS)
2. Backend API (Go Fiber)
3. MySQL database
4. Internal AIOps service (content analysis and moderation support)
5. Admin/Moderator operations console

## Deployment Model
- Deploy on Cloud VPS/AWS/GCP.
- Scale API instances horizontally for traffic bursts.
- Keep moderation and sensitive data services under controlled infrastructure.

## Why This Structure
- Easy future expansion to mobile app clients.
- Better security boundaries.
- Faster independent deployment of frontend and backend.
- Better control for high voting traffic periods.
