import { useState } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import { rssApi, githubApi } from '../../services/api';
import { useLanguage } from '../../contexts/LanguageContext';
import './RSSStyles.css';

export const InfoQTable = () => {
  const [modalOpen, setModalOpen] = useState(false);
  const [summaryText, setSummaryText] = useState('');
  const { language } = useLanguage();

  const { data: rssData, isLoading, error } = useQuery({
    queryKey: ['infoq-rss', language],
    queryFn: () => rssApi.getInfoQ(language),
  });

  const summaryMutation = useMutation({
    mutationFn: (url: string) => githubApi.getArticleSummary(url),
    onSuccess: (data) => {
      setSummaryText(data.summary);
      setModalOpen(true);
    },
    onError: (error) => {
      console.error('AIサマリー取得エラー:', error);
      setSummaryText('AIサマリーの取得に失敗しました。');
      setModalOpen(true);
    },
  });

  const handleGenerateSummary = (url: string) => {
    setSummaryText('');
    setModalOpen(true);
    summaryMutation.mutate(url);
  };

  const closeModal = () => {
    setModalOpen(false);
  };

  if (isLoading) {
    return (
      <div className="container">
        <h1>InfoQ latest news</h1>
        <p>Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container">
        <h1>InfoQ latest news</h1>
        <p style={{ color: 'red' }}>RSSフィードの取得に失敗しました。</p>
      </div>
    );
  }

  return (
    <div className="container">
      <h1>InfoQ latest news</h1>
      <table id="rss-feed-table">
        <thead>
          <tr>
            <th>title</th>
            <th>description</th>
            <th>publication date</th>
            <th>link</th>
            <th>AI</th>
          </tr>
        </thead>
        <tbody>
          {rssData?.items.map((item, index) => (
            <tr key={index}>
              <td>{item.title || 'No title'}</td>
              <td>
                <div
                  dangerouslySetInnerHTML={{ __html: item.description || 'No description' }}
                  className="description-content"
                />
              </td>
              <td>{item.published || 'No date'}</td>
              <td>
                <a href={item.link} target="_blank" rel="noopener noreferrer">
                  URL
                </a>
              </td>
              <td>
                <button
                  className="rss-button"
                  onClick={() => handleGenerateSummary(item.link)}
                  disabled={summaryMutation.isPending}
                >
                  Generate
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* モーダル */}
      {modalOpen && (
        <div className="modal" onClick={closeModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <span className="close-button" onClick={closeModal}>
              &times;
            </span>
            <h2>AIサマリー</h2>
            {summaryMutation.isPending ? (
              <div id="loading-indicator">processing...</div>
            ) : (
              <p id="summary-text">{summaryText}</p>
            )}
          </div>
        </div>
      )}
    </div>
  );
};
