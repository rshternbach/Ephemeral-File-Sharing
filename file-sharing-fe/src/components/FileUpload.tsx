import React, { useState, ChangeEvent } from 'react';
import useFileUpload from '../hooks/useFileUpload';
import './FileUpload.css';

const FileUpload: React.FC = () => {
  const [file, setFile] = useState<File | null>(null);
  const [retentionTime, setRetentionTime] = useState<number>(1);
  const { shareableUrl, error, uploadFile } = useFileUpload();

  const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      setFile(event.target.files[0]);
    }
  };

  const handleRetentionTimeChange = (event: ChangeEvent<HTMLInputElement>) => {
    setRetentionTime(Number(event.target.value));
  };

  const handleSubmit = () => {
    if (file) {
      uploadFile(file, retentionTime);
    }
  };

  const handleCopyUrl = () => {
    if (shareableUrl) {
      navigator.clipboard.writeText(shareableUrl);
    }
  };

  return (
    <div className="file-upload-card">
      <h2>Upload Your File</h2>
      <label htmlFor="file-input">Upload your file</label>
      <input id="file-input" type="file" onChange={handleFileChange} className="file-input" />
      <label htmlFor="retention-input">Retention time in minutes</label>
      <input
        type="number"
        value={retentionTime}
        onChange={handleRetentionTimeChange}
        placeholder="Retention time in minutes"
        className="retention-input"
        min="1"
      />
      <button onClick={handleSubmit} className="upload-button">Upload</button>

      {error && <p className="error">{error}</p>}

      {shareableUrl && (
        <div className="modal">
          <p>File uploaded successfully!</p>
          <p className="shareable-url">
            Shareable URL: <a href={shareableUrl} target="_blank" rel="noopener noreferrer" title={shareableUrl}>{shareableUrl}</a>
          </p>
          <button onClick={handleCopyUrl} className="copy-button">Copy URL</button>
        </div>
      )}
    </div>
  );
};

export default FileUpload;