# CLAUDE.md

## プロジェクト概要

Raindrop.io のブックマーク・コレクションを操作する MCP (Model Context Protocol) サーバー。Go 実装。

## 開発環境

- Go 1.25（mise で管理: `mise.toml`）
- golangci-lint（mise で latest を管理）
- モック生成: `go.uber.org/mock/mockgen`
- MCP ライブラリ: `github.com/mark3labs/mcp-go`

## よく使うコマンド

```sh
# ビルド
go build ./...

# テスト（全体）
go test ./...

# テスト（レース検出付き）
go test -race ./...

# リント
golangci-lint run

# モック再生成
go run go.uber.org/mock/mockgen@latest \
  -source=internal/domain/repository/bookmark.go \
  -destination=internal/domain/repository/mock/bookmark_mock.go \
  -package=mock

go run go.uber.org/mock/mockgen@latest \
  -source=internal/domain/repository/collection.go \
  -destination=internal/domain/repository/mock/collection_mock.go \
  -package=mock
```

## アーキテクチャ

DDD ベースのレイヤードアーキテクチャ。依存方向: `handler → usecase → (domain/entity + domain/repository) ← infra/raindrop`

- **handler/**: MCP ツール登録・リクエスト変換。ツールごとに1ファイル（`create.go`, `search.go` 等）、登録は `register.go` の `RegisterAll` で集約
- **usecase/**: 入力バリデーション・ビジネスロジック
- **domain/entity/**: エンティティ定義とドメインエラー（`ErrorKind` による分類）
- **domain/repository/**: リポジトリインターフェース。`mock/` に mockgen 生成モック
- **infra/raindrop/**: Raindrop.io REST API クライアント実装。`client.go` に共通 HTTP ロジック

## コーディング規約

- gofmt / goimports（ローカルプレフィックス: `github.com/nnno/my-raindrop-io-mcp-server`）で整形
- golangci-lint の standard リンタセット + misspell (US locale) を適用
- テストには `github.com/stretchr/testify` を使用

## CI

GitHub Actions（`.github/workflows/ci.yml`）: テスト → リント → 脆弱性スキャン（govulncheck）
