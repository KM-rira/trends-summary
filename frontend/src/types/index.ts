// GitHub Trending の型定義
export interface GitHubTrendingItem {
  name: string;
  description: string;
  language: string;
  stars: string;
  url: string;
}

// AI Summary の型定義
export interface AISummaryResponse {
  summary: string;
}

// RSS Feed の型定義
export interface RSSFeedItem {
  title: string;
  description: string;
  published: string;
  link: string;
}

export interface RSSFeedResponse {
  items: RSSFeedItem[];
}
