/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php

*/

package main

import (
	//	"io" //　ログ出力設定用。必要に応じて。
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"time"

	"github.com/Chouette2100/srhandler"
)

/*
	ルームの状態を調べroom_url_keyからroom_idがわかるApiRoomStatus()をテストする。

	ソースのダウンロード、ビルドについて以下簡単に説明します。詳細は以下の記事を参考にしてください。
	WindowsもLinuxも特記した部分以外は同じです。

		【Windows】かんたんなWebサーバーの作り方
			https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/c5cab5

		---------------------

		【Windows】Githubにあるサンプルプログラムの実行方法
			https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/e27fc9

		【Unix/Linux】Githubにあるサンプルプログラムの実行方法
			https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/220e38

			ロードモジュールさえできればいいということでしたらコマンド一つでできます。

【Unix/Linux】

	$ cd ~/go/src
	$ curl -OL https://github.com/Chouette2100/t010srapi/archive/refs/tags/v0.1.0.tar.gz
	$ tar xvf v0.1.0.tar.gz
	$ mv t010srapi-0.1.0 t010srapi
	$ cd t010srapi
	$ go mod init
	$ go mod tidy
	$ go build t010srapi.go
	$ ./t010srapi

	ここでブラウザを起動し

	http://localhost:8080/apiroomstatus?room_url_key=48_Hitomi_Otake

	のようにURLを入力すると、（上の例であれば大竹 ひとみさんの）ルーム情報が表示されます。

	【Windows】

	Microsoft Windows [Version 10.0.22000.856]
	(c) Microsoft Corporation. All rights reserved.

	C:\Users\chouette>cd go

	C:\Users\chouette\go>cd src

	作業はかならず %HOMEPATH%\go\src の下で行います。

	以下、要するに https://github.com/Chouette2100/t010srapi/releases にあるv0.1.0のZIPファイルSource code (zip) からソースをとりだしてくださいということなので、ブラウザでダウンロードしてエクスプローラで解凍というこでもけっこうです。なんならこの記事の最後にあるgithubのソースをエディターにコピペで作るということでもかまいません（この場合文字コードはかならずUTF-8にしてください 改行はLFになっています。というようなことを考えるとやっぱりダウンロードして解凍が安全かも）

	C:\Users\chouette\go\src>mkdir t010srapi

	C:\Users\chouette\go\src>cd t010srapi

	C:\Users\chouette\go\src\t010srapi>curl -OL https://github.com/Chouette2100/t010srapi/archive/refs/tags/v0.1.0.zip
	  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
	                                 Dload  Upload   Total   Spent    Left  Speed
	  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
	100  6265    0  6265    0     0   6777      0 --:--:-- --:--:-- --:--:-- 16400

	C:\Users\chouette\go\src\t010srapi>call powershell -command "Expand-Archive v0.1.0.zip"

	C:\Users\chouette\go\src\t010srapi>tree
	フォルダー パスの一覧
	ボリューム シリアル番号は E2CD-BDF1 です
	C:.
	└─v0.1.0
	    └─t010srapi-0.1.0
	        ├─public
	        └─templates

	C:\Users\chouette\go\src\t010srapi>xcopy /e v0.1.0\t010srapi-0.1.0\*.*
	v0.1.0\t010srapi-0.1.0\freebsd.bat
	v0.1.0\t010srapi-0.1.0\freebsd.sh
	v0.1.0\t010srapi-0.1.0\LICENSE
	v0.1.0\t010srapi-0.1.0\README.md
	v0.1.0\t010srapi-0.1.0\t010srapi.go
	v0.1.0\t010srapi-0.1.0\public\index.html
	v0.1.0\t010srapi-0.1.0\templates\top.gtpl
	以下省略（リリースによって内容が異なる場合があります。）

	C:\Users\chouette\go\src\t010srapi>rmdir /s /q v0.1.0
	以下省略（リリースによって内容が異なる場合があります。）

	C:\Users\chouette\go\src\t010srapi>del v0.1.0.zip

	ここで次のような構成になっていればOKです。top.gtpl と index.html が所定の場所にあることをかならず確かめてください。

	C:%HOMEPATH%\go\src\t010srapi --+-- t010srapi.go
	                                |
	                                +-- \templates --- t010top.gtpl
	                                |
	                                +-- \public    --- index.html

	ここからはコマンド三つでビルドが完了します。

	C:\Users\chouette\go\src\t010srapi>go mod init
	go: creating new go.mod: module t010srapi
	go: to add module requirements and sums:
	        go mod tidy

	C:\Users\chouette\go\src\t010srapi>go mod tidy
	go: finding module for package github.com/dustin/go-humanize
	go: downloading github.com/dustin/go-humanize v1.0.0
	go: found github.com/dustin/go-humanize in github.com/dustin/go-humanize v1.0.0

	C:\Users\chouette\go\src\t010srapi>go build t010srapi.go

	あとは

	C:\Users\chouette\go\src\t010srapi>t010srapi

	でWebサーバが起動します。ここでセキュリティー上の警告が出ると思いますが、説明をよく読んで問題ないと思ったらアクセスを許可してください（もちろん許可しなければWebサーバは使えなくなります）

	Webサーバを起動したままにしておいてブラウザを開き

	http://localhost:8080/apiroomstatus?room_url_key=48_Hitomi_Otake

	のようにURLを入力すると、（上の例であれば大竹 ひとみさんの）ルーム情報が表示されます。

	Ver. 0.1.0

*/



//Webサーバーを起動する。
func main() {

	logfilename := time.Now().Format("20060102") + ".txt"
	logfile, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open logfile: " + logfilename + err.Error())
	}
	defer logfile.Close()

	/*
		//	ログをコンソールにも出力する。原則これは行わないこと。とくにCGIの場合はぜったいダメ。
		//	"io"をimportすること。
		log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	*/

	//	ログ出力を設定する。
	log.SetOutput(logfile)

	rootPath := os.Getenv("SCRIPT_NAME")
	log.Printf("rootPath: \"%s\"\n", rootPath)

	//	URLに対するハンドラーを定義する。この例では /top の1行しかないが、実際はURLのある分だけ羅列する。
	http.HandleFunc(rootPath+"/t008top", srhandler.HandlerT008topForm) //	http://localhost:8080/t008top で呼び出される。
	http.HandleFunc(rootPath+"/t009top", srhandler.HandlerT009topForm) //	http://localhost:8080/t009top で呼び出される。
	http.HandleFunc(rootPath+"/apiroomstatus", srhandler.HandlerApiRoomStatus) //	http://localhost:8080/apiroomstatus?room_url_key=... で呼び出される。

	//	ポートは8080などを使います。
	//	Webサーバーはroot権限のない（特権昇格ができない）ユーザーで起動した方が安全だと思います。
	//	その場合80や443のポートはlistenできないので、ルータやOSの設定でポートフォワーディングするか
	//	ケーパビリティを設定してください。
	//	# setcap cap_net_bind_service=+ep ShowroomCGI
	//　（設置したあとこの操作を行うこと）
	httpport := "8080"

	//		CGIとして起動されたときはWebサーバーやCGIの設置場所にあわせて変更すること。
	//		さくらのレンタルサーバーでwwwの直下にCGIを設置したときはこのままでいいです。

	if rootPath == "" {
		//	Webサーバーとして起動された
		//		URLがホスト名だけのときは public/index.htmlが表示されるようにしておきます。
		http.Handle("/", http.FileServer(http.Dir("public"))) // http://localhost:8080/ で呼び出される。
		err = http.ListenAndServe(":"+httpport, nil)
	} else {
		//	cgiとして起動された
		err = cgi.Serve(nil)
	}
	if err != nil {
		log.Printf("%s\n", err.Error())
	}
}
