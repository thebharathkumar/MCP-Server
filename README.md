# pg-mcp

Production-grade MCP server for PostgreSQL with a safety-first policy model.

## Current development status

> Stop-the-line item: real `pg_query_go` parser integration is still pending in this environment.
>
> The classifier/policy architecture is in place, but parser-backed enforcement and parser-native AST fingerprints must land before handler implementation starts.

## v0.1 scope

- stdio transport only
- read-only `run_query`
- write workflow via `enable_writes` -> `preview_write` -> `commit_pending_write`
- top-level `INSERT`/`UPDATE`/`DELETE` supported for preview/commit
- writable CTEs blocked in all tools for v0.1

## Security-first setup

The primary security boundary is your PostgreSQL role privileges, not application logic.

Use a dedicated low-privilege role:

```sql
CREATE ROLE claude_readonly LOGIN PASSWORD 'replace_me';
GRANT CONNECT ON DATABASE appdb TO claude_readonly;
GRANT USAGE ON SCHEMA public TO claude_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO claude_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT ON TABLES TO claude_readonly;
```

If you provide superuser credentials, code-level read-only defaults are not sufficient.

## Claude Desktop note

Claude Desktop stores MCP config in plaintext on disk. Treat connection strings as sensitive and prefer dedicated least-privilege roles.

## SQL policy (v0.1)

- Single statement only (statement chaining rejected)
- `run_query`: SELECT-only, EXPLAIN denied (use explain tools)
- `explain_query`: supported for SELECT/INSERT/UPDATE/DELETE
- `explain_analyze_query`: SELECT-only in v0.1 (writes blocked)
- `preview_write`: top-level INSERT/UPDATE/DELETE only
- `TRUNCATE`, DDL, admin statements, LISTEN/UNLISTEN, LOCK, role/session auth changes are denied
- dangerous function denylist for read paths includes file/large object and dblink families

## Write semantics

- `enable_writes` toggles session write mode.
- In v0.1 (stdio), `WritesEnabled` persists for the process/session lifetime.
- `preview_write` returns:
  - affected row count
  - sampled rows
  - warnings (including UPDATE/DELETE without WHERE)
  - `ast_fingerprint`
  - `params_hash`
- `commit_pending_write` requires matching fingerprint metadata and unexpired pending-write token.

## Fingerprints

Two values are used:

- `ast_fingerprint`: normalized statement-shape identifier (**currently placeholder until parser integration lands**)
- `params_hash`: hash of canonical parameter payload

Both must match between preview and commit.

## Known limitations

- Preview is best-effort: sequences/triggers/FDW/external effects can escape rollback semantics.
- Writable CTEs are intentionally unsupported in v0.1.
- Function denylist is conservative and may expand in v0.2.
- Parser integration is pending in this environment due module-fetch restrictions; safety matrix re-run is required before handlers.

## Existing servers

For a baseline Postgres MCP server, see the official implementation:
- https://github.com/modelcontextprotocol/servers/tree/main/src/postgres

`pg-mcp` focuses on stricter safety workflow and explicit write controls for production use.
