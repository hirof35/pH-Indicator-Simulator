# pH Indicator Simulator (Go + Lorca)

Go言語のバックエンド処理と、HTML5/CSS3/JavaScriptによるフロントエンドGUIを組み合わせた、インタラクティブな液体属性（酸性・中性・アルカリ性）判定シミュレーターです。

Cコンパイラ（cgo）への依存を排除しているため、Windows環境でも軽量かつ安全に動作します。
<img width="1032" height="903" alt="スクリーンショット 2026-05-24 190726" src="https://github.com/user-attachments/assets/b29d3616-8ffb-40fd-87cd-6bb16a32f740" />

## 🧪 機能特徴

- **リアルタイムシミュレーション**: スライダーを動かすと、pH（0.0〜14.0）の値が即座に反映されます。
- **各種指示薬の完全再現**:
  - **赤色リトマス紙**: アルカリ性（pH 8.0以上）で青色に変色。
  - **青色リトマス紙**: 酸性（pH 5.0以下）で赤色に変色。
  - **BTB溶液**: 酸性で黄色、中性で緑色、アルカリ性で青色に変色。
  - **フェノールフタレイン溶液**: アルカリ性（pH 8.3以上）で赤紫色に変色し、強アルカリになるほど色が濃くなります。
- **ハイブリッド構造**: フロントエンドの描画をブラウザ（Chrome/Edge）に任せ、液性の判定ロジックはGo言語側で厳密に計算（`window.goJudgePH` を通じてシームレスに双方向通信）しています。

## 🛠️ 動作環境

- **Go言語**: 1.16以上推奨
- **Google Chrome** または **Microsoft Edge**（Chromium版）がシステムにインストールされていること
- **OS**: Windows, macOS, Linux (CGO不要のPure Goで動作)

## 📦 必要外部パッケージ

本アプリは軽量GUIライブラリ `Lorca` を使用しています。

```bash
go get [github.com/zserge/lorca](https://github.com/zserge/lorca)
🚀 起動方法
Windows (PowerShell) 環境では、cgoを無効化（CGO_ENABLED="0"）して実行します。

PowerShell
$env:CGO_ENABLED="0"
go run main.go
📂 プロジェクト構造
Plaintext
go-webview-ph/
├── main.go       # Goロジック ＆ 画面HTML/CSS/JSデータ
├── go.mod        # Goモジュール定義ファイル
├── go.sum        # 依存パッケージのチェックサム
└── README.md     # 本ドキュメント
💡 技術的なこだわり
CGOフリーによる高移植性:
webview_go などで発生しがちなWindows上でのC++コンパイラ（MinGW等）のバージョンエラーを回避するため、内部通信にWebSocketを採用した Lorca を選定。これにより、安全かつ一瞬でコンパイル・実行が可能です。

Base64エンコードによる安全な画面読み込み:
HTML内のCSSやJavaScriptで使われる記法（rgba()のカンマや、文字列結合の+など）が、Go言語側やURLエスケープ側で誤解釈されて構文エラーを起こすのを防ぐため、HTMLリソース全体をBase64文字列に暗号化してから安全にブラウザエンジンへと流し込んでいます。

📄 ライセンス
MIT License
