# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

webview_goは、軽量なクロスプラットフォームWebViewライブラリのGoバインディングです。GoアプリケーションにネイティブWebViewを組み込み、HTML/CSS/JavaScriptでUIを構築しながら、バックエンドロジックをGoで実装できます。

## アーキテクチャ

### 主要コンポーネント
- **webview.go**: WebViewインターフェースとCgoバインディング
- **glue.c**: GoコールバックとCコールバック間の変換層
- **libs/webview/**: ネイティブwebview実装（C/C++）
- **libs/mswebview2/**: Windows用Edge WebView2ヘッダー

### プラットフォーム実装
- Linux: GTK+ 3.0 + WebKit2GTK 4.0
- macOS: WebKit framework
- Windows: Microsoft Edge WebView2

## 開発コマンド

### ビルド
```bash
# 通常のビルド
go build

# Windows用（コンソールウィンドウなし）
go build -ldflags="-H windowsgui"
```

### テスト
```bash
# 全テスト実行
go test

# 詳細出力（実際にウィンドウが表示される）
go test -v

# 特定のテスト実行
go test -run TestWebView
```

### 依存関係
Linux環境では以下のパッケージが必要:
```bash
sudo apt-get install libgtk-3-dev libwebkit2gtk-4.0-dev
```

## 重要な実装パターン

### GoとJavaScript間の通信
1. **Bind()**: Go関数をJavaScriptに公開
2. **Eval()**: JavaScriptコードを実行
3. **Dispatch()**: Goコードをメインスレッドで実行

### スレッドセーフティ
- WebViewの操作は全てメインスレッドで実行される必要がある
- `Run()`メソッドが自動的に`runtime.LockOSThread()`を呼ぶ
- `Dispatch()`を使用して他のgoroutineからメインスレッドに処理を委譲

### エラーハンドリング
- WebView作成の失敗は`New()`がnilを返すことで示される
- JavaScript実行エラーは`Eval()`の戻り値でチェック
- Bind関数内のエラーは適切なJSONレスポンスで返す

## コード実装時の注意点

1. **初期化順序**: `Run()`を呼ぶ前に`Eval()`や`Dispatch()`を呼ばない
2. **メモリ管理**: `Destroy()`を適切に呼んでリソースを解放
3. **JSON通信**: Bind関数の引数と戻り値はJSON形式でやり取りされる
4. **デバッグモード**: `SetDebug(true)`でDevToolsを有効化できる