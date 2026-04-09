# AC1 Add User Authentication (EXAMPLE)

> This is an example AC document, not an active change request.
> It demonstrates the AC pattern used in this repo. Use `ac-template.md` as the starting point for real ACs.

## Objective Fit

1. Users need to log in before accessing protected resources.
2. This is the highest-priority feature gap blocking the beta launch.
3. The app already has a session store; this adds the authentication layer on top.
4. Direct roadmap work.

## Summary

Add email/password authentication with session-based login. Users can register, log in, and log out. Protected routes redirect unauthenticated users to the login page.

## In Scope

- Registration endpoint with email and password
- Login endpoint that creates a session
- Logout endpoint that destroys the session
- Middleware that redirects unauthenticated requests to login
- Password hashing with bcrypt
- Tests for all endpoints and the middleware

## Out Of Scope

- OAuth or social login (future work)
- Password reset flow (separate AC)
- Role-based access control (not needed for beta)

## Implementation Notes

- Use the existing session store rather than adding a new dependency
- Store hashed passwords only; never log or return plaintext passwords
- Rate-limit login attempts to prevent brute-force attacks

## Acceptance Tests

- [Automated] Registration creates a user and returns success
- [Automated] Login with valid credentials creates a session
- [Automated] Login with invalid credentials returns 401
- [Automated] Protected route redirects to login when unauthenticated
- [Automated] Logout destroys the session
- [Manual] Verify login flow works end-to-end in the browser

## Documentation Updates

- Update API docs with new authentication endpoints
- Add a "Getting Started" section covering registration and login

## Status

EXAMPLE — not an active AC
