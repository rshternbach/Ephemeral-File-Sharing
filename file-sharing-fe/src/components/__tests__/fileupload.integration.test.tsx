import React from 'react';
import { render, fireEvent, waitFor } from '@testing-library/react';
import FileUpload from '../FileUpload';
import path from 'path';
import fs from 'fs';
import '@testing-library/jest-dom'; 
import axios from 'axios';


jest.mock('../../config', () => ({
  config: {
    MODE: 'test',
    BACKEND_API_URL: 'http://localhost:8080',
  },
}));

describe('FileUpload Integration Test', () => {
  it('should upload file and return shareable URL', async () => {
    const { getByText, getByPlaceholderText, getByLabelText, getByRole, container } = render(<FileUpload />);

    const filePath = path.resolve(__dirname, 'test-meme.jpg');
    const fileContent = fs.readFileSync(filePath);
    const file = new File([fileContent], 'test-meme.jpg', { type: 'image/jpeg' });

    const fileInput = getByLabelText(/upload your file/i) as HTMLInputElement;
    const retentionInput = getByPlaceholderText(/retention time in minutes/i) as HTMLInputElement;
    const uploadButton = getByRole('button', { name: /upload/i });

    fireEvent.change(fileInput, { target: { files: [file] } });
    fireEvent.change(retentionInput, { target: { value: '1' } });
    fireEvent.click(uploadButton);

    await waitFor(() => {
      expect(getByText(/file uploaded successfully/i)).toBeInTheDocument();
    });

    const shareableUrlElement = await waitFor(() => container.querySelector('.shareable-url a'));

    expect(shareableUrlElement).toBeInTheDocument();
    expect((shareableUrlElement as HTMLAnchorElement).href).toMatch(/http:\/\/localhost:8080\/v1\/.*/);
        
    // Verify 200 OK response
    const shareableUrl = (shareableUrlElement as HTMLAnchorElement).href;
    const response = await axios.get(shareableUrl);
    expect(response.status).toBe(200);
    
    // Wait for 1 minute
    await new Promise(resolve => setTimeout(resolve, 60000));
    
    // Verify 404 Not Found response
    try {
    await axios.get(shareableUrl);
    } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
    expect(error.response.status).toBe(404);
      } else {
    throw error;
      }
    }
  }, 70000);

});