import type { GitHubTrendingItem, SummaryResponse } from '../types/github';
import type { RSSFeedResponse, XMLRSSItem } from '../types/rss';

// 開発環境ではViteのプロキシを使用するため、相対パスを使用
const API_BASE_URL = '';

// GitHub関連のAPI（個別にエクスポート）
export async function fetchGitHubTrending(): Promise<GitHubTrendingItem[]> {
  const response = await fetch(`${API_BASE_URL}/github-trending`);
  if (!response.ok) {
    throw new Error(`HTTPエラー: ${response.status}`);
  }
  return response.json();
}

export async function fetchGolangRepositoryTrending(): Promise<GitHubTrendingItem[]> {
  const response = await fetch(`${API_BASE_URL}/golang-repository-trending`);
  if (!response.ok) {
    throw new Error(`HTTPエラー: ${response.status}`);
  }
  return response.json();
}

export async function fetchAIRepositorySummary(url: string): Promise<SummaryResponse> {
  const response = await fetch(
    `${API_BASE_URL}/ai-repository-summary?url=${encodeURIComponent(url)}`
  );
  if (!response.ok) {
    throw new Error(`APIエラー: ${response.status}`);
  }
  return response.json();
}

export async function fetchAIArticleSummary(url: string): Promise<SummaryResponse> {
  const response = await fetch(
    `${API_BASE_URL}/ai-article-summary?url=${encodeURIComponent(url)}`
  );
  if (!response.ok) {
    throw new Error(`APIエラー: ${response.status}`);
  }
  return response.json();
}

// 後方互換性のためにオブジェクト形式でもエクスポート
export const githubApi = {
  getTrending: fetchGitHubTrending,
  getGolangRepositoryTrending: fetchGolangRepositoryTrending,
  getRepositorySummary: fetchAIRepositorySummary,
  getArticleSummary: fetchAIArticleSummary,
};

export const rssApi = {
  getInfoQ: async (language: 'ja' | 'en' = 'ja'): Promise<RSSFeedResponse> => {
    const endpoint = language === 'ja' ? '/rss-ja' : '/rss';
    const response = await fetch(`${API_BASE_URL}${endpoint}`);
    if (!response.ok) {
      throw new Error(`HTTPエラー: ${response.status}`);
    }
    return response.json();
  },

  getGolangWeekly: async (): Promise<XMLRSSItem[]> => {
    const response = await fetch(`${API_BASE_URL}/golang-weekly-content`);
    if (!response.ok) {
      throw new Error(`HTTPエラー: ${response.status}`);
    }
    const xmlText = await response.text();
    return parseXMLRSS(xmlText);
  },

  getGoogleCloud: async (language: 'ja' | 'en' = 'ja'): Promise<XMLRSSItem[]> => {
    const endpoint = language === 'ja' ? '/google-cloud-content-ja' : '/google-cloud-content';
    const response = await fetch(`${API_BASE_URL}${endpoint}`);
    if (!response.ok) {
      throw new Error(`HTTPエラー: ${response.status}`);
    }
    const xmlText = await response.text();
    return parseXMLRSS(xmlText);
  },

  getAWS: async (language: 'ja' | 'en' = 'ja'): Promise<XMLRSSItem[]> => {
    const endpoint = language === 'ja' ? '/aws-content-ja' : '/aws-content';
    const response = await fetch(`${API_BASE_URL}${endpoint}`);
    if (!response.ok) {
      throw new Error(`HTTPエラー: ${response.status}`);
    }
    const xmlText = await response.text();
    return parseXMLRSS(xmlText);
  },

  getAzure: async (language: 'ja' | 'en' = 'ja'): Promise<XMLRSSItem[]> => {
    const endpoint = language === 'ja' ? '/azure-content-ja' : '/azure-content';
    const response = await fetch(`${API_BASE_URL}${endpoint}`);
    if (!response.ok) {
      throw new Error(`HTTPエラー: ${response.status}`);
    }
    const xmlText = await response.text();
    return parseXMLRSS(xmlText);
  },
};

// XMLをパースしてRSSアイテムの配列を返す
function parseXMLRSS(xmlText: string): XMLRSSItem[] {
  const parser = new DOMParser();
  const xmlDoc = parser.parseFromString(xmlText, 'application/xml');
  const items = xmlDoc.querySelectorAll('item');

  const result: XMLRSSItem[] = [];
  items.forEach((item) => {
    const title = item.querySelector('title')?.textContent || 'No Title';
    const link = item.querySelector('link')?.textContent || '#';
    const pubDate = item.querySelector('pubDate')?.textContent || 'No Date';
    
    // descriptionはCDATAセクション内のHTMLを含むため、textContentで取得
    // textContentはCDATAの中身を正しく取得し、HTMLタグを文字列として返す
    // その後、Reactコンポーネント側でdangerouslySetInnerHTMLを使用してレンダリング
    const description = item.querySelector('description')?.textContent || 'No Description';

    result.push({ title, link, pubDate, description });
  });

  return result;
}
