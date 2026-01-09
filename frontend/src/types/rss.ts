export interface RSSFeedItem {
  title: string;
  link: string;
  published: string;
  description: string;
}

export interface RSSFeedResponse {
  title: string;
  description: string;
  items: RSSFeedItem[];
}

export interface XMLRSSItem {
  title: string;
  link: string;
  pubDate: string;
  description: string;
}
