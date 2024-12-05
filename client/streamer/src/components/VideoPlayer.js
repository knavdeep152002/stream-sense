import Hls from "hls.js";
import React, { useEffect, useState } from 'react';
import { useParams } from "react-router-dom";
import './test.css';
import { Box, Button, Text, VStack, useToast } from '@chakra-ui/react';


function VideoPlayer() {

  // const [adminID, setAdminID] = useState(0);
  const [error, setError] = useState('');
  const toast = useToast();
  const { link } = useParams();
  let adminID = 2;
  useEffect(() => {
    const video = document.getElementById('videoPlayer');
    let token = localStorage.getItem('token');
    // get userId from token
    let user_id = '';
    if (token) {
      user_id = JSON.parse(atob(token.split('.')[1]))?.id || '';
    }
    // Check if HLS is supported
    if (Hls.isSupported()) {
      const hls = new Hls({
        xhrSetup: xhr => {
          // Set authorization header
          xhr.setRequestHeader('Authorization', `Bearer ${token}`);
        }
      });

      // HLS stream URL
      const src = `${process.env.REACT_APP_SERVER_HOST_URL}/serve/${link}/playlist.m3u8`;
      
      // Load and attach HLS stream
      hls.loadSource(src);
      hls.attachMedia(video);

      // WebSocket connection to server
      const ws = new WebSocket(
        `${process.env.REACT_APP_SERVER_HOST_URL_WS}/room/${link}`,
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

        switch (data.MsgType) {
          case 'current_status':
            // Set video playback position
            if (data.TimeStamp > 0) {
              video.currentTime = data.TimeStamp;
            }
            // if (data.AdminId) {
            //   setAdminID(data.AdminId);
            // }
            // Play video
            if (data.Status === 1) {
              video.play();
            } else {
              video.pause();
            }
            break;
          case 'play':
            if (user_id !== adminID) {
              video.play();
            }
            break;
          case 'pause':
            if (user_id !== adminID) {
              video.pause();
            }
            break;
          case 'update':
            if (user_id !== adminID) {
              console.log('update event', user_id, adminID, data.TimeStamp);
              video.currentTime = data.TimeStamp;
            }
            break;
          default:
            break;
        }
      };
    
      if (user_id === adminID) {
        video.controls = true;
        // video.controlsList = 'download remoteplayback foobar';
        console.log('Admin controls enabled');
        // Event listeners for video playback control
        video.addEventListener('play', () => {
          if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ 'play': true, "msgType": "control" }));
          }
        });

        video.addEventListener('pause', () => {
          console.log("asdasd", ws.readyState, WebSocket.OPEN)
          if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ 'pause': true, "msgType": "control" }));
          }
        });

        video.addEventListener('seeked', () => {
          if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ updateTimestamp: video.currentTime, "msgType": "control" }));
          }
        });
      } else{
        console.log('Admin controls disabled');
        video.classList.add('no-controls');
      }


    } else {
      // HLS is not supported
      setError('HLS is not supported by your browser.');
    }
  }, [adminID, link]);

  const handleCopyLink = () => {
    navigator.clipboard.writeText(`${window.location.origin}/room/${link}`);
    toast({
      title: "Link copied!",
      description: "Video link has been copied to your clipboard.",
      status: "success",
      duration: 2000,
      isClosable: true,
    });
  };


  return (
    <VStack spacing={4} align="center" mt={12}>
      <Text fontSize="xl" fontWeight="bold">
        Stream Sense Video Player
      </Text>
      <Box>
        <video controls controlsList="noplay noremoteplayback noplaybackrate foobar" muted="muted" playsInline="" autoPlay="" id="videoPlayer" width="640" height="360">
          {/* <source src=`${serprocess.env.REACT_APP_SERVER_HOST_URLve/${link}/playlist.m3u8" type="application/x-mpegURL" `> */}
          {/* Your browser does not support the video tag. */}
        </video>
        </Box>
        <Button
        colorScheme="teal"
        onClick={handleCopyLink}
      >
        Share Link
      </Button>
    </VStack>
  );
}

export default VideoPlayer;
