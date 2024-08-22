import React, { useState } from 'react';

function UploadVideo() {
  const [file, setFile] = useState(null);
  const [totalChunks, setTotalChunks] = useState(0);
  const [totalFileSize, setTotalFileSize] = useState(0);

  const handleFileChange = (event) => {
    const selectedFile = event.target.files[0];
    setFile(selectedFile);

    const chunkSize = 1024 * 1024; // 1MB
    const fileSize = selectedFile.size;
    const chunks = Math.ceil(fileSize / chunkSize);

    setTotalChunks(chunks);
    setTotalFileSize(fileSize);
  };

  const handleFileUpload = () => {
    if (!file) return;

    const chunkSize = 1024 * 1024; // 1MB
    const uploadId = Math.random().toString(36).substring(7);
    let start = 0;

    let chunkNumber = 0;
    while (start < file.size) {
      const chunk = file.slice(start, start + chunkSize);
      chunkNumber++;
      uploadChunk(chunk, uploadId, chunkNumber);
      start += chunkSize;
    }

    setTimeout(() => {
      handleFileComplete(uploadId, file.name);
    }, 1000);
  };

  const handleFileComplete = (uploadId, fileName) => {
    fetch(`${process.env.REACT_APP_SERVER_HOST_URL}/complete?upload_id=${uploadId}&file_name=${fileName}`, {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('token'),
      },
    });
  };

  const uploadChunk = (chunk, uploadId, chunkNumber) => {
    const formData = new FormData();
    formData.append('upload_id', uploadId);
    formData.append('chunk_number', chunkNumber);
    formData.append('total_chunks', totalChunks);
    formData.append('total_file_size', totalFileSize);
    formData.append('file_name', file.name);
    formData.append('file', chunk);

    fetch(`${process.env.REACT_APP_SERVER_HOST_URL}/upload`, {
      method: 'POST',
      body: formData,
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('token'),
      },
    });
  };

  return (
    <div>
      <input type="file" onChange={handleFileChange} />
      <button type="button" onClick={handleFileUpload}>Upload File</button>
    </div>
  );
}

export default UploadVideo;
