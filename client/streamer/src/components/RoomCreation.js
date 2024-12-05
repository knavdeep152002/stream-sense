import React, { useState } from 'react';
import { useHistory } from 'react-router-dom';
import { Alert, AlertIcon, Box } from '@chakra-ui/react';

function RoomCreation() {
  const history = useHistory();
  const [error, setError] = useState('');

  const createRoom = async () => {
    const token = localStorage.getItem('token');
    const videoId = 'video-id';
    const response = await fetch(`${process.env.REACT_APP_SERVER_HOST_URL}/room/${videoId}`, {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    });

    if (response.status > 299) {
      if (response.status === 401) {
        history.push('/login');
      } else {
        setError('Failed to create room. Please try again.');
      }
      return;
    }

    const data = await response.json();
    if (data.roomId) {
      history.push(`/room/${data.roomId}`);
    }
  };

  return (
    <Box>
      {error && (
        <Alert status="error">
          <AlertIcon />
          {error}
        </Alert>
      )}
      <button onClick={createRoom}>Create Room</button>
    </Box>
  );
}

export default RoomCreation;
