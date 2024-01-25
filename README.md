# private-isu-bench-lambda

「[ISUCON](https://isucon.net/)」は、LINE株式会社の商標または登録商標です。

本リポジトリは[private-isu](https://github.com/catatsuy/private-isu)のベンチマーカーを[AWS Lambda](https://aws.amazon.com/jp/lambda/)で動かすためのリポジトリです。

チーム情報をスプレッドシートで管理し、Lambda関数URLをcurlで叩くことでベンチマーカーを動かし、その結果を[Mackerel](https://ja.mackerel.io/)に送信します。

## 使い方
1. [private-isu/benchmarker/](https://github.com/catatsuy/private-isu/tree/master/benchmarker)をビルドし、そのファイルをbinディレクトリに配置する。
```sh
GOOS=linux GOARCH=amd64 go build -o benchmarker
```
2. userdataを[こちら](https://github.com/catatsuy/private-isu/tree/master#mac%E3%82%84linux%E4%B8%8A%E3%81%A7%E9%81%A9%E5%BD%93%E3%81%AB%E5%8B%95%E3%81%8B%E3%81%99)を参照して必要なファイルを配置する。Lambdaにデプロイする際のデプロイパッケージの制限を超えないように注意してください。
3. スプレッドシートを作成し、スプレッドシートIDを取得する。
4. Googleの認証情報Jsonファイルを取得する。
5. MackerelでサービスとAPIKeyを作成する。
6. 本リポジトリのmain.goをビルドし、ファイルをzip化し、S3にアップロードする。
7. Lambdaの「コードソース」の「アップロード元」から先ほどアップロードしたzipファイルのURLを入力する。
8. Lambdaで環境変数を設定する。

| Environment Variable             | Description   |
|----------------------------------|---------------|
| `MACKEREL_API_KEY`              | Mackerelから発行されるAPIキー |
| `MACKEREL_SERVICE_NAME`         | Mackerelのサービス名          |
| `SPREADSHEETID`                | GoogleスプレッドシートのID    |
| `SPREADSHEET_CREDENTIALS_JSON` | Google APIの認証情報を含むJSON（サービスアカウントからダウンロード）|
| `SPREADSHEET_RANGE`            | チーム情報が記載されているシート名と範囲（シート1のA2からB50の場合。`シート1!A2:B50`） |

9.  private-isuを動かしているサーバー上で、Lambda関数のURLを叩く。
```sh
curl <your-lambda-function-url>
# レスポンス例
{"pass":true,"score":623,"success":565,"fail":2,"messages":["リクエストがタイムアウトしました (POST /login)"]}
```

### [lambroll](https://github.com/fujiwara/lambroll) を活用したデプロイ
1~5については上記と同様。Makefileで以下の環境変数を設定する。

| Environment Variable             | Description   |
|----------------------------------|---------------|
| `MACKEREL_API_KEY`              | Mackerelから発行されるAPIキー |
| `MACKEREL_SERVICE_NAME`         | Mackerelのサービス名          |
| `SPREADSHEETID`                | GoogleスプレッドシートのID    |
| `SPREADSHEET_CREDENTIALS_JSON` | Google APIの認証情報を含むJSON（サービスアカウントからダウンロード）|
| `SPREADSHEET_RANGE`            | チーム情報が記載されているシート名と範囲（シート1のA2からB50の場合。`シート1!A2:B50`） |
| `FUNCTION_NAME`                 | AWS Lambda関数の名前  |
| `FUNCTION_ROLE`                 | 関数を実行するときのAWSロール |
| `S3_BUCKET`                    | 関数のデプロイ先のS3バケット  |
| `S3_KEY`                        | デプロイするときのS3のキー    |

その後`make deploy_private_isu_bench_lambda` を実行し、デプロイが完了する。

## ディレクトリ構成
```
├── bin  # private-isuのベンチマーカーのバイナリ
├── userdata  # ベンチマーカーで必要な画像などのデータ
└── main.go  # Lambda上で動かすメインプログラム
```

## スプレッドシートについて
左側にIPアドレス、右側にチーム名になるようにしてください。

例

| IPアドレス | チーム名 |
|-----------|---------|
| 192.0.2.1 | test1   |

## 参考ブログ
- [デジタル創作同好会traPさんと社内ISUCONイベントを開催しました](https://developers.prtimes.jp/2023/01/19/private-isu-with-trap-2023/)
