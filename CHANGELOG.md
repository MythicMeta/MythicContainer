# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [v1.0.9-rc12] - 2023-07-17

### Changed

- Fixed the tracking for c2 service binaries

## [v1.0.9-rc11] - 2023-07-12

### Changed

- Fixed the taskData.Args.GetArrayArg to properly cast to []string from []interface{}

## [v1.0.9-rc10] - 2023-06-26

### Changed

- Added the `WrappedPayloadUUID` value to a payload build message so you don't just get the raw bytes

## [v1.0.9-rc09] - 2023-06-09

### Changed

- Updated the grpc code to set maxInt for the send/recv limits with the translation containers

## [v1.0.9-rc08] - 2023-06-08

### Changed

- Added additional check if given a string and no parseArgString function defined, to just default to the raw command line

## [v1.0.9-rc07] - 2023-06-01

### Changed

- Updating queue name for logging/webhooks to be unique so we don't round robin the information

## [v1.0.9-rc06] - 2023-06-01

### Changed

- Added a fix to register new response logging data

## [v1.0.9-rc05] - 2023-05-31

### Changed

- Added new logging type for responses

## [v1.0.9-rc04] - 2023-05-23

### Changed

- Updated the SendMythicRPCFileUpdate function to support changing the DeleteAfterFetch attribute

## [v1.0.9-rc02] - 2023-05-23

### Changed

- Modified many of the similar C2 message structs to support new helper functions for getting arguments
- Modified the use of the supplied parameter group from the Mythic UI to be a tie breaker rather than as a manually set group name

## [v1.0.9-rc01] - 2023-05-22

### Changed

- Added base functionality for two new C2 RPC functions - GetIOC and SampleMessage
- Changed PayloadBuildMessage.BuildParameters to be a struct with a Parameters map inside of it
  - Added a suite of helper functions on it to get build parameters of various types
- Updated PTTaskMessageArgsData.Get*Arg functions to return default type-based blank values if nil
- Added suite of helper functions to PayloadBuildMessage.PayloadBuildC2Profile entries for getting C2 Parameter arguments
- Bumped the container version to v1.1.0 to account for new getIOC and SampleMessage C2 RPC Functionality

## [v1.0.8] - 2023-05-22

### Changed

- Updated tasking to make sure specified parameter groups in the UI carry over
- Updated tasking to list out unused parameters via the task's stdout/stderr modal

## [v1.0.7] - 2023-05-17

### Changed

- Updated the constant definitions for SupportedOS values to match the PyPi side with a capital first letter for all but macOS

## [v1.0.6] - 2023-05-10

### Changed

- Fixed the logging service capabilities to respect the log level defined (it was being overridden by Mythic's logging level)
- Fixed translation services gRPC connections that weren't reconnecting

## [v1.0.5] - 2023-05-09

### Changed
- Updated the way manual parameter group name is set during create tasking - now use `taskData.Args.SetManualParameterGroup`

## [v1.0.0-rc13] - 2023-04-25

### Changed

- Fixed a bug where new alert and new custom webhook fields weren't tracked for existence

## [v1.0.0-rc12] - 2023-04-21

### Changed

- Added the ability to return updated filename when building payloads
- Added a lot of docstrings for agent structures/building

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
