#!/bin/bash

set -e  # エラーが発生したら停止

echo "=== 本番環境へのデプロイ開始 ==="
echo ""

# サービス停止
echo "📦 trends-summaryサービスを停止中..."
sudo systemctl stop trends-summary

echo "🔒 Caddyを停止中..."
sudo systemctl stop caddy

# バイナリをコピー
echo "📁 バイナリをコピー中..."
sudo cp trends-summary /opt/trends-summary/trends-summary

# staticディレクトリをコピー（Reactビルド成果物含む）
echo "📁 staticディレクトリをコピー中..."
sudo rm -rf /opt/trends-summary/static
sudo cp -r static /opt/trends-summary/

# 必要な権限を設定
echo "🔐 権限を設定中..."
sudo setcap 'cap_net_bind_service=+ep' /opt/trends-summary/trends-summary

# サービス開始
echo "🚀 trends-summaryサービスを起動中..."
sudo systemctl start trends-summary

echo "🔒 Caddyを起動中..."
sudo systemctl start caddy

# ステータス確認
echo ""
echo "✅ デプロイ完了！"
echo ""
echo "📊 サービスステータス："
echo ""
echo "=== trends-summary ==="
sudo systemctl status trends-summary --no-pager
echo ""
echo "=== Caddy ==="
sudo systemctl status caddy --no-pager

echo ""
echo "🌐 アクセスURL: ${MY_DOMAIN_URL}"
echo "📝 ログ確認: sudo journalctl -u trends-summary -f"
