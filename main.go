package main

func main() {
	chat := NewChatServer()
	err := chat.StartChatServer(8000)
	if err != nil {
		panic(err)
	}
	return
}
