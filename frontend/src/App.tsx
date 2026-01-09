import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { LanguageProvider } from './contexts/LanguageContext';
import { LanguageSelector } from './components/LanguageSelector/LanguageSelector';
import { GitHubTrendingTable } from './components/GitHubTrending/GitHubTrendingTable';
import { GolangRepositoryTrendingTable } from './components/GitHubTrending/GolangRepositoryTrendingTable';
import { InfoQTable } from './components/RSS/InfoQTable';
import { RSSFeedList } from './components/RSS/RSSFeedList';
import { rssApi } from './services/api';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <LanguageProvider>
        <div style={{ fontFamily: 'Arial, sans-serif', backgroundColor: '#f9f9f9', minHeight: '100vh', padding: '20px 0' }}>
          <LanguageSelector />
          <InfoQTable />
          <GitHubTrendingTable />
          <GolangRepositoryTrendingTable />
          <RSSFeedList
            title="GolangWeekly"
            queryKey="golang-weekly"
            fetchFn={rssApi.getGolangWeekly}
            supportsLanguage={false}
          />
          <RSSFeedList
            title="Google Cloud GCP RSS Feed"
            queryKey="gcp-rss"
            fetchFn={rssApi.getGoogleCloud}
            supportsLanguage={true}
          />
          <RSSFeedList
            title="AWS RSS Feed"
            queryKey="aws-rss"
            fetchFn={rssApi.getAWS}
            supportsLanguage={true}
          />
          <RSSFeedList
            title="Azure RSS Feed"
            queryKey="azure-rss"
            fetchFn={rssApi.getAzure}
            supportsLanguage={true}
          />
        </div>
      </LanguageProvider>
    </QueryClientProvider>
  );
}

export default App;
