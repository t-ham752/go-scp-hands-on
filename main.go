package main

// 準備
// - 秘密鍵ファイルを用意する
// 	 - mkdir .ssh && ssh-keygen -t ed25519 -f .ssh/go-scp
// - SSHサーバーを用意する
// 	 - docker-compose build --no-cache && docker-compose up -d
// - DockerコンテナにSSHでログインする
// 	 - ssh -i .ssh/go-scp casone@localhost -p 20021
// - Dockerコンテナにファイルを転送する(scpが使えることを確認する)
// 	 - scp -i .ssh/go-scp -P 20021 main.go casone@localhost:/home/casone

func main() {
	// 秘密鍵ファイルの読み込み

	// 秘密鍵ファイルを解析し、署名情報を取得

	// SSHクライアントの設定を行う

	// SSHサーバーに接続

	// セッションを開始

	// SSHサーバーにコマンドを送信
}
