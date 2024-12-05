import React, { useState } from 'react';
import { Box, Button, Input, Text, VStack, useToast } from '@chakra-ui/react';
import { useNavigate } from 'react-router-dom';

function UploadVideo() {
  const [file, setFile] = useState(null);
  const [isUploading, setIsUploading] = useState(false);
  const toast = useToast();
  const navigate = useNavigate();

  const handleFileChange = (event) => {
    const selectedFile = event.target.files[0];
    setFile(selectedFile);
  };

  const handleFileUpload = async () => {
    if (!file) {
      toast({
        title: "No file selected.",
        status: "warning",
        duration: 2000,
        isClosable: true,
      });
      return;
    }

    setIsUploading(true);

    const chunkSize = 1024 * 1024; // 1MB
    const uploadId = Math.random().toString(36).substring(7);
    let start = 0;
    let chunkNumber = 0;

    while (start < file.size) {
      const chunk = file.slice(start, start + chunkSize);
      chunkNumber++;
      await uploadChunk(chunk, uploadId, chunkNumber);
      start += chunkSize;
    }

    setTimeout(() => {
      handleFileComplete(uploadId, file.name);
    }, 1000);
  };

  const uploadChunk = (chunk, uploadId, chunkNumber) => {
    const formData = new FormData();
    formData.append('upload_id', uploadId);
    formData.append('chunk_number', chunkNumber);
    formData.append('file_name', file.name);
    formData.append('file', chunk);

    return fetch(`${process.env.REACT_APP_SERVER_HOST_URL}/upload`, {
      method: 'POST',
      body: formData,
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('token'),
      },
    });
  };

  const handleFileComplete = (uploadId, fileName) => {
    fetch(`${process.env.REACT_APP_SERVER_HOST_URL}/complete?upload_id=${uploadId}&file_name=${fileName}`, {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('token'),
      },
    }).then(response => {
      if (response.status > 299) {
        toast({
          title: "Upload failed.",
          status: "error",
          duration: 2000,
          isClosable: true,
        });
      } else {
        toast({
          title: "Video uploaded successfully!",
          status: "success",
          duration: 2000,
          isClosable: true,
        });
        navigate('/');  // Navigate back to the homepage after upload
      }
    });
  };

  return (
    <VStack spacing={6} align="center">
      <Box
        p={6}
        borderWidth={2}
        borderColor="gray.400"
        borderRadius="md"
        bg="gray.50"
        _hover={{ bg: "gray.100" }}
        cursor="pointer"
      >
        <Text mb={2} fontWeight="bold" fontSize="lg">
          {file ? file.name : "Click to select a file"}
        </Text>
        <Input type="file" onChange={handleFileChange} hidden />
      </Box>

      <Button
        colorScheme="teal"
        size="lg"
        isLoading={isUploading}
        onClick={handleFileUpload}
      >
        Upload File
      </Button>
    </VStack>
  );
}

export default UploadVideo;
