## [0.1.11] - 2018-08-10
### Added
- Hourly money income
- Fleets model and migration scripts
- Fleet creation webservice
- Player fleets and planet fleets retrieval webservices

## [0.1.10] - 2018-07-17
### Added
- Player wallet
- Ship model price stored at its creation
- Player ship models retrieval
- Planet hangar ships retrieval
- Ship creation
- Ship construction with military points
- Storage resources spending

### Fixed
- Price diff between the front and the API

## [0.1.9] - 2018-07-13
### Added
- Factions banner
- Configuration for test servers

### Fixed
- Error messages when faction select query fails

## [0.1.8] - 2018-07-05
### Added
- Exception catchers
- Planet settings model and update service
- Building construction with building points
- Ship, ship models and fleet model structures with migration scripts
- Player ship models retrieval webservice
- Ship model creation webservice

## [0.1.7] - 2018-03-30
### Added
- Diplomatic relations between factions
- Planet demography and population points model
- Population points update webservice
- Faction retrieval webservice

### Changed
- Buildings can now be built only once per planet
- Buildings now require building points instead of resources

## [0.1.6] - 2018-03-17
### Added
- Hourly takss in scheduling
- Player retrieval webservice
- Faction retrieval webservice
- Hourly calculation of planet resources production
- Planet storages
- File logger

## [0.1.5] - 2018-02-16
### Added
- Scheduling component
- Planet buildings

## [0.1.4] - 2018-01-08
### Added
- Hybrid encryption for portal communications
- Diplomatic relations initialization at game creation
- Initial player relations with starter planet
- Faction colors

## [0.1.3] - 2018-01-01
### Added
- Optimize server map generation with goroutines
- Planet resources
- Player first connection course
- Factions model
- Player binding to a planet
- Player binding to a faction
- Faction retrieval route
- Player planets retrieval

## [0.1.2] - 2017-12-25
### Added
- System retrieval route
- Planet retrieval route
- Planet types choice at map generation

### Fixed
- Registration of same player on multiple servers

## [0.1.1] - 2017-12-23
### Added
- Map model structures
- Map model migration scripts
- Map generation after server creation
- Map systems retrieval route

## [0.1.0] - 2017-12-08
### Added
- RSA keypair generation
- RSA encryption manager
- Crypted communication between portal and server
- Server DB model
- Player DB model
- Server registration
- Player registration
- JSON WebToken generation
- PostgreSQL ORM library
- Database migration external tool
- Mux routing utils
