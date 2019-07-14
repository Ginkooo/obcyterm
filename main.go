package main

import (
    "github.com/marcusolsson/tui-go"
    "github.com/sacOO7/gowebsocket"
)

func makeStatusLbl(theme *tui.Theme) *tui.Label {
    statusLbl := tui.NewLabel("Disconnected")

    statusLbl.SetSizePolicy(tui.Maximum, tui.Minimum)
    statusLbl.SetStyleName("")

    theme.SetStyle("label.disconnected", tui.Style{Fg: tui.ColorRed})
    theme.SetStyle("label.inTalk", tui.Style{Fg: tui.ColorRed})
    theme.SetStyle("label.connected", tui.Style{Fg: tui.ColorBlue})

    statusLbl.SetStyleName("disconnected")

    return statusLbl
}

func makeChatScroll(chatBox *tui.Box) *tui.ScrollArea {
    chatScroll := tui.NewScrollArea(chatBox)
    chatScroll.SetSizePolicy(tui.Expanding, tui.Expanding)

    return chatScroll
}

func makeChatBox() *tui.Box {
    return tui.NewVBox()
}

func makeInputBox(chat *chat) *tui.Box {
    input := tui.NewEntry()
    input.SetFocused(true)
    input.SetSizePolicy(tui.Minimum, tui.Minimum)

    input.OnSubmit(func(e *tui.Entry) {
        chat.SendMessage(e.Text())
        input.SetText("")
    })

    inputBox := tui.NewHBox(input)
    inputBox.SetBorder(true)
    inputBox.SetSizePolicy(tui.Minimum, tui.Minimum)

    return inputBox
}

func makeTypingIndicator() *tui.Label {
    typingIndicator := tui.NewLabel("")
    typingIndicator.SetSizePolicy(tui.Minimum, tui.Minimum)

    return typingIndicator
}

func main() {
    refreshChan := make(chan int)
    theme := tui.NewTheme()
    chatBox := makeChatBox()
    chatScroll := makeChatScroll(chatBox)
    statusLbl := makeStatusLbl(theme)
    typingIndicator := makeTypingIndicator()

    socket := gowebsocket.New("wss://server.6obcy.pl:7008/6eio/?EIO=3&transport=websocket")

    chatObj := chat{
        socket: &socket,
        ckey: "",
        state: "disconnected",
        chatBox: chatBox,
        statusLbl: statusLbl,
        typingIndicator: typingIndicator,
        refreshChannel: refreshChan,
    }

    inputBox := makeInputBox(&chatObj)

    socket.OnConnected = func(socket gowebsocket.Socket) {
    }

    socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
        panic(err)
    }

    socket.Connect()


    chatObj.InitializeTalk()

    spacer := tui.NewSpacer()
    spacer.SetSizePolicy(tui.Minimum, tui.Expanding)

    root := tui.NewVBox(
        tui.NewHBox(
            statusLbl,
            spacer,
            typingIndicator,
        ),
        chatScroll,
        inputBox,
    )

    ui, err := tui.New(root)

    socket.OnTextMessage = func(msg string, socket gowebsocket.Socket) {
        msg = msg[1:]
        chatObj.ReactToMsg([]byte(msg))
        ui.Repaint()
    }

    ui.SetTheme(theme)

    if err != nil {
        panic(err)
    }

    ui.SetKeybinding("Ctrl+n", func() {
        for chatBox.Length() > 0 {
            chatBox.Remove(0)
        }
        chatObj.InitializeTalk()
    })
    ui.SetKeybinding("Esc", func() { ui.Quit() })

    if err := ui.Run(); err != nil {
        panic(err)
    }
}
