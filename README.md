# Raindrop.io MCP Server

[![CI](https://github.com/nnno/my-raindrop-io-mcp-server/actions/workflows/ci.yml/badge.svg)](https://github.com/nnno/my-raindrop-io-mcp-server/actions/workflows/ci.yml)

[Raindrop.io](https://raindrop.io) のブックマークを操作する [Model Context Protocol (MCP)](https://modelcontextprotocol.io) サーバー。Go で実装。

## 提供ツール

### ブックマーク

| ツール | 説明 | アノテーション |
|--------|------|----------------|
| `search_bookmarks` | ブックマークを検索 | ReadOnly |
| `get_bookmark` | ID を指定して詳細取得 | ReadOnly |
| `create_bookmark` | 新しいブックマークを作成 | — |
| `update_bookmark` | 既存のブックマークを更新 | — |
| `delete_bookmark` | ブックマークを削除 | Destructive |

### コレクション

| ツール | 説明 | アノテーション |
|--------|------|----------------|
| `list_collections` | コレクション一覧を取得 | ReadOnly |
| `create_collection` | 新しいコレクションを作成 | — |
| `update_collection` | 既存のコレクションを更新 | — |

## セットアップ

### 前提条件

- [Raindrop.io のテストトークン](https://developer.raindrop.io/v1/authentication/token)

### インストール

インストールスクリプト（`~/.local/bin/` に配置）:

```sh
curl -fsSL https://raw.githubusercontent.com/nnno/my-raindrop-io-mcp-server/main/install.sh | sh
```

Go ユーザーの場合:

```sh
go install github.com/nnno/my-raindrop-io-mcp-server@latest
```

ソースからビルドする場合:

```sh
make install
```

### MCP クライアントへの設定

Claude Desktop (`claude_desktop_config.json`) の場合:

```json
{
  "mcpServers": {
    "raindrop": {
      "command": "~/.local/bin/my-raindrop-io-mcp-server",
      "env": {
        "RAINDROP_TOKEN": "your-test-token"
      }
    }
  }
}
```

Claude Code の場合:

```sh
claude mcp add raindrop ~/.local/bin/my-raindrop-io-mcp-server -e RAINDROP_TOKEN=your-test-token
```

## アーキテクチャ

DDD ベースのレイヤードアーキテクチャを採用。

```
main.go                          # エントリポイント（DI・stdio 起動）
internal/
├── domain/
│   ├── entity/                  # Bookmark, Collection, DomainError
│   └── repository/              # インターフェース定義
│       └── mock/                # mockgen 生成モック
├── application/
│   └── usecase/                 # バリデーション・ビジネスロジック
├── infra/
│   └── raindrop/                # Raindrop.io REST API クライアント
└── handler/                     # MCP ツール登録・リクエスト変換
```

依存方向: `handler → usecase → (domain/entity + domain/repository) ← infra/raindrop`

## 開発

```sh
make build       # ビルド（バージョン情報を埋め込み）
make test        # テスト
make test-race   # テスト（レース検出付き）
make lint        # リント
make mock        # モック再生成
```

## ライセンス

MIT
