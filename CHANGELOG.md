# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.0.0-rc11] - 2023-04-10

### Changed

- Fixed an issue with RabbitMQ Channels not closing resulting in an ID leak

## [v1.0.0-0.0.10] - 2023-03-20

### Added

- Added new structs for connection information command parameters to be more verbose

## [v1.0.0-0.0.9] - 2023-03-15

### Changed

- Updated create tasking functions to take pointer rather than value
- Started adding text descriptions for structs to make it easier for development

## [v1.0.0-0.0.8] - 2023-03-14

### Changed

- updated some structs to uint64 from int to match Mythic

## [v1.0.0-0.0.7] - 2023-03-12

### Changed

- fixed an issue with the Process response message routing to itself

## [v1.0.0-0.0.6] - 2023-03-03

### Changed

- fixed an issue where default int values weren't getting processed properly

## [v1.0.0-0.0.4] - 2023-03-01

### Changed

- updated the utils submodule to initialize on init() so that Mythic configuration can more easily be used in other projects

## [v1.0.0-0.0.3] - 2023-03-01

### Changed

- updated the logging submodule to initialize on init() for easier inclusion in other projects


## [v1.0.0-0.0.0] - 2023-02-28

### Added

- Created the initial push of this code
