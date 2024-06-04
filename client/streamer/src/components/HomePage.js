import { Flex, Heading, List } from '@chakra-ui/react';
import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

function HomePage() {
  const [videos, setVideos] = useState([]);

  useEffect(() => {
    const fetchVideos = async () => {
      const token = localStorage.getItem('token');
      const response = await fetch(`${process.env.REACT_APP_SERVER_HOST_URL}/uploads`, {
        headers: {
          'Authorization': 'Bearer ' + token,
        },
      });
      const data = await response.json();
      setVideos(data.data);
    };

    fetchVideos();
  }, []);

  return (
    <Flex direction="column" align="center">
      <Flex direction="row" align="center" justify="space-between" width="100%">
        <Link to="/upload">Upload Video</Link>
        <Link onClick={() => localStorage.removeItem('token')} to="/login">Logout</Link>
      </Flex>
      <Heading>Your Videos</Heading>
      <Flex direction="column" align="center">
        <List spacing={3}>
          {videos.map((video) => (
            <li key={video.id}>{video.filename}</li>
          ))}
        </List>
      </Flex>  
    </Flex>
  );
}

export default HomePage;
