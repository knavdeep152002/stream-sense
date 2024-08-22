import { Children, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
const Auth = ({token, children}) => {
    console.log("asdasdsad")
    const navigate = useNavigate();
    useEffect(() => {
      if (!token) {
        navigate('/login')
      }
    }, [token]);

    return (
        <div>
            {Children.map(children, child => 
            child
            )}
        </div>
        
    )
  
}

export default Auth