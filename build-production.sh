#!/bin/bash

set -e  # エラーが発生したら停止

echo "=== Reactフロントエンドのビルド開始 ==="
cd frontend
npm run build
cd ..

echo "=== ビルド成果物をstaticディレクトリにコピー ==="
# 既存のstatic内のReact関連ファイルを削除（古いCSSやJSは残す）
rm -rf static/assets
rm -f static/index.html

# ビルドされたファイルをstaticにコピー
cp -r frontend/dist/* static/

echo "=== Goバイナリのビルド ==="
go build -o trends-summary

echo ""
echo "✅ ビルド完了！"
echo ""
echo "本番環境へのデプロイコマンド："
echo "sudo systemctl stop trends-summary"
echo "sudo cp trends-summary /opt/trends-summary/trends-summary"
echo "sudo setcap 'cap_net_bind_service=+ep' /opt/trends-summary/trends-summary"
echo "sudo systemctl start trends-summary"
