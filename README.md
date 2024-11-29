# GoZilla

## A Learning Project to Understand BitTorrent Trackers

GoZilla is an educational project designed to explore the technology behind a BitTorrent tracker. This implementation serves as a simple and fast BitTorrent tracker, built using Go and leveraging the Cobra library for command-line interface (CLI) functionality.

### Background

The purpose of this project is to learn about the inner workings of a BitTorrent tracker. A tracker plays a crucial role in facilitating peer-to-peer file sharing within the BitTorrent network by maintaining a list of connected peers and managing upload/download data. This implementation aims to demonstrate how these core concepts are implemented.

### Features

- Simple and fast BitTorrent tracker
- Command-line interface (CLI) using Cobra library
- Basic peer management: storing, updating, and deleting peer information
- Tracker commands for checking the version and server status
  
### Dependencies

This project relies on the following dependencies:

- github.com/spf13/cobra for CLI functionality
- github.com/gin-gonic/gin for HTTP server functionality
- gorm.io/gorm and gorm.io/driver/sqlite for database management
- github.com/jackpal/bencode-go for encoding/decoding BitTorrent metadata files

### Running the Project

To run this project, navigate to the project directory and execute:

```bash
go build && ./gozilla version
```

This will compile the project and display the GoZilla version.

**Note**: This README.md file provides a high-level overview of the project's purpose, features, dependencies, and running instructions. If you'd like me to add anything specific or make any changes, please let me know!
