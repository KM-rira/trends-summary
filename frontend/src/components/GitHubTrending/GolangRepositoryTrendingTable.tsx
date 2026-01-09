import { useState } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import ReactMarkdown from 'react-markdown';
import { githubApi } from '../../services/api';
import type { GitHubTrendingItem } from '../../types/github';
import './GitHubTrendingTable.css';

export const GolangRepositoryTrendingTable = () => {
  const [modalOpen, setModalOpen] = useState(false);
  const [summaryText, setSummaryText] = useState('');
  
  const { data: trendingData, isLoading, error } = useQuery({
    queryKey: ['golang-repository-trending'],
    queryFn: githubApi.getGolangRepositoryTrending,
  });

  const summaryMutation = useMutation({
    mutationFn: (url: string) => githubApi.getRepositorySummary(url),
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
        <h1>Golang repository daily trends</h1>
        <p>Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container">
        <h1>Golang repository daily trends</h1>
        <p style={{ color: 'red' }}>Golangリポジトリトレンドデータの取得に失敗しました。</p>
      </div>
    );
  }

  return (
    <div className="container">
      <h1>Golang repository daily trends</h1>
      <table id="golang-repository-trending-table">
        <thead>
          <tr>
            <th>name</th>
            <th>description</th>
            <th>language</th>
            <th>stars</th>
            <th>link</th>
            <th>AI</th>
          </tr>
        </thead>
        <tbody>
          {trendingData?.map((item: GitHubTrendingItem, index: number) => (
            <tr key={index}>
              <td>{item.name || 'No name'}</td>
              <td>{item.description || 'No description'}</td>
              <td>{item.language || 'N/A'}</td>
              <td>{item.stars || '0'}</td>
              <td>
                <a href={item.url} target="_blank" rel="noopener noreferrer">
                  URL
                </a>
              </td>
              <td>
                <button
                  className="rss-button"
                  onClick={() => handleGenerateSummary(item.url)}
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
              <div id="summary-text" className="markdown-content">
                <ReactMarkdown>{summaryText}</ReactMarkdown>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};
