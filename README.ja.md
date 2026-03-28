# jviz

JSONデータをブラウザでビジュアライズするCLIツールです。

JSONの配列を `jviz` にパイプするとローカルWebページが開き、インタラクティブなチャートを表示します。データが更新されるとチャートもリアルタイムで更新されます。

## インストール

```sh
go install github.com/nlink-jp/jviz@latest
```

または[リリースページ](https://github.com/nlink-jp/jviz/releases)からビルド済みバイナリをダウンロードしてください。

## 使い方

```sh
# 標準入力からJSONをパイプ
cat data.json | jviz

# ストリーミング：新しいJSON配列が来るたびにチャートを更新
while true; do generate-stats.sh; sleep 5; done | jviz

# ファイルの変更を監視
jviz --watch data.json

# ポートを変更（デフォルト: 8765）
jviz --port 9000 < data.json

# ブラウザの自動起動を無効化
jviz --no-open < data.json
```

### 入力形式

`jviz` はオブジェクトのJSON配列を受け付けます：

```json
[
  {"label": "1月", "売上": 120},
  {"label": "2月", "売上": 95},
  {"label": "3月", "売上": 140}
]
```

### チャートタイプ

| タイプ | 説明                  |
|--------|-----------------------|
| Bar    | 棒グラフ              |
| Line   | 折れ線グラフ（塗りつぶしあり）|
| Pie    | 円グラフ              |
| Table  | データテーブル        |

**X / Label** と **Y / Value** のセレクターで描画するカラムを選択できます。

## 組み合わせて使えるツール

- [jstats](https://github.com/nlink-jp/jstats) — JSONをフィールドで集計
- [csv-to-json](https://github.com/nlink-jp/csv-to-json) — CSVをJSONに変換
- [json-filter](https://github.com/nlink-jp/json-filter) — JSON配列をフィルタリング

## ビルド

```sh
make build        # 現在のプラットフォーム → dist/jviz
make build-all    # 全5プラットフォーム  → dist/
make test
```

## ライセンス

MIT — [LICENSE](LICENSE) を参照してください。

Chart.js は独自のMITライセンス（© 2023 Chart.js Contributors）の下でバンドルされています。
