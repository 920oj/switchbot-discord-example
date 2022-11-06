package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func getDeviceList(s *discordgo.Session, m *discordgo.MessageCreate) {
	// botユーザ自身のメッセージであれば無視する
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != "!devices" {
		return
	}

	b, err := requestGetDeviceList()
	if err != nil {
		fmt.Println("デバイス一覧取得に失敗しました。", err.Error())
	}

	message := ""
	for _, v := range b.Body.DeviceList {
		message += v.DeviceName + ": " + v.DeviceID + "\r"
	}

	if len(b.Body.DeviceList) == 0 {
		s.ChannelMessageSend(m.ChannelID, "デバイスが登録されていません。")
	}

	s.ChannelMessageSend(m.ChannelID, message)
}

func toggleBotLight(s *discordgo.Session, m *discordgo.MessageCreate) {
	// botユーザ自身のメッセージであれば無視する
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != "!kitchen" {
		return
	}

	botDeviceId := switchbotMac

	status, err := requestGetBotDeviceStatus(botDeviceId)
	if err != nil {
		fmt.Println("ステータス取得失敗: ", err.Error())
		return
	}

	if status.Body.Power == "off" {
		res, err := requestPostBotCommand(botDeviceId, "turnOn")
		if err != nil {
			fmt.Println("キッチン照明操作に失敗しました。: ", err.Error())
			return
		}
		if res.Message == "success" {
			s.ChannelMessageSend(m.ChannelID, "キッチン照明: ONにしました")
		} else {
			s.ChannelMessageSend(m.ChannelID, "キッチン照明: ONにできませんでした")
		}
	} else {
		res, err := requestPostBotCommand(botDeviceId, "turnOff")
		if err != nil {
			fmt.Println("キッチン照明操作に失敗しました。: ", err.Error())
			return
		}
		if res.Message == "success" {
			s.ChannelMessageSend(m.ChannelID, "キッチン照明: OFFにしました")
		} else {
			s.ChannelMessageSend(m.ChannelID, "キッチン照明: OFFにできませんでした")
		}
	}
}
