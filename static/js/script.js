// DOMが完全に読み込まれた後にスクリプトを実行する
document.addEventListener('DOMContentLoaded', () => {
    // モーダル要素の取得
    const modal = document.getElementById('summary-modal');
    const closeButton = document.querySelector('.close-button');
    const summaryText = document.getElementById('summary-text');
    const loadingIndicator = document.getElementById('loading-indicator');

    // モーダルを表示する関数
    function showModal(summary) {
        summaryText.textContent = summary;
        loadingIndicator.style.display = 'none';
        summaryText.style.display = 'block';
        modal.style.display = 'block';
    }

    // モーダルを表示中にローディングインジケーターを表示する関数
    function showLoading() {
        summaryText.style.display = 'none';
        loadingIndicator.style.display = 'block';
        modal.style.display = 'block';
    }

    // モーダルを閉じる関数
    function closeModal() {
        modal.style.display = 'none';
    }

    // 閉じるボタンにイベントリスナーを追加
    if (closeButton) {
        closeButton.addEventListener('click', closeModal);
    } else {
        console.error('Close button not found in the DOM.');
    }

    // モーダル外をクリックしたら閉じる
    window.addEventListener('click', (event) => {
        if (event.target == modal) {
            closeModal();
        }
    });

    // RSSフィードのJSONデータを取得
    fetch('/rss')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTPエラー: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            const tbody = document.getElementById('rss-feed-body');

            // データをテーブルに追加
            data.items.forEach(item => {
                const row = document.createElement('tr');

                // タイトルセル
                const titleCell = document.createElement('td');
                titleCell.textContent = item.title || 'No title';
                row.appendChild(titleCell);

                // 説明セル
                const descriptionCell = document.createElement('td');

                // HTMLを挿入
                descriptionCell.innerHTML = item.description || 'No description';

                // 画像のサイズを調整
                const images = descriptionCell.querySelectorAll('img'); // <img>タグを取得
                images.forEach(img => {
                    img.style.maxWidth = '100px';  // 最大幅を指定
                    img.style.maxHeight = '100px'; // 最大高さを指定
                    img.style.objectFit = 'contain'; // 画像のアスペクト比を保ったままフィットさせる
                });

                // セルを行に追加
                row.appendChild(descriptionCell);

                // 公開日セル
                const publishedCell = document.createElement('td');
                publishedCell.textContent = item.published || 'No date';
                row.appendChild(publishedCell);

                // リンクセル
                const linkCell = document.createElement('td');
                const link = document.createElement('a');
                link.href = item.link;
                link.textContent = '記事を見る';
                link.target = '_blank';
                linkCell.appendChild(link);
                row.appendChild(linkCell);

                // ボタンセル（新規追加）
                const buttonCell = document.createElement('td');
                const button = document.createElement('button');
                button.textContent = 'AIサマリー'; // ボタンに表示するテキスト
                button.classList.add('rss-button'); // スタイルクラスを追加（オプション）

                // ボタンクリック時のイベントリスナーを追加
                button.addEventListener('click', () => {
                    // 親行（<tr>）を取得
                    const parentRow = button.parentElement.parentElement;

                    // リンクセル（4番目の<td>）を取得
                    const linkTd = parentRow.children[3];
                    const articleUrl = linkTd.querySelector('a').href;

                    // モーダルを表示してローディングインジケーターを表示
                    showLoading();

                    // /ai-summary API に GET リクエストを送信
                    fetch(`/ai-summary?url=${encodeURIComponent(articleUrl)}`)
                        .then(response => {
                            if (!response.ok) {
                                throw new Error(`APIエラー: ${response.status}`);
                            }
                            return response.json();
                        })
                        .then(summaryData => {
                            // サマリーをモーダルで表示
                            showModal(summaryData.summary);
                        })
                        .catch(error => {
                            console.error('AIサマリー取得エラー:', error);
                            loadingIndicator.style.display = 'none';
                            summaryText.textContent = 'AIサマリーの取得に失敗しました。';
                            summaryText.style.display = 'block';
                        });
                });

                buttonCell.appendChild(button);
                row.appendChild(buttonCell);

                tbody.appendChild(row);
            });
        })
        .catch(error => {
            console.error('Fetch Error:', error);
            const tbody = document.getElementById('rss-feed-body');
            const row = document.createElement('tr');
            const errorCell = document.createElement('td');
            errorCell.colSpan = 5; // ボタン列も含めて5に変更
            errorCell.textContent = 'RSSフィードの取得に失敗しました。';
            row.appendChild(errorCell);
            tbody.appendChild(row);
        });
});

// TIOBEグラフの取得と表示
fetch('/tiobe-graph')
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTPエラー: ${response.status}`);
        }
        return response.json(); // JSONとして受け取る
    })
    .then(data => {
        const container = document.getElementById('tiobe-graph-container');
        container.innerHTML = data; // HTMLとして挿入
    })
    .catch(error => {
        console.error('Fetch Error:', error);
        const container = document.getElementById('tiobe-graph-container');
        container.innerHTML = '<p>グラフの取得に失敗しました。</p>';
    });
