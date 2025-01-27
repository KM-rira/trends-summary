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

            const titleCell = document.createElement('td');
            titleCell.textContent = item.title || 'No title';
            row.appendChild(titleCell);

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

            const publishedCell = document.createElement('td');
            publishedCell.textContent = item.published || 'No date';
            row.appendChild(publishedCell);

            const linkCell = document.createElement('td');
            const link = document.createElement('a');
            link.href = item.link;
            link.textContent = '記事を見る';
            link.target = '_blank';
            linkCell.appendChild(link);
            row.appendChild(linkCell);

            tbody.appendChild(row);
        });
    })
    .catch(error => {
        console.error('Fetch Error:', error);
        const tbody = document.getElementById('rss-feed-body');
        const row = document.createElement('tr');
        const errorCell = document.createElement('td');
        errorCell.colSpan = 4;
        errorCell.textContent = 'RSSフィードの取得に失敗しました。';
        row.appendChild(errorCell);
        tbody.appendChild(row);
    });

        // GitHubトレンドデータのJSONを取得
    fetch('/github-trending')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTPエラー: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            const tbody = document.getElementById('github-trending-body');

            // データをテーブルに追加
            data.forEach(item => {
                const row = document.createElement('tr');

                // リポジトリ名
                const nameCell = document.createElement('td');
                nameCell.textContent = item.name || 'No name';
                row.appendChild(nameCell);

                // 説明
                const descriptionCell = document.createElement('td');
                descriptionCell.textContent = item.description || 'No description';
                row.appendChild(descriptionCell);

                // 言語
                const languageCell = document.createElement('td');
                languageCell.textContent = item.language || 'N/A';
                row.appendChild(languageCell);

                // スター数
                const starsCell = document.createElement('td');
                starsCell.textContent = item.stars || '0';
                row.appendChild(starsCell);

                // リンク
                const linkCell = document.createElement('td');
                const link = document.createElement('a');
                link.href = item.url;
                link.textContent = 'リポジトリを見る';
                link.target = '_blank';
                linkCell.appendChild(link);
                row.appendChild(linkCell);

                tbody.appendChild(row);
            });

            // JSON全データ表示用
            const rawOutput = document.getElementById('github-trending-raw');
            rawOutput.textContent = JSON.stringify(data, null, 2);
        })
        .catch(error => {
            console.error('Fetch Error:', error);

            // テーブルのエラー表示
            const tbody = document.getElementById('github-trending-body');
            const row = document.createElement('tr');
            const errorCell = document.createElement('td');
            errorCell.colSpan = 5;
            errorCell.textContent = 'GitHubトレンドデータの取得に失敗しました。';
            row.appendChild(errorCell);
            tbody.appendChild(row);

            // JSON全データのエラー表示
            const rawOutput = document.getElementById('github-trending-raw');
            rawOutput.textContent = 'エラー: データの取得に失敗しました。';
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
