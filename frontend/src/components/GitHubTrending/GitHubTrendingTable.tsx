import { useState } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import ReactMarkdown from 'react-markdown';
import { fetchGitHubTrending, fetchAIRepositorySummary } from '../../services/api';
import type { GitHubTrendingItem } from '../../types/github';
import './GitHubTrendingTable.css';

export function GitHubTrendingTable() {
  const [modalOpen, setModalOpen] = useState(false);
  const [summaryText, setSummaryText] = useState('');
  const [isLoadingSummary, setIsLoadingSummary] = useState(false);

  // GitHubトレンドデータの取得
  const { data: trendingData, isLoading, error } = useQuery({
    queryKey: ['github-trending'],
    queryFn: fetchGitHubTrending,
  });

  // AIサマリー取得のMutation
  const summaryMutation = useMutation({
    mutationFn: fetchAIRepositorySummary,
    onMutate: () => {
      setIsLoadingSummary(true);
      setSummaryText('');
      setModalOpen(true);
    },
    onSuccess: (data) => {
      setSummaryText(data.summary);
      setIsLoadingSummary(false);
    },
    onError: (error) => {
      console.error('AIサマリー取得エラー:', error);
      setSummaryText('AIサマリーの取得に失敗しました。');
      setIsLoadingSummary(false);
    },
  });

  const handleGenerateSummary = (url: string) => {
    summaryMutation.mutate(url);
  };

  const closeModal = () => {
    setModalOpen(false);
  };

  if (isLoading) {
    return (
      <div className="container">
        <h1>GitHub daily trends</h1>
        <p>読み込み中...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container">
        <h1>GitHub daily trends</h1>
        <p>GitHubトレンドデータの取得に失敗しました。</p>
      </div>
    );
  }

  return (
    <>
      <div className="container">
        <h1>GitHub daily trends</h1>
        <table id="github-trending-table">
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
      </div>

      {/* モーダル */}
      {modalOpen && (
        <div className="modal" onClick={closeModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <span className="close-button" onClick={closeModal}>
              &times;
            </span>
            <h2>AIサマリー</h2>
            {isLoadingSummary ? (
              <div id="loading-indicator">処理中...</div>
            ) : (
              <div id="summary-text" className="markdown-content">
                <ReactMarkdown>{summaryText}</ReactMarkdown>
              </div>
            )}
          </div>
        </div>
      )}
    </>
  );
}
