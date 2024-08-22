import { Button, Flex, Input } from '@chakra-ui/react';
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';


function Login({ setToken }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();
  const handleLogin = async (e) => {
    let host = process.env.REACT_APP_SERVER_HOST_URL;
    console.log('host', process.env);
    const response = await fetch(`${host}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });
    const data = await response.json();
    if (data.token) {
      localStorage.setItem('token', data.token);
      setToken(data.token);
    }
    navigate('/home');
  };

  return (
    <Flex justify={'center'} align={'center'} width={'100vh'} h={'100vh'}>
        <Flex direction={'column'} justify={'center'} width={'500px'} height={'400px'} gap={'20px'} border={'1px solid #ccc'} p={'20px'} justifyContent={'center'} alignItems={'center'}>
          <Input 
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <Input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <Button onClick={()=>{handleLogin()}} type="submit">Login</Button>
          <Link to="/register">Create a user</Link>
        </Flex>
      </Flex>
  );
}

export default Login;
