package model

var (
    DbSuccess string
    SerialKeyLength int = 10
)

type Common struct {
    Title   string
    Content string
}

type Model struct {
    SerialKeyLength int
    SerialKeysCount int
    CommonData      Common
}

func New() *Model {
    m := new(Model)
    m.CommonData = Common{Title : "Some Title", Content: "next date"}

    return m
}







