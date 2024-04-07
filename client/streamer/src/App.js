import './App.css';
import { useState } from 'react';

function App() {
  const [file, setFile] = useState(null);
  const [totalChunks, setTotalChunks] = useState(0);
  const [totalFileSize, setTotalFileSize] = useState(0);

  const handleFileChange = (event) => {
    const selectedFile = event.target.files[0];
    setFile(selectedFile);
    
    // Calculate total chunks and total file size
    const chunkSize = 1024 * 1024; // size of each chunk (1MB)
    const fileSize = selectedFile.size;
    const chunks = Math.ceil(fileSize / chunkSize);
    
    setTotalChunks(chunks);
    setTotalFileSize(fileSize);
  };

  const handleFileUpload = () => {
    if (!file) return;

    const chunkSize = 1024 * 1024; // size of each chunk (1MB)
    const uploadId = Math.random().toString(36).substring(7);;
    let start = 0;

    console.log('file size:', file.size);
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
    fetch(`http://localhost:8000/stream-sense/api/v1/complete?upload_id=${uploadId}&file_name=${fileName}`, {
      method: 'POST',
    });
  }

  const uploadChunk = (chunk, uploadId, chunkNumber) => {
    const formData = new FormData();
    formData.append('upload_id', uploadId);
    formData.append('chunk_number', chunkNumber); // Assuming chunk_number starts from 1
    formData.append('total_chunks', totalChunks); // Use state value
    formData.append('total_file_size', totalFileSize); // Use state value
    formData.append('file_name', file.name);
    formData.append('file', chunk); // 'file' is the actual file to be uploaded

    // Make a request to the server
    fetch('http://localhost:8000/stream-sense/api/v1/upload', {
      method: 'POST',
      body: formData,
    });
  };

  return (
    <div className="App">
      <header className="App-header">
        <form>
          <input type="file" onChange={handleFileChange}></input>
          <button type="button" onClick={handleFileUpload}>Upload File</button>
        </form>
      </header>
    </div>
  );
}

export default App;
