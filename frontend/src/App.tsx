import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { LanguageProvider } from './contexts/LanguageContext';
import { LoginPage } from './components/Login/LoginPage';
import { LanguageSelector } from './components/LanguageSelector/LanguageSelector';
import { GitHubTrendingTable } from './components/GitHubTrending/GitHubTrendingTable';
import { GolangRepositoryTrendingTable } from './components/GitHubTrending/GolangRepositoryTrendingTable';
import { InfoQTable } from './components/RSS/InfoQTable';
import { RSSFeedList } from './components/RSS/RSSFeedList';
import { rssApi } from './services/api';
import './App.css';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function MainApp() {
  const { isAuthenticated, isLoading, logout } = useAuth();

  if (isLoading) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        minHeight: '100vh',
        fontSize: '1.2rem',
        color: '#666'
      }}>
        読み込み中...
      </div>
    );
  }

  if (!isAuthenticated) {
    return <LoginPage />;
  }

  return (
    <div style={{ fontFamily: 'Arial, sans-serif', backgroundColor: '#f9f9f9', minHeight: '100vh', padding: '20px 0' }}>
      <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '0 20px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
          <LanguageSelector />
          <button 
            onClick={logout}
            style={{
              padding: '8px 16px',
              backgroundColor: '#dc3545',
              color: 'white',
              border: 'none',
              borderRadius: '6px',
              cursor: 'pointer',
              fontSize: '0.9rem',
              fontWeight: '500'
            }}
          >
            ログアウト
          </button>
        </div>
        <InfoQTable />
        <GitHubTrendingTable />
        <GolangRepositoryTrendingTable />
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
        <RSSFeedList
          title="GolangWeekly"
          queryKey="golang-weekly"
          fetchFn={rssApi.getGolangWeekly}
          supportsLanguage={false}
        />
      </div>
    </div>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <LanguageProvider>
          <MainApp />
        </LanguageProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
