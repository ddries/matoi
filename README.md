
# matoi

matoi is a basic websocket client mainly created to test websocket servers under development in testing environments.

It is a basic echo websocket client: it will read from console (stdin) and send data as plain text to the server,
as well as printing the server messages to console (stdout)

## Installation

Download the source code and compile it using go >= 1.17. Then simply run `matoi`

## Usage

`matoi <url> [flags]`