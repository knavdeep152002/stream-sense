import { Flex, Heading, List, Alert, AlertIcon, Box, Text, Button, SimpleGrid } from '@chakra-ui/react';
import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';

function HomePage() {
  const [videos, setVideos] = useState([]);
  const [error, setError] = useState('');
  const [hoveredVideo, setHoveredVideo] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchVideos = async () => {
      const token = localStorage.getItem('token');
      const response = await fetch(`${process.env.REACT_APP_SERVER_HOST_URL}/uploads`, {
        headers: {
          'Authorization': 'Bearer ' + token,
        },
      });

      if (response.status > 299) {
        if (response.status === 401) {
          navigate('/login');  // Redirect if unauthorized
        } else {
          setError('Failed to load videos. Please try again.');
        }
        return;
      }

      const data = await response.json();
      setVideos(data.data);
    };

    fetchVideos();
  }, [navigate]);

  const createRoom = (videoID) => async () => {
    const token = localStorage.getItem('token');
    const response = await fetch(`${process.env.REACT_APP_SERVER_HOST_URL}/room/${videoID}`, {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + token,
        'Content-Type': 'application/json',
      },
    });

    if (response.status > 299) {
      if (response.status === 401) {
        navigate('/login');
      } else {
        setError('Error creating room. Please try again.');
      }
      return;
    }

    const data = await response.json();
    window.location.href = `/room/${data.data.room_link}`;
  };

  return (
    <Flex direction="column" align="center">
      <Flex direction="row" align="center" justify="space-between" width="100%" mb={4}>
        <Button as={Link} to="/upload" colorScheme="teal" size="lg">
          Upload Video
        </Button>
        <Button
          onClick={() => {
            localStorage.removeItem('token');
            navigate('/login');
          }}
          colorScheme="red"
          size="lg"
        >
          Logout
        </Button>
      </Flex>

      {error && (
        <Alert status="error" mt={4}>
          <AlertIcon />
          {error}
        </Alert>
      )}

      <Heading as="h1" fontSize="3xl" mt={8} mb={6} color="teal.600" textAlign="center">
        Your Videos
      </Heading>

      {/* Video Grid */}
      <SimpleGrid columns={[1, 2, 3, 4, 5]} spacing={6} width="100%" maxWidth="1200px" mx="auto">
        {videos?.map((video) => (
          <Box
            key={video.id}
            width="100%"
            height="150px"
            bg="gray.200"
            borderRadius="md"
            boxShadow="md"
            p={4}
            cursor="pointer"
            onClick={createRoom(video.id)}
            onMouseEnter={() => setHoveredVideo(video.id)}
            onMouseLeave={() => setHoveredVideo(null)}
            _hover={{ bg: 'gray.300', transform: 'scale(1.02)' }}
            transition="transform 0.2s"
          >
            <Text>{video.filename}</Text>
            {hoveredVideo === video.id && (
              <Text fontSize="sm" color="gray.600">
                Click to create a room
              </Text>
            )}
          </Box>
        ))}
      </SimpleGrid>
    </Flex>
  );
}

export default HomePage;
