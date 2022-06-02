# Demodata Service Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

Please note, that this project, while following numbering syntax, it DOES NOT
adhere to [Semantic Versioning](http://semver.org/spec/v2.0.0.html) rules.

## Types of changes

* ```Added``` for new features.
* ```Changed``` for changes in existing functionality.
* ```Deprecated``` for soon-to-be removed features.
* ```Removed``` for now removed features.
* ```Fixed``` for any bug fixes.
* ```Security``` in case of vulnerabilities.

## [2022.2.3.6] - 2022-06-02

### Changed

- faster logging

## [2022.2.2.18] - 2022-05-18

### Changed

- updated checkDatabaseConnection()

## [2022.2.1.28] - 2022-04-28

### Changed

- updated for creating docker images on apple silicon

## [2022.1.3.23] - 2022-03-23

### Changed
- one database connection for whole service


## [2022.1.3.21] - 2022-03-21

### Changed
- updated to latest Go and all libraries

## [2021.4.2.1] - 2021-11-01

### Changed
- updated to latest database with cache records

## [2021.2.3.28] - 2021-06-28

### Changed
- when updating order, adding user_id

## [2021.2.3.14] - 2021-06-15

### Changed
- inserting digital data every 20 seconds

## [2021.2.3.14] - 2021-06-14

### Added
- copyright
- updated libraries

## [2021.2.3.6] - 2021-06-06

### Bugfix

- proper updating terminal data

## [2021.2.3.5] - 2021-06-05

### Changed

- performance improvements: reduce database calls

## [2021.2.3.1] - 2021-06-01

### Changed

- updated max and min for generating production and downtime data

## [2021.2.2.27] - 2021-05-27

### Changed

- adding workplace and workplace workshifts
- changing default order, downtime and user records

## [2021.2.2.13] - 2021-05-13

### Changed

- updated to latest go 1.16.4
- updated to latest go libraries

## [2021.2.2.3] - 2021-05-03

### Changed
- updated to latest go
- updated to latest go libraries

## [2021.2.1.1] - 2021-04-01

### Fixed
- proper opening and closing database connection

## [2021.1.3.30] - 2021-03-30

### Changed
- updated to latest go
- updated to latest libraries
- changes according to state_id


## [2021.1.3.18] - 2021-03-18

### Changed
- updated to latest go
- updated to latest libraries

## [2021.1.2.22] - 2021-02-22

### Changed
- reduced database calling

## [2021.1.2.21] - 2021-02-21

### Changed
- updated to latest go
- updated to latest libraries

## [2020.4.3.14] - 2020-12-14

### Changed
- updated to latest go
- updated to latest libraries

## [2020.4.2.24] - 2020-11-24

### Fixed
- calling goroutines for updating data in for loop (instead of ordinary functions)

## [2020.4.2.17] - 2020-11-17

### Changed
- added creating additional orders, products, users, downtimes
- updating default downtime, order and user records

## [2020.4.1.26] - 2020-10-26

### Fixed
- fixed leaking goroutine bug when opening sql connections, the right way is this way

## [2020.4.1.1] - 2020-10-1

### Changed
- constant number for creating devices and workplaces

## [2020.3.3.28] - 2020-09-28

### Changed
- updated readme.md
- updated dockerfile
- updated create.sh script

### Added
- creating terminals and link them with workplaces
- creating proper workplace workshifts

## [2020.3.2.22] - 2020-08-29

### Changed
- functions naming changed to idiomatic go (ThisFunction -> thisFunction)

## [2020.3.2.22] - 2020-08-22

### Added
- automatic go get -u all when creating docker image


## [2020.3.2.4] - 2020-08-04

### Changed
- update to latest libraries and latest database changes
- removed all about config and logging to file

## [2020.1.2.29] - 2020-02-29

### Change
- update for latest database changes
- minor changes after testing new postgres, mariadb and mssql
- when searching for active devices, changed from "true" to "1"

## [2020.1.1.1] - 2020-01-01

### Added
- creates 20 devices and 20 workplace if not present
- generates pseudorandom analog and digital data for those 20 devices
