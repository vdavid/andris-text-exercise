# AGENTS.md - Message Board App

## Project Overview

A simple Next.js "Message Board" web application with TypeScript, Supabase for the backend, and Vercel for deployment.

## Stack

- **Framework**: Next.js (App Router)
- **Language**: TypeScript (strict)
- **Database**: Supabase
- **Deployment**: Vercel
- **Check system**: Go-based check runner in `scripts/check/`

## File Structure

```
.
├── AGENTS.md              # This file
├── eslint.config.mjs      # ESLint flat config (strict TypeScript + Next.js)
├── .prettierrc             # Prettier config
├── scripts/
│   ├── check.sh           # Entry point: runs all checks
│   └── check/             # Go check runner
│       ├── main.go
│       ├── runner.go
│       └── checks/        # Individual check definitions
└── src/                   # Next.js application source (app/, components/, lib/)
```

## Running Checks

```bash
./scripts/check.sh              # Run all checks
./scripts/check.sh -app app     # Run only Next.js app checks
./scripts/check.sh -app scripts # Run only Go script checks
./scripts/check.sh -check app-typecheck  # Run a specific check
./scripts/check.sh -list        # List all available checks
./scripts/check.sh -fail-fast   # Stop on first failure
./scripts/check.sh -j 4         # Set parallelism
```

## Critical Rules

1. **Always run checks before considering work complete.** Run `./scripts/check.sh` and ensure all checks pass.
2. **Never skip or disable checks.** If a check fails, fix the underlying issue.
3. **Never ignore type errors.** Do not use `@ts-ignore`, `@ts-expect-error`, or `any` (unless there is genuinely no alternative, which is rare).
4. **Keep files under 300 lines.** The file-length check enforces this. Split large files into smaller, focused modules.
5. **Format code.** Prettier and gofmt handle this. Run them before committing.
6. **No `console.log` in production code.** The ESLint rule warns on this. Use proper error handling and logging.

## Apps

This project has two "apps" for the purpose of the check system:

- **app** - The Next.js web application (ESLint, Prettier, TypeScript checks)
- **scripts** - The Go check scripts themselves (gofmt, go vet, staticcheck)
