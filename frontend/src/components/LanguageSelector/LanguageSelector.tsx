import { useLanguage } from '../../contexts/LanguageContext';
import './LanguageSelector.css';

export function LanguageSelector() {
  const { language, setLanguage } = useLanguage();
  
  console.log('🌐 LanguageSelector rendered, current language:', language);

  return (
    <div className="language-selector" style={{ display: 'flex' }}>
      <label htmlFor="language-select">🌐 Language: </label>
      <select
        id="language-select"
        value={language}
        onChange={(e) => {
          console.log('🔄 Language changed to:', e.target.value);
          setLanguage(e.target.value as 'ja' | 'en');
        }}
      >
        <option value="ja">🇯🇵 日本語</option>
        <option value="en">🇬🇧 English</option>
      </select>
    </div>
  );
}
