package main

import (
    "fmt"
    "github.com/sacOO7/gowebsocket"
    "github.com/marcusolsson/tui-go"
    "github.com/buger/jsonparser"
)

type chat struct {
    socket *gowebsocket.Socket
    ckey string
    state string
    chatBox *tui.Box
    statusLbl *tui.Label
    typingIndicator *tui.Label
}


func (chat *chat) setConnected() {
    status := "connected"

    chat.statusLbl.SetText(status)
    chat.statusLbl.SetStyleName(status)
}

func (chat *chat) sendPong() {
    msg := "4{\"ev_name\": \"_gdzie\"}"
    chat.socket.SendText(msg)
}

func (chat *chat) startTalk(data []byte) {
    cKey, _ := jsonparser.GetString(data, "ev_data", "ckey")
    chat.ckey = cKey

    chat.statusLbl.SetText("Talking")
    chat.statusLbl.SetStyleName("inTalk")
}

func (chat *chat) addStrangerMessage(data []byte) {
    msg, _ := jsonparser.GetString(data, "ev_data", "msg")
    chat.chatBox.Append(tui.NewLabel("Obcy: " + msg))
}

func (chat *chat) SendMessage(msg string) {
    toSend := fmt.Sprintf(`4{"ev_name":"_pmsg","ev_data":{"ckey":"%s","msg":"%s","idn":0},"ceid":3}`, chat.ckey, msg)
    chat.socket.SendText(toSend)
    chat.chatBox.Append(tui.NewLabel(fmt.Sprintf(`Ty: %s`, msg)))
}

func (chat *chat) InitializeTalk() {
    msg := `4{"ev_name":"_sas","ev_data":{"channel":"main","myself":{"sex":0,"loc":2},"preferences":{"sex":0,"loc":2}},"ceid":1}`
    chat.socket.SendText(msg)
    chat.typingIndicator.SetText("")
}

func (chat *chat) disconnectFromChat(sendRequest bool) {
    status := "connected"

    chat.statusLbl.SetText(status)
    chat.statusLbl.SetStyleName(status)

    chat.typingIndicator.SetText("Stranger disconnected")
}

func (chat *chat) setStrangetTyping(data []byte) {
    isTyping, _ := jsonparser.GetBoolean(data, "ev_data")

    typingText := ""
    if isTyping {
        typingText = "Stranger typing..."
    }

    chat.typingIndicator.SetText(typingText)
}


func (chat *chat) ReactToMsg(data []byte) {
    evName, err := jsonparser.GetString(data, "ev_name")
    if err != nil {
        return
    }
    switch evName {
        case "piwo":
            chat.sendPong()
        case "cn_acc":
            chat.setConnected()
        case "talk_s":
            chat.startTalk(data)
        case "rmsg":
            chat.addStrangerMessage(data)
        case "sdis":
            chat.disconnectFromChat(false)
        case "styp":
            chat.setStrangetTyping(data)
        case "count":
        default:
            chat.chatBox.Append(tui.NewLabel("Unknown message: " + string(data)))
    }
}
