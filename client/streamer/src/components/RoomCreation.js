import React from 'react';
import { useHistory } from 'react-router-dom';

function RoomCreation() {
  const history = useHistory();

  const createRoom = async () => {
    const token = localStorage.getItem('token');
    const videoId = 'video-id';
    const response = await fetch(`${process.env.REACT_APP_SERVER_HOST_URL}/room/${videoId}`, {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + token,
      },
    });
    const data = await response.json();
    if (data.roomId) {
      history.push(`/room/${data.roomId}`);
    }
  };

  return (
    <div>
      <button onClick={createRoom}>Create Room</button>
    </div>
  );
}

export default RoomCreation;
