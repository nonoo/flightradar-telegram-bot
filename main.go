package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"golang.org/x/exp/slices"
)

const errorStr = "❌ Error"

var telegramBot *bot.Bot
var cmdHandler cmdHandlerType

func sendReplyToMessage(ctx context.Context, replyToMsg *models.Message, s string) (msg *models.Message) {
	var err error
	msg, err = telegramBot.SendMessage(ctx, &bot.SendMessageParams{
		ReplyToMessageID: replyToMsg.ID,
		ChatID:           replyToMsg.Chat.ID,
		Text:             s,
	})
	if err != nil {
		fmt.Println("  reply send error:", err)
	}
	return
}

// func editReplyToMessage(ctx context.Context, msg *models.Message, s string) error {
// 	var err error
// 	_, err = telegramBot.EditMessageText(ctx, &bot.EditMessageTextParams{
// 		MessageID: msg.ID,
// 		ChatID:    msg.Chat.ID,
// 		Text:      s,
// 	})
// 	if err != nil {
// 		fmt.Println("  reply edit error:", err)
// 	}
// 	return err
// }

func sendTextToAdmins(ctx context.Context, s string) {
	for _, chatID := range params.AdminUserIDs {
		_, _ = telegramBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   s,
		})
	}
}

func handleMessage(ctx context.Context, update *models.Update) {
	fmt.Print("msg from ", update.Message.From.Username, "#", update.Message.From.ID, ": ", update.Message.Text, "\n")

	if update.Message.Chat.ID >= 0 { // From user?
		if !slices.Contains(params.AllowedUserIDs, update.Message.From.ID) {
			fmt.Println("  user not allowed, ignoring")
			return
		}
	} else { // From group ?
		fmt.Print("  msg from group #", update.Message.Chat.ID)
		if !slices.Contains(params.AllowedGroupIDs, update.Message.Chat.ID) {
			fmt.Println(", group not allowed, ignoring")
			return
		}
		fmt.Println()
	}

	// Check if message is a command.
	if update.Message.Text[0] == '/' || update.Message.Text[0] == '!' {
		cmd := strings.Split(update.Message.Text, " ")[0]
		if strings.Contains(cmd, "@") {
			cmd = strings.Split(cmd, "@")[0]
		}
		update.Message.Text = strings.TrimPrefix(update.Message.Text, cmd+" ")
		update.Message.Text = strings.TrimPrefix(update.Message.Text, cmd)
		cmdChar := string(cmd[0])
		cmd = cmd[1:] // Cutting the command character.
		switch cmd {
		case "frloc":
			fmt.Println("  interpreting as cmd frloc")
			cmdHandler.Location(ctx, update.Message)
		case "frrange":
			fmt.Println("  interpreting as cmd frrange")
			cmdHandler.Range(ctx, update.Message)
		case "frhelp":
			fmt.Println("  interpreting as cmd aaihelp")
			cmdHandler.Help(ctx, update.Message, cmdChar)
			return
		case "start":
			fmt.Println("  interpreting as cmd start")
			if update.Message.Chat.ID >= 0 { // From user?
				sendReplyToMessage(ctx, update.Message, "🤖 Welcome! This is the Flightradar Telegram Bot\n\n"+
					"More info: https://github.com/nonoo/flightradar-telegram-bot")
			}
			return
		default:
			fmt.Println("  invalid cmd")
			if update.Message.Chat.ID >= 0 {
				sendReplyToMessage(ctx, update.Message, errorStr+": invalid command")
			}
			return
		}
	}

	// if update.Message.Chat.ID >= 0 { // From user?
	// 	cmdHandler.TTS(ctx, update.Message.Text, update.Message)
	// }
}

func telegramBotUpdateHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.Text != "" {
		handleMessage(ctx, update)
	}
}

func main() {
	fmt.Println("flightradar-telegram-bot starting...")

	if err := params.Init(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	if err := settings.Load(); err != nil {
		fmt.Println("error: can't load settings:", err)
	}

	var cancel context.CancelFunc
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := airports.Load(ctx); err != nil {
		fmt.Println("error: can't load airports:", err)
		os.Exit(1)
	}
	if err := airlines.Load(ctx); err != nil {
		fmt.Println("error: can't load airlines:", err)
		os.Exit(1)
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(telegramBotUpdateHandler),
	}

	var err error
	telegramBot, err = bot.New(params.BotToken, opts...)
	if nil != err {
		panic(fmt.Sprint("can't init telegram bot: ", err))
	}

	flightData.Init(ctx)

	sendTextToAdmins(ctx, "🤖 Bot started")

	telegramBot.Start(ctx)
}
