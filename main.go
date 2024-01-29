package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

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
	privateKey, err := os.ReadFile(".ssh/go-scp")
	if err != nil {
		log.Fatalf("failed to read private key: %v", err)
	}

	// 秘密鍵ファイルを解析し、署名情報を取得
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", err)
	}

	// SSHクライアントの設定を行う
	config := &ssh.ClientConfig{
		User: "casone",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// SSHサーバーに接続
	client, err := ssh.Dial("tcp", "localhost:20021", config)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	// セッションを開始
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}

	// SSHサーバーにコマンドを送信
	rw, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("failed to create stdin pipe: %v", err)
	}
	defer rw.Close()

	rr, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("failed to create stdout pipe: %v", err)
	}

	f, err := os.Open("Dockerfile")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	err = session.Start("scp -t /home/casone")
	if err != nil {
		log.Fatalf("failed to start scp: %v", err)
	}

	fs, err := f.Stat()
	if err != nil {
		log.Fatalf("failed to stat: %v", err)
	}

	_, err = fmt.Fprintf(rw, "C0644 %d %s\n", fs.Size(), fs.Name())
	if err != nil {
		log.Fatalf("failed to send header: %v", err)
	}

	// これが無いと、ファイルの中身が送信されない
	buffer := make([]uint8, 1)
	_, err = rr.Read(buffer)
	if err != nil {
		log.Fatalf("failed to read: %v", err)
	}

	_, err = io.Copy(rw, f)
	if err != nil {
		log.Fatalf("failed to send body: %v", err)
	}

	// 無くても動く
	_, err = fmt.Fprint(rw, "\x00")
	if err != nil {
		return
	}
}
