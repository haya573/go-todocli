// :=は型推論しながら代入してくれる、Go言語はこれを使わないと型定義不足でエラーになる
// deferは関数の終了直前に実行する
package main

import (
    "encoding/json" // タスクをJSON形式で保存・読み込むために使用。
    "fmt" // 標準出力にテキストを表示するために使用。
    "os" // ファイル操作（保存・読み込み）やコマンドライン引数の取得に使用。
    "strconv" // コマンドライン引数（文字列）を数値に変換するために使用。
)

type Task struct { // TODOリストの1つのタスクを表す構造体。
    ID       int    `json:"id"` // タスクを識別する番号。
    Title    string `json:"title"` // タスクの内容（例: "買い物に行く"）。
    Completed bool   `json:"completed"` // タスクが完了しているかどうかを示すフラグ。
}

var tasks []Task // 全タスクを格納するスライス（リスト）。
const taskFile = "tasks.json" // タスクデータを保存するJSONファイルの名前。

func main() {
    if len(os.Args) < 2 {
        // コマンドライン引数が不足している場合にエラーメッセージを表示し終了します。
        // 例: go run main.go だけ実行した場合はエラー。
        fmt.Println("使い方: [add|list|done]")
        return
    }

    loadTasks()

    command := os.Args[1] // コマンドラインの最初の引数（例: add, list, done）を取得。
    switch command {
    case "add": // タスクの追加
        if len(os.Args) < 3 {
            fmt.Println("タスクの内容を指定してください")
            return
        }
        title := os.Args[2]
        addTask(title)
    case "list": // 現在のタスクリストを表示します。
        listTasks()
    case "done": // os.Args[2] に指定されたIDのタスクを完了にします。
        if len(os.Args) < 3 {
            fmt.Println("完了するタスクIDを指定してください")
            return
        }
        id, err := strconv.Atoi(os.Args[2]) // strconv.Atoiは整数に変換
        if err != nil {
            fmt.Println("タスクIDは数字で指定してください")
            return
        }
        markTaskAsDone(id)
    default:
        fmt.Println("不明なコマンドです: [add|list|done]")
    }
}

func saveTasks() {
    file, err := os.Create(taskFile) // ファイルを作成（既存ファイルがある場合は上書き）。
    if err != nil {
        fmt.Println("タスクを保存できませんでした:", err)
        return
    }
    defer file.Close() // 関数終了時にファイルを自動的に閉じる。

    encoder := json.NewEncoder(file)
    err = encoder.Encode(tasks) // タスクのスライスをJSON形式でファイルに保存。
    if err != nil {
        fmt.Println("タスクを保存できませんでした:", err)
    }
}

func loadTasks() { // プログラム起動時にタスクをファイルから読み込みます。
    file, err := os.Open(taskFile) // タスクデータファイルを開く。ファイルが存在しない場合は新規作成扱い（空のスライスをセット）。
    if err != nil {
        if os.IsNotExist(err) {
            tasks = []Task{}
            return
        }
        fmt.Println("タスクを読み込めませんでした:", err)
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&tasks) // ファイル内容をタスクスライスにデコード。
    if err != nil {
        fmt.Println("タスクを読み込めませんでした:", err)
        tasks = []Task{}
    }
}

func addTask(title string) {
    newTask := Task{
        ID:       len(tasks) + 1, // 新しいタスクのIDを設定。
        Title:    title,
        Completed: false,
    }
    tasks = append(tasks, newTask) // 新しいタスクをスライスに追加。
    saveTasks() // 変更後のタスクリストを保存。
    fmt.Printf("タスクを追加しました: %s\n", title)
}

func listTasks() {
    if len(tasks) == 0 { // タスクリストが空の場合にメッセージを表示。
        fmt.Println("タスクはありません")
        return
    }

    for _, task := range tasks { // 各タスクをループで表示。
        status := "未完了"
        if task.Completed { // 完了状態に応じてステータスを設定。
            status = "完了"
        }
        fmt.Printf("[%d] %s (%s)\n", task.ID, task.Title, status)
    }
}

func markTaskAsDone(id int) {
    for i, task := range tasks {
        if task.ID == id { // 指定されたIDのタスクを検索。
            tasks[i].Completed = true // 該当タスクの完了フラグをtrueに設定。
            saveTasks() // 変更内容をファイルに保存。
            fmt.Printf("タスクを完了にしました: %s\n", task.Title)
            return
        }
    }
    fmt.Println("指定されたIDのタスクが見つかりません")
}