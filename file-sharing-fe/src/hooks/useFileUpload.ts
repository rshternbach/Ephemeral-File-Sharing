import { useState } from 'react';
import axios from 'axios';
import log from '../utils/logger';
import { config } from '../config';

interface UseFileUploadReturn {
  shareableUrl: string | null;
  error: string | null;
  uploadFile: (file: File, retentionTime: number) => Promise<void>;
}

const useFileUpload = (): UseFileUploadReturn => {
  const [shareableUrl, setShareableUrl] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const uploadFile = async (file: File, retentionTime: number): Promise<void> => {
    const formData = new FormData();
    formData.append('file', file);
    try {
        log.info('Starting file upload', { fileName: file.name, retentionTime });
        const formData = new FormData();
        formData.append('file', file);
        formData.append('retentionTime', retentionTime.toString());
    
   
          const response = await axios.put(`${config.BACKEND_API_URL}/v1/file`, formData);
          log.info('File upload successful', { url: response.data.url });
          setShareableUrl(`${config.BACKEND_API_URL}/v1/${response.data.url}`);
        setError(null);
    } catch (error) {
      log.error('Error uploading file', error);

      setError('Error uploading file');
    }
  };

  return { shareableUrl, error, uploadFile };
};

export default useFileUpload;