# Changelog

## UNRELEASED
## 1.2.0
### Added
- Serve with sane default timeouts & graceful shutdown
## 1.1.0
### Added
- TTL cache

## 1.0.0
### Fixed
* Pass args to all middlewares in DefaultMiddleware

## 0.14.3

### Fixed

* LoganMiddleware entry stacking

## 0.14.2

### Fixed

* Add stack to recover log

## 0.14.1

### Fixed

* Proper type assert on log getter

## 0.14.0

### Added

* Typo fixing alias for `CtxMiddleware`
* `DefaultMiddlewares` helper to init router with safe defaults
* `Log` entry getter from request context

## 0.13.1

### Fixed

* proper dep chi constraint

## 0.13.0

### Changed

* problems.BadRequest determines status code based on interfaces, not types
