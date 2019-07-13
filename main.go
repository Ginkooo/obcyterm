package main

import (
    "strings"
    "encoding/json"
    "github.com/marcusolsson/tui-go"
    "github.com/sacOO7/gowebsocket"
)


type chat struct {
    socket *gowebsocket.Socket
    ckey string
    state string
}

func sendPong(chat *chat) {
    msg := "4{\"ev_name\": \"_gdzie\"}"
    chat.socket.SendText(msg)
}

func reactToMsg(chatBox *tui.Box, data map[string]*json.RawMessage, chat *chat) {
    evNameRaw, err := json.Marshal(data["ev_name"])
    evName := strings.Trim(string(evNameRaw), "\"")
    if err != nil || evName == "null" {
        return
    }
    if evName == "piwo" {
        sendPong(chat)
    }
    label := tui.NewLabel(string(evName))
    label.SetSizePolicy(tui.Maximum, tui.Minimum)
    chatBox.Append(label)
}

func sendMessage(socket *gowebsocket.Socket, msg string, chat *chat) {
    panic("end")
}


func main() {

    socket := gowebsocket.New("wss://server.6obcy.pl:7008/6eio/?EIO=3&transport=websocket")

    chatObj := chat{socket: &socket, ckey: "", state: "disconnected"}

    socket.OnConnected = func(socket gowebsocket.Socket) {
    }

    socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
        panic(err)
    }

    socket.Connect()

    chatBox := tui.NewVBox()


    chatScroll := tui.NewScrollArea(chatBox)
    chatScroll.SetSizePolicy(tui.Expanding, tui.Expanding)

    socket.OnTextMessage = func(msg string, socket gowebsocket.Socket) {
        msg = msg[1:]
        data := make(map[string]*json.RawMessage)
        err := json.Unmarshal([]byte(msg), &data)
        if err != nil {
            return
        }
        reactToMsg(chatBox, data, &chatObj)
    }

    chatBox.SetSizePolicy(tui.Expanding, tui.Expanding)

    input := tui.NewEntry()
    input.SetFocused(true)
    input.SetSizePolicy(tui.Minimum, tui.Minimum)

    input.OnSubmit(func(e *tui.Entry) {
        sendMessage(&socket, e.Text(), &chatObj)
    })

    inputBox := tui.NewHBox(input)
    inputBox.SetBorder(true)
    inputBox.SetSizePolicy(tui.Minimum, tui.Minimum)

    root := tui.NewVBox(
        chatScroll,
        inputBox,
    )

    ui, err := tui.New(root)

    if err != nil {
        panic(err)
    }

    ui.SetKeybinding("Esc", func() { ui.Quit() })

    if err := ui.Run(); err != nil {
        panic(err)
    }
}
