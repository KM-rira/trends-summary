#!/bin/bash

set -e  # エラーが発生したら停止

# 環境変数のチェック
if [ -z "$SSH_KEY_PATH" ]; then
    echo "❌ エラー: SSH_KEY_PATH 環境変数が設定されていません"
    echo "例: export SSH_KEY_PATH=\"~/.ssh/key.txt\""
    exit 1
fi

if [ -z "$TARGET_INSTANCE" ]; then
    echo "❌ エラー: TARGET_INSTANCE 環境変数が設定されていません"
    echo "例: export TARGET_INSTANCE=\"os@domain\""
    exit 1
fi

# チルダ展開
SSH_KEY_PATH="${SSH_KEY_PATH/#\~/$HOME}"

# SSH鍵の存在確認
if [ ! -f "$SSH_KEY_PATH" ]; then
    echo "❌ エラー: SSH鍵が見つかりません: $SSH_KEY_PATH"
    exit 1
fi

echo "=== 本番環境へのファイル転送開始 ==="
echo "🔑 SSH鍵: $SSH_KEY_PATH"
echo "🎯 ターゲット: $TARGET_INSTANCE"
echo ""

# staticディレクトリの転送
echo "📁 staticディレクトリを転送中..."
scp -i "$SSH_KEY_PATH" -r static "$TARGET_INSTANCE:~/repo/trends-summary/"

echo ""

# trends-summaryバイナリの転送
echo "📦 trends-summaryバイナリを転送中..."
scp -i "$SSH_KEY_PATH" trends-summary "$TARGET_INSTANCE:~/repo/trends-summary/"

echo ""
echo "✅ 転送完了！"
echo ""
echo "次のステップ:"
echo "  ssh -i $SSH_KEY_PATH $TARGET_INSTANCE"
echo "  cd ~/repo/trends-summary"
echo "  ./deploy-production.sh"
