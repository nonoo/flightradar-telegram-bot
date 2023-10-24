# flightradar-telegram-bot

This bot announces flights which are passing over a given range and location.

The bot acquires data from [FlightRadar24.com](https://flightradar24.com/) and
uses the [Telegram Bot API](https://github.com/go-telegram-bot-api/telegram-bot-api).
Configuration is saved in a JSON file in the operating system's default app
configuration directory.

Settings are stored per chat ID, so you can use the bot in multiple groups in
case you need different settings (for example: different location settings).

Tested on Linux, but should be able to run on other operating systems.

## Compiling

You'll need Go installed on your computer. Install a recent package of `golang`.
Then:

```
go get github.com/nonoo/flightradar-telegram-bot
go install github.com/nonoo/flightradar-telegram-bot
```

This will typically install `flightradar-telegram-bot` into `$HOME/go/bin`.

Or just enter `go build` in the cloned Git source repo directory.

## Prerequisites

Create a Telegram bot using [BotFather](https://t.me/BotFather) and get the
bot's `token`.

## Running

You can get the available command line arguments with `-h`.
Mandatory arguments are:

- `-bot-token`: set this to your Telegram bot's `token`

Set your Telegram user ID as an admin with the `-admin-user-ids` argument.
Admins will get a message when the bot starts.

Other user/group IDs can be set with the `-allowed-user-ids` and
`-allowed-group-ids` arguments. IDs should be separated by commas.

You can get Telegram user IDs by writing a message to the bot and checking
the app's log, as it logs all incoming messages.

All command line arguments can be set through OS environment variables.
Note that using a command line argument overwrites a setting by the environment
variable. Available OS environment variables are:

- `BOT_TOKEN`
- `ALLOWED_USERIDS`
- `ADMIN_USERIDS`
- `ALLOWED_GROUPIDS`

## Supported commands

- `frloc (location)` - set or show current location
- `frrange (range)` - set or show current range

Only admins can set the location and range.

## Contributors

- Norbert Varga [nonoo@nonoo.hu](mailto:nonoo@nonoo.hu)

## Donations

If you find this bot useful then [buy me a beer](https://paypal.me/ha2non). :)
