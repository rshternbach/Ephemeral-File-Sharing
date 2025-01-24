import React from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import FileUpload from './components/FileUpload';
import './App.css';

const queryClient = new QueryClient();

const App = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="App">
        <h1>Ephemeral File Sharing</h1>
        <FileUpload />
      </div>
    </QueryClientProvider>
  );
};

export default App;