import { useQuery } from '@tanstack/react-query';
import type { XMLRSSItem } from '../../types/rss';
import { useLanguage } from '../../contexts/LanguageContext';
import './RSSStyles.css';

interface RSSFeedListProps {
  title: string;
  queryKey: string;
  fetchFn: (language?: 'ja' | 'en') => Promise<XMLRSSItem[]>;
  supportsLanguage?: boolean;
}

export const RSSFeedList = ({ title, queryKey, fetchFn, supportsLanguage = true }: RSSFeedListProps) => {
  const { language } = useLanguage();
  
  const { data: rssData, isLoading, error } = useQuery({
    queryKey: supportsLanguage ? [queryKey, language] : [queryKey],
    queryFn: () => supportsLanguage ? fetchFn(language) : fetchFn(),
  });

  if (isLoading) {
    return (
      <div className="container">
        <h1>{title}</h1>
        <p>Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container">
        <h1>{title}</h1>
        <p style={{ color: 'red' }}>フィードの取得に失敗しました。</p>
      </div>
    );
  }

  if (!rssData || rssData.length === 0) {
    return (
      <div className="container">
        <h1>{title}</h1>
        <p>No articles found.</p>
      </div>
    );
  }

  return (
    <div className="container">
      <h1>{title}</h1>
      <ul className="rss-list">
        {rssData.map((item, index) => (
          <li key={index} className="rss-list-item">
            <h3>
              <a href={item.link} target="_blank" rel="noopener noreferrer">
                {item.title}
              </a>
            </h3>
            <p>
              <strong>Published:</strong> {item.pubDate}
            </p>
            <div
              className="rss-description"
              dangerouslySetInnerHTML={{ __html: item.description }}
            />
            <hr />
          </li>
        ))}
      </ul>
    </div>
  );
};
