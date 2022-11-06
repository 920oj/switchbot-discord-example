package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	discordToken    string = ""
	switchbotToken  string = ""
	switchbotSecret string = ""
	switchbotMac    string = ""
)

func main() {
	// Tokenの読み込み
	err := loadEnv()
	if err != nil {
		fmt.Println("環境変数を読み込めませんでした。: ", err.Error())
		os.Exit(1)
		return
	}

	// Discord Botを初期化
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("Discord Botを初期化できませんでした。: ", err.Error())
		os.Exit(1)
		return
	}

	// Handlerを追加
	discord.AddHandler(getDeviceList)
	discord.AddHandler(toggleBotLight)

	// Websocketコネクションを開始
	err = discord.Open()
	if err != nil {
		fmt.Println("Websocketコネクション開始に失敗しました。: ,", err)
		return
	}

	fmt.Println("Botを開始しました。 CTRL-C で終了します。")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// キルシグナルを受信した場合Websocketコネクションをクローズする
	discord.Close()
}

func loadEnv() (err error) {
	err = godotenv.Load(".env")
	if err != nil {
		return err
	}

	discordToken = os.Getenv("DISCORD_TOKEN")
	switchbotToken = os.Getenv("SWITCHBOT_TOKEN")
	switchbotSecret = os.Getenv("SWITCHBOT_SECRET")
	switchbotMac = os.Getenv("SWITCHBOT_MAC")

	if discordToken == "" {
		return errors.New("discord Token が指定されていません。")
	}

	if switchbotToken == "" {
		return errors.New("switchbot Token が指定されていません。")
	}

	if switchbotSecret == "" {
		return errors.New("switchbot Secret が指定されていません。")
	}

	if switchbotMac == "" {
		return errors.New("switchbot Mac が指定されていません。")
	}

	return nil
}
