import Hls from "hls.js";
import React, { useEffect, useState } from 'react';

function VideoPlayer() {

  const [adminID, setAdminID] = useState('');

  useEffect(() => {
    const video = document.getElementById('videoPlayer');
    let token = localStorage.getItem('token');
    // get userId from token
    // let user_id = JSON.parse(atob(token.split('.')[1])).id;

    // Check if HLS is supported
    if (Hls.isSupported()) {
      const hls = new Hls({
        xhrSetup: xhr => {
          // Set authorization header
          xhr.setRequestHeader('Authorization', `Bearer ${token}`);
        }
      });

      // HLS stream URL
      const src = `${process.env.REACT_APP_SERVER_HOST_URL}/serve/5@18b6c71e-8921-4115-8861-13f73bb9e69f/playlist.m3u8`;
      
      // Load and attach HLS stream
      hls.loadSource(src);
      hls.attachMedia(video);

      // WebSocket connection to server
      const ws = new WebSocket(
        `${process.env.REACT_APP_SERVER_HOST_URL_WS}/room/5@18b6c71e-8921-4115-8861-13f73bb9e69f`,
      );

      // WebSocket event handlers
      ws.onopen = () => {
        console.log('WebSocket connection opened');
        ws.send(JSON.stringify({ 'msgType': 'current_status' }));
      }

      ws.onclose = (event) => {
        console.log('WebSocket connection closed:', event);
      }

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      }

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log('Message received from server:', data);
        // Handle server messages
        if (data.MsgType == 'current_status') {
          console.log('Current status:', data);
          // Set video playback position
          if (data.TimeStamp > 0) {
            video.currentTime = data.TimeStamp;
          }
          // Play video
          if (data.Status == 1) {
            video.play();
          } else {
            video.pause();
          }
        } else if (data.MsgType == 'play') {
          video.play();
        } else if (data.MsgType == 'pause') {
          console.log('pause');
          video.pause();
        } else if (data.MsgType == 'update') {
          video.currentTime = data.TimeStamp;
        }
         

        //   switch (data.action) {
        //     case 'play':
        //       video.play();
        //       break;
        //     case 'pause':
        //       video.pause();
        //       break;
        //     case 'seek':
        //       video.currentTime = data.timestamp;
        //       break;
        //     default:
        //       break;
        //   }
        // } else {
        //   console.log('Action denied by the server:', data.action);
        // }
      };

      // Event listeners for video playback control
      // video.addEventListener('play', () => {
      //   // if (ws.readyState === WebSocket.OPEN) {
      //   //   ws.send(JSON.stringify({ 'play': true }));
      //   // }
      // });

      // video.addEventListener('pause', () => {
      //   // if (ws.readyState === WebSocket.OPEN) {
      //   //   ws.send(JSON.stringify({ 'pause': true }));
      //   // }
      // });

      // video.addEventListener('seeked', () => {
      //   // if (ws.readyState === WebSocket.OPEN) {
      //   //   ws.send(JSON.stringify({ updateTimestamp: video.currentTime }));
      //   // }
      // });


    } else {
      // HLS is not supported
      console.log('HLS is not supported');
    }
  }, []);

  return (
    <div>
      <h1>HLS Video Player</h1>
      <video muted="muted" id="videoPlayer" width="640" height="360" controls>
        {/* <source src=`${serprocess.env.REACT_APP_SERVER_HOST_URLve/5@18b6c71e-8921-4115-8861-13f73bb9e69f/playlist.m3u8" type="application/x-mpegURL" `> */}
        {/* Your browser does not support the video tag. */}
      </video>
    </div>
  );
}

export default VideoPlayer;
