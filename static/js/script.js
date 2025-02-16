document.addEventListener("DOMContentLoaded", () => {
  // モーダル要素の取得
  const modal = document.getElementById("summary-modal");
  const closeButton = document.querySelector(".close-button");
  const summaryText = document.getElementById("summary-text");
  const loadingIndicator = document.getElementById("loading-indicator");

  // モーダルを表示する関数
  function showModal(summary) {
    summaryText.textContent = summary;
    loadingIndicator.style.display = "none";
    summaryText.style.display = "block";
    modal.style.display = "block";
  }

  // モーダルを表示中にローディングインジケーターを表示する関数
  function showLoading() {
    summaryText.style.display = "none";
    loadingIndicator.style.display = "block";
    modal.style.display = "block";
  }

  // モーダルを閉じる関数
  function closeModal() {
    modal.style.display = "none";
  }

  // 閉じるボタンにイベントリスナーを追加
  if (closeButton) {
    closeButton.addEventListener("click", closeModal);
  } else {
    console.error("Close button not found in the DOM.");
  }

  // モーダル外をクリックしたら閉じる
  window.addEventListener("click", (event) => {
    if (event.target == modal) {
      closeModal();
    }
  });

  // RSSフィードのJSONデータを取得
  fetch("/rss")
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTPエラー: ${response.status}`);
      }
      return response.json();
    })
    .then((data) => {
      const tbody = document.getElementById("rss-feed-body");

      // データをテーブルに追加
      data.items.forEach((item) => {
        const row = document.createElement("tr");

        // タイトルセル
        const titleCell = document.createElement("td");
        titleCell.textContent = item.title || "No title";
        row.appendChild(titleCell);

        // 説明セル
        const descriptionCell = document.createElement("td");

        // HTMLを挿入
        descriptionCell.innerHTML = item.description || "No description";

        // 画像のサイズを調整
        const images = descriptionCell.querySelectorAll("img"); // <img>タグを取得
        images.forEach((img) => {
          img.style.maxWidth = "100px"; // 最大幅を指定
          img.style.maxHeight = "100px"; // 最大高さを指定
          img.style.objectFit = "contain"; // 画像のアスペクト比を保ったままフィットさせる
        });

        // セルを行に追加
        row.appendChild(descriptionCell);

        // 公開日セル
        const publishedCell = document.createElement("td");
        publishedCell.textContent = item.published || "No date";
        row.appendChild(publishedCell);

        // リンクセル
        const linkCell = document.createElement("td");
        const link = document.createElement("a");
        link.href = item.link;
        link.textContent = "URL";
        link.target = "_blank";
        linkCell.appendChild(link);
        row.appendChild(linkCell);

        // ボタンセル（新規追加）
        const buttonCell = document.createElement("td");
        const button = document.createElement("button");
        button.textContent = "Generate"; // ボタンに表示するテキスト
        button.classList.add("rss-button"); // スタイルクラスを追加（オプション）

        // ボタンクリック時のイベントリスナーを追加
        button.addEventListener("click", () => {
          // 親行（<tr>）を取得
          const parentRow = button.parentElement.parentElement;

          // リンクセル（4番目の<td>）を取得
          const linkTd = parentRow.children[3];
          const articleUrl = linkTd.querySelector("a").href;

          // モーダルを表示してローディングインジケーターを表示
          showLoading();

          // /ai-article-summary API に GET リクエストを送信
          fetch(`/ai-article-summary?url=${encodeURIComponent(articleUrl)}`)
            .then((response) => {
              if (!response.ok) {
                throw new Error(`APIエラー: ${response.status}`);
              }
              return response.json();
            })
            .then((summaryData) => {
              // サマリーをモーダルで表示
              showModal(summaryData.summary);
            })
            .catch((error) => {
              console.error("AIサマリー取得エラー:", error);
              loadingIndicator.style.display = "none";
              summaryText.textContent = "AIサマリーの取得に失敗しました。";
              summaryText.style.display = "block";
            });
        });

        buttonCell.appendChild(button);
        row.appendChild(buttonCell);

        tbody.appendChild(row);
      });
    })
    .catch((error) => {
      console.error("Fetch Error:", error);
      const tbody = document.getElementById("rss-feed-body");
      const row = document.createElement("tr");
      const errorCell = document.createElement("td");
      errorCell.colSpan = 5; // ボタン列も含めて5に変更
      errorCell.textContent = "RSSフィードの取得に失敗しました。";
      row.appendChild(errorCell);
      tbody.appendChild(row);
    });

  // GitHubトレンドデータのJSONを取得
  fetch("/github-trending")
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTPエラー: ${response.status}`);
      }
      return response.json();
    })
    .then((data) => {
      const tbody = document.getElementById("github-trending-body");

      // データをテーブルに追加
      data.forEach((item) => {
        const row = document.createElement("tr");

        // リポジトリ名
        const nameCell = document.createElement("td");
        nameCell.textContent = item.name || "No name";
        row.appendChild(nameCell);

        // 説明
        const descriptionCell = document.createElement("td");
        descriptionCell.textContent = item.description || "No description";
        row.appendChild(descriptionCell);

        // 言語
        const languageCell = document.createElement("td");
        languageCell.textContent = item.language || "N/A";
        row.appendChild(languageCell);

        // スター数
        const starsCell = document.createElement("td");
        starsCell.textContent = item.stars || "0";
        row.appendChild(starsCell);

        // リンク
        const linkCell = document.createElement("td");
        const link = document.createElement("a");
        link.href = item.url;
        link.textContent = "URL";
        link.target = "_blank";
        linkCell.appendChild(link);
        row.appendChild(linkCell);

        // ボタンセル（新規追加）
        const buttonCell = document.createElement("td");
        const button = document.createElement("button");
        button.textContent = "Generate"; // ボタンに表示するテキスト
        button.classList.add("rss-button"); // スタイルクラスを追加（オプション）

        // ボタンクリック時のイベントリスナーを追加
        button.addEventListener("click", () => {
          // 親行（<tr>）を取得
          const parentRow = button.parentElement.parentElement;

          // リンクセル（5番目の<td>）を取得
          const linkTd = parentRow.children[4];
          const repositoryUrl = linkTd.querySelector("a").href;

          // モーダルを表示してローディングインジケーターを表示
          showLoading();

          // /ai-repository-summary API に GET リクエストを送信
          fetch(
            `/ai-repository-summary?url=${encodeURIComponent(repositoryUrl)}`,
          )
            .then((response) => {
              if (!response.ok) {
                throw new Error(`APIエラー: ${response.status}`);
              }
              return response.json();
            })
            .then((summaryData) => {
              // サマリーをモーダルで表示
              showModal(summaryData.summary);
            })
            .catch((error) => {
              console.error("AIサマリー取得エラー:", error);
              loadingIndicator.style.display = "none";
              summaryText.textContent = "AIサマリーの取得に失敗しました。";
              summaryText.style.display = "block";
            });
        });

        buttonCell.appendChild(button);
        row.appendChild(buttonCell);

        tbody.appendChild(row);
      });

      // JSON全データ表示用
      // const rawOutput = document.getElementById('github-trending-raw');
      // rawOutput.textContent = JSON.stringify(data, null, 2);
    })
    .catch((error) => {
      console.error("Fetch Error:", error);

      // テーブルのエラー表示
      const tbody = document.getElementById("github-trending-body");
      const row = document.createElement("tr");
      const errorCell = document.createElement("td");
      errorCell.colSpan = 6;
      errorCell.textContent = "GitHubトレンドデータの取得に失敗しました。";
      row.appendChild(errorCell);
      tbody.appendChild(row);

      // JSON全データのエラー表示
      // const rawOutput = document.getElementById('github-trending-raw');
      // rawOutput.textContent = 'エラー: データの取得に失敗しました。';
    });

  // Golang RepositotyトレンドデータのJSONを取得
  fetch("/golang-repository-trending")
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTPエラー: ${response.status}`);
      }
      return response.json();
    })
    .then((data) => {
      const tbody = document.getElementById("golang-repository-trending-body");

      // データをテーブルに追加
      data.forEach((item) => {
        const row = document.createElement("tr");

        // リポジトリ名
        const nameCell = document.createElement("td");
        nameCell.textContent = item.name || "No name";
        row.appendChild(nameCell);

        // 説明
        const descriptionCell = document.createElement("td");
        descriptionCell.textContent = item.description || "No description";
        row.appendChild(descriptionCell);

        // 言語
        const languageCell = document.createElement("td");
        languageCell.textContent = item.language || "N/A";
        row.appendChild(languageCell);

        // スター数
        const starsCell = document.createElement("td");
        starsCell.textContent = item.stars || "0";
        row.appendChild(starsCell);

        // リンク
        const linkCell = document.createElement("td");
        const link = document.createElement("a");
        link.href = item.url;
        link.textContent = "URL";
        link.target = "_blank";
        linkCell.appendChild(link);
        row.appendChild(linkCell);

        // ボタンセル（新規追加）
        const buttonCell = document.createElement("td");
        const button = document.createElement("button");
        button.textContent = "Generate"; // ボタンに表示するテキスト
        button.classList.add("rss-button"); // スタイルクラスを追加（オプション）

        // ボタンクリック時のイベントリスナーを追加
        button.addEventListener("click", () => {
          // 親行（<tr>）を取得
          const parentRow = button.parentElement.parentElement;

          // リンクセル（5番目の<td>）を取得
          const linkTd = parentRow.children[4];
          const repositoryUrl = linkTd.querySelector("a").href;

          // モーダルを表示してローディングインジケーターを表示
          showLoading();

          // /ai-repository-summary API に GET リクエストを送信
          fetch(
            `/ai-repository-summary?url=${encodeURIComponent(repositoryUrl)}`,
          )
            .then((response) => {
              if (!response.ok) {
                throw new Error(`APIエラー: ${response.status}`);
              }
              return response.json();
            })
            .then((summaryData) => {
              // サマリーをモーダルで表示
              showModal(summaryData.summary);
            })
            .catch((error) => {
              console.error("AIサマリー取得エラー:", error);
              loadingIndicator.style.display = "none";
              summaryText.textContent = "AIサマリーの取得に失敗しました。";
              summaryText.style.display = "block";
            });
        });

        buttonCell.appendChild(button);
        row.appendChild(buttonCell);

        tbody.appendChild(row);
      });

      // JSON全データ表示用
      // const rawOutput = document.getElementById('github-trending-raw');
      // rawOutput.textContent = JSON.stringify(data, null, 2);
    })
    .catch((error) => {
      console.error("Fetch Error:", error);

      // テーブルのエラー表示
      const tbody = document.getElementById("golang-repository-trending-body");
      const row = document.createElement("tr");
      const errorCell = document.createElement("td");
      errorCell.colSpan = 6;
      errorCell.textContent = "GitHubトレンドデータの取得に失敗しました。";
      row.appendChild(errorCell);
      tbody.appendChild(row);

      // JSON全データのエラー表示
      // const rawOutput = document.getElementById('github-trending-raw');
      // rawOutput.textContent = 'エラー: データの取得に失敗しました。';
    });

  // TIOBEグラフの取得と表示
  fetch("/tiobe-graph")
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTPエラー: ${response.status}`);
      }
      return response.json(); // JSONとして受け取る
    })
    .then((data) => {
      const container = document.getElementById("tiobe-graph-container");
      container.innerHTML = data; // HTMLとして挿入
    })
    .catch((error) => {
      console.error("Fetch Error:", error);
      const container = document.getElementById("tiobe-graph-container");
      container.innerHTML = "<p>グラフの取得に失敗しました。</p>";
    });

  const summaryBox = document.getElementById("summary-box");
  const generateButton = document.getElementById("generate-summary-button");

  generateButton.addEventListener("click", function () {
    // 現在のページのHTMLを取得
    const pageContent = document.body.innerHTML.trim();

    fetch(`/ai-trends-summary`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ data: pageContent }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(`APIエラー: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        const markdownText = data.summary || "AIサマリーの取得に失敗しました。";
        const htmlContent = marked.parse(markdownText);
        summaryBox.innerHTML = htmlContent;
      })
      .catch((error) => {
        console.error("AIサマリー取得エラー:", error);
        summaryBox.textContent = "AIサマリーの取得に失敗しました。";
      });
  });
  function updateCharacterCounts() {
    const rssFeedBody = document.getElementById("rss-feed-body");
    const githubTrendingBody = document.getElementById("github-trending-body");
    const tiobeGraphContainer = document.getElementById(
      "tiobe-graph-container",
    );

    const rssFeedCount = rssFeedBody.innerHTML.length;
    const githubTrendingCount = githubTrendingBody.innerHTML.length;
    const tiobeGraphCount = tiobeGraphContainer.innerHTML.length;

    document.getElementById("rss-feed-count").textContent = rssFeedCount;
    document.getElementById("github-trending-count").textContent =
      githubTrendingCount;
    document.getElementById("tiobe-graph-count").textContent = tiobeGraphCount;
  }

  // データ挿入後に呼び出す
  updateCharacterCounts();

  // 監視対象の要素を取得
  const rssFeedBody = document.getElementById("rss-feed-body");
  const githubTrendingBody = document.getElementById("github-trending-body");
  const tiobeGraphContainer = document.getElementById("tiobe-graph-container");

  // MutationObserverのコールバック関数
  const observerCallback = function (mutationsList, observer) {
    // 変更があった場合にカウントを更新
    updateCharacterCounts();
  };

  // オブザーバのオプション
  const observerOptions = {
    childList: true, // 直接の子要素の追加・削除を監視
    subtree: true, // 全ての子孫要素を監視
    characterData: true, // テキストノードの変更を監視
  };

  // 各要素に対してオブザーバを設定
  const rssObserver = new MutationObserver(observerCallback);
  rssObserver.observe(rssFeedBody, observerOptions);

  const githubObserver = new MutationObserver(observerCallback);
  githubObserver.observe(githubTrendingBody, observerOptions);

  const tiobeObserver = new MutationObserver(observerCallback);
  const container = document.getElementById("golangweekly-container");

  fetch("/golang-weekly-content")
    .then((response) => {
      if (!response.ok) {
        throw new Error(`Network error: ${response.status}`);
      }
      return response.text();
    })
    .then((xmlText) => {
      const parser = new DOMParser();
      const xmlDoc = parser.parseFromString(xmlText, "application/xml");

      const items = xmlDoc.querySelectorAll("item");
      if (items.length === 0) {
        container.textContent = "No articles found.";
        return;
      }

      const list = document.createElement("ul");
      items.forEach((item) => {
        const title = item.querySelector("title")?.textContent || "No Title";
        const link = item.querySelector("link")?.textContent || "#";
        const pubDate = item.querySelector("pubDate")?.textContent || "No Date";
        const description =
          item.querySelector("description")?.textContent || "No Description";

        const listItem = document.createElement("li");
        listItem.innerHTML = `
                    <h3><a href="${link}" target="_blank">${title}</a></h3>
                    <p><strong>Published:</strong> ${pubDate}</p>
                    <p>${description}</p>
                    <hr>
                `;
        list.appendChild(listItem);
      });

      container.innerHTML = ""; // 既存の内容をクリア
      container.appendChild(list);
    })
    .catch((error) => {
      console.error("Error fetching GolangWeekly feed:", error);
      container.textContent = "Error loading feed";
    });

  const gcpContainer = document.getElementById("gcp-rss-container");

  fetch("/google-cloud-content")
    .then((response) => {
      if (!response.ok) {
        throw new Error(`Network error: ${response.status}`);
      }
      return response.text();
    })
    .then((xmlText) => {
      const parser = new DOMParser();
      const xmlDoc = parser.parseFromString(xmlText, "application/xml");

      const items = xmlDoc.querySelectorAll("item");
      if (items.length === 0) {
        gcpContainer.textContent = "No articles found.";
        return;
      }

      const list = document.createElement("ul");
      items.forEach((item) => {
        const title = item.querySelector("title")?.textContent || "No Title";
        const link = item.querySelector("link")?.textContent || "#";
        const pubDate = item.querySelector("pubDate")?.textContent || "No Date";
        const description =
          item.querySelector("description")?.textContent || "No Description";

        const listItem = document.createElement("li");
        listItem.innerHTML = `
                    <h3><a href="${link}" target="_blank">${title}</a></h3>
                    <p><strong>Published:</strong> ${pubDate}</p>
                    <p>${description}</p>
                    <hr>
                `;
        list.appendChild(listItem);
      });

      gcpContainer.innerHTML = ""; // 既存の内容をクリア
      gcpContainer.appendChild(list);
    })
    .catch((error) => {
      console.error("Error fetching GCP feed:", error);
      gcpContainer.textContent = "Error loading feed";
    });

  const awsContainer = document.getElementById("aws-rss-container");

  fetch("/aws-content")
    .then((response) => {
      if (!response.ok) {
        throw new Error(`Network error: ${response.status}`);
      }
      return response.text();
    })
    .then((xmlText) => {
      const parser = new DOMParser();
      const xmlDoc = parser.parseFromString(xmlText, "application/xml");
      const items = xmlDoc.querySelectorAll("item");

      if (items.length === 0) {
        awsContainer.textContent = "No articles found.";
        return;
      }

      const list = document.createElement("ul");
      items.forEach((item) => {
        // 各要素の存在チェックを行いながらテキストを抽出
        const titleElem = item.querySelector("title");
        const linkElem = item.querySelector("link");
        const pubDateElem = item.querySelector("pubDate");
        const descElem = item.querySelector("description");

        const title = titleElem ? titleElem.textContent : "No Title";
        const link = linkElem ? linkElem.textContent : "#";
        const pubDate = pubDateElem ? pubDateElem.textContent : "No Date";
        const description = descElem ? descElem.textContent : "No Description";

        const listItem = document.createElement("li");
        listItem.innerHTML = `
        <h3><a href="${link}" target="_blank">${title}</a></h3>
        <p><strong>Published:</strong> ${pubDate}</p>
        <p>${description}</p>
        <hr>
      `;
        list.appendChild(listItem);
      });

      awsContainer.innerHTML = ""; // 既存の内容をクリア
      awsContainer.appendChild(list);
    })
    .catch((error) => {
      console.error("Error fetching AWS feed:", error);
      awsContainer.textContent = "Error loading feed";
    });
  const azureContainer = document.getElementById("azure-rss-container");

  fetch("/azure-content")
    .then((response) => {
      if (!response.ok) {
        throw new Error(`Network error: ${response.status}`);
      }
      return response.text();
    })
    .then((xmlText) => {
      const parser = new DOMParser();
      const xmlDoc = parser.parseFromString(xmlText, "application/xml");
      const items = xmlDoc.querySelectorAll("item");

      if (items.length === 0) {
        azureContainer.textContent = "No articles found.";
        return;
      }

      const list = document.createElement("ul");
      items.forEach((item) => {
        // 各要素の存在チェックを行いながらテキストを抽出
        const titleElem = item.querySelector("title");
        const linkElem = item.querySelector("link");
        const pubDateElem = item.querySelector("pubDate");
        const descElem = item.querySelector("description");

        const title = titleElem ? titleElem.textContent : "No Title";
        const link = linkElem ? linkElem.textContent : "#";
        const pubDate = pubDateElem ? pubDateElem.textContent : "No Date";
        const description = descElem ? descElem.textContent : "No Description";

        const listItem = document.createElement("li");
        listItem.innerHTML = `
        <h3><a href="${link}" target="_blank">${title}</a></h3>
        <p><strong>Published:</strong> ${pubDate}</p>
        <p>${description}</p>
        <hr>
      `;
        list.appendChild(listItem);
      });

      azureContainer.innerHTML = ""; // 既存の内容をクリア
      azureContainer.appendChild(list);
    })
    .catch((error) => {
      console.error("Error fetching AZURE feed:", error);
      azureContainer.textContent = "Error loading feed";
    });
});
