# Changelog

## [2.0.0] - 2025-09-27

### Added
- `CmdRouter.GroupWithMiddleware()` method — allows creating submenus that inherit parent router middleware.

### Changed
- Replaced `Middlewares` with `Middleware` everywhere for clarity.
- Updated examples and tests to reflect new API changes.

### Fixed
- README updated: fixed minor errors and added table of contents.

---

## [1.0.0] - 2025-06-22

### Added
- Middleware signature updated to be more intuitive and easier to use.
- Menu path display added — the current path in nested menus is now shown to provide better context to users.
- Router configuration via functional options — added support for configuring the router with flexible functional options, including custom input/output streams.
- Standard middlewares included — built-in middlewares for panic recovery and error logging have been added to improve robustness.
- Continuous Integration and tests added — set up CI pipeline with automated tests to ensure code quality and reliability.
- Renaming for clarity: `OptionHandler` is now called `Option`, and `OptionHandler.Exec` renamed to `Option.Handler`.
- Introduced a new type `Handler func(ctx context.Context) error` for handler functions.
- All other logic regarding groups, adding options, and so forth remains unchanged.